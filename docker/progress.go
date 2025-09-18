package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"

	"moodle-prototype-manager/utils"
)

// DockerPullEvent represents a single progress event from Docker
type DockerPullEvent struct {
	Status         string                 `json:"status"`
	ID             string                 `json:"id"`
	Progress       string                 `json:"progress"`
	ProgressDetail ProgressDetail         `json:"progressDetail"`
	Error          string                 `json:"error"`
}

// ProgressDetail contains download/extract progress information
type ProgressDetail struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

// LayerProgress tracks progress for a single layer
type LayerProgress struct {
	ID              string
	Status          string
	DownloadCurrent int64
	DownloadTotal   int64
	ExtractCurrent  int64
	ExtractTotal    int64
}

// PullProgress manages overall pull progress
type PullProgress struct {
	layers    map[string]*LayerProgress
	mu        sync.RWMutex
	callbacks []func(float64, string)
}

// NewPullProgress creates a new progress tracker
func NewPullProgress() *PullProgress {
	return &PullProgress{
		layers:    make(map[string]*LayerProgress),
		callbacks: make([]func(float64, string), 0),
	}
}

// AddCallback registers a callback for progress updates
func (p *PullProgress) AddCallback(callback func(float64, string)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.callbacks = append(p.callbacks, callback)
}

// ProcessStream reads and processes Docker output stream
func (p *PullProgress) ProcessStream(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Try to parse as JSON first
		var event DockerPullEvent
		if err := json.Unmarshal([]byte(line), &event); err == nil {
			// Process JSON event
			if err := p.processEvent(&event); err != nil {
				utils.LogError("Failed to process Docker event", err)
			}
		} else {
			// Not JSON, parse plain text output
			p.processPlainTextLine(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading Docker output: %w", err)
	}

	return nil
}

// processPlainTextLine handles plain text Docker output
func (p *PullProgress) processPlainTextLine(line string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Parse different plain text formats
	// Example formats:
	// "4f4fb700ef54: Pulling fs layer"
	// "8cc6894b165e: Downloading  12.3MB/45.6MB"
	// "8cc6894b165e: Extracting  12.3MB/45.6MB"
	// "4f4fb700ef54: Pull complete"
	// "Status: Downloaded newer image for..."

	// Extract layer ID if present
	layerID := ""
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && len(parts[0]) == 12 { // Docker layer IDs are 12 chars
			layerID = parts[0]
			line = strings.TrimSpace(parts[1])
		}
	}

	// Handle status messages
	if strings.HasPrefix(line, "Status:") {
		status := strings.TrimPrefix(line, "Status:")
		p.notifyCallbacks(-1, strings.TrimSpace(status))
		return
	}

	// Handle "Pulling from" message
	if strings.Contains(line, "Pulling from") {
		p.notifyCallbacks(0, "Starting download...")
		return
	}

	// Skip non-layer lines
	if layerID == "" {
		utils.LogDebug(fmt.Sprintf("Non-layer Docker output: %s", line))
		return
	}

	// Get or create layer progress
	layer, exists := p.layers[layerID]
	if !exists {
		layer = &LayerProgress{
			ID: layerID,
		}
		p.layers[layerID] = layer
	}

	// Parse the status
	if strings.Contains(line, "Pulling fs layer") {
		layer.Status = "Preparing"
	} else if strings.Contains(line, "Waiting") {
		layer.Status = "Waiting"
	} else if strings.Contains(line, "Already exists") {
		layer.Status = "Already exists"
		// Set as complete
		layer.DownloadCurrent = 1
		layer.DownloadTotal = 1
		layer.ExtractCurrent = 1
		layer.ExtractTotal = 1
	} else if strings.Contains(line, "Pull complete") {
		layer.Status = "Pull complete"
		// Ensure everything is marked as complete
		if layer.DownloadTotal > 0 {
			layer.DownloadCurrent = layer.DownloadTotal
		} else {
			layer.DownloadCurrent = 1
			layer.DownloadTotal = 1
		}
		if layer.ExtractTotal > 0 {
			layer.ExtractCurrent = layer.ExtractTotal
		} else {
			layer.ExtractCurrent = 1
			layer.ExtractTotal = 1
		}
	} else if strings.Contains(line, "Download complete") {
		layer.Status = "Download complete"
		if layer.DownloadTotal > 0 {
			layer.DownloadCurrent = layer.DownloadTotal
		}
	} else if strings.Contains(line, "Downloading") || strings.Contains(line, "Extracting") {
		// Parse size information
		// Format: "Downloading  12.3MB/45.6MB" or "Extracting  [====>  ] 12.3MB/45.6MB"
		isExtracting := strings.Contains(line, "Extracting")

		// Find the size information (look for MB/GB pattern)
		sizePattern := `(\d+(?:\.\d+)?)\s*([KMGT]B)\/(\d+(?:\.\d+)?)\s*([KMGT]B)`
		re := regexp.MustCompile(sizePattern)
		matches := re.FindStringSubmatch(line)

		if len(matches) == 5 {
			current := parseSize(matches[1], matches[2])
			total := parseSize(matches[3], matches[4])

			if isExtracting {
				layer.Status = "Extracting"
				layer.ExtractCurrent = current
				layer.ExtractTotal = total
			} else {
				layer.Status = "Downloading"
				layer.DownloadCurrent = current
				layer.DownloadTotal = total
			}
		}
	}

	// Calculate and notify progress
	percentage := p.calculateOverallProgress()
	status := p.getOverallStatus()
	p.notifyCallbacks(percentage, status)
}

// parseSize converts size string to bytes
func parseSize(value string, unit string) int64 {
	var multiplier int64 = 1
	switch unit {
	case "KB":
		multiplier = 1024
	case "MB":
		multiplier = 1024 * 1024
	case "GB":
		multiplier = 1024 * 1024 * 1024
	case "TB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}

	var size float64
	fmt.Sscanf(value, "%f", &size)
	return int64(size * float64(multiplier))
}

// processEvent handles a single Docker pull event
func (p *PullProgress) processEvent(event *DockerPullEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check for errors
	if event.Error != "" {
		return fmt.Errorf("docker error: %s", event.Error)
	}

	// Skip events without layer ID
	if event.ID == "" {
		// These are usually summary messages like "Pulling from library/..."
		utils.LogDebug(fmt.Sprintf("Docker status: %s", event.Status))
		p.notifyCallbacks(-1, event.Status)
		return nil
	}

	// Get or create layer progress
	layer, exists := p.layers[event.ID]
	if !exists {
		layer = &LayerProgress{
			ID: event.ID,
		}
		p.layers[event.ID] = layer
	}

	// Update layer status
	layer.Status = event.Status

	// Update progress based on status
	switch event.Status {
	case "Downloading":
		if event.ProgressDetail.Total > 0 {
			layer.DownloadCurrent = event.ProgressDetail.Current
			layer.DownloadTotal = event.ProgressDetail.Total
		}
	case "Extracting":
		if event.ProgressDetail.Total > 0 {
			layer.ExtractCurrent = event.ProgressDetail.Current
			layer.ExtractTotal = event.ProgressDetail.Total
		}
	case "Pull complete":
		// Mark both download and extract as complete
		layer.DownloadCurrent = layer.DownloadTotal
		layer.ExtractCurrent = layer.ExtractTotal
	case "Already exists":
		// This layer is already downloaded
		layer.DownloadCurrent = layer.DownloadTotal
		layer.ExtractCurrent = layer.ExtractTotal
	}

	// Calculate and notify overall progress
	percentage := p.calculateOverallProgress()
	status := p.getOverallStatus()
	p.notifyCallbacks(percentage, status)

	return nil
}

// calculateOverallProgress calculates the total progress percentage
func (p *PullProgress) calculateOverallProgress() float64 {
	if len(p.layers) == 0 {
		return 0
	}

	// Count different types of layers and collect meaningful byte data
	var cachedLayers, activeLayers, completedLayers int
	var totalDownloadBytes, currentDownloadBytes int64
	var totalExtractBytes, currentExtractBytes int64
	var layersWithMeaningfulBytes int

	for _, layer := range p.layers {
		switch layer.Status {
		case "Already exists":
			cachedLayers++
		case "Pull complete", "Download complete":
			completedLayers++
		default:
			activeLayers++
		}

		// Only count bytes for layers that are actually being downloaded/extracted
		// Skip cached layers and only include meaningful byte data (> 1 byte)
		if layer.Status != "Already exists" {
			// Check if this layer has meaningful byte data (not dummy values)
			hasMeaningfulDownload := layer.DownloadTotal > 1
			hasMeaningfulExtract := layer.ExtractTotal > 1

			if hasMeaningfulDownload {
				totalDownloadBytes += layer.DownloadTotal
				currentDownloadBytes += layer.DownloadCurrent
				layersWithMeaningfulBytes++
			}
			if hasMeaningfulExtract {
				totalExtractBytes += layer.ExtractTotal
				currentExtractBytes += layer.ExtractCurrent
				if !hasMeaningfulDownload {
					layersWithMeaningfulBytes++
				}
			}
		}
	}

	totalLayers := len(p.layers)
	actualWorkLayers := totalLayers - cachedLayers // Layers that need actual work

	utils.LogDebug(fmt.Sprintf("Layer analysis - Total: %d, Cached: %d, Active: %d, Completed: %d, Work needed: %d, With meaningful bytes: %d",
		totalLayers, cachedLayers, activeLayers, completedLayers, actualWorkLayers, layersWithMeaningfulBytes))

	// If all layers are cached, we're done
	if cachedLayers == totalLayers {
		utils.LogDebug("All layers cached, returning 100%")
		return 100
	}

	// If no layers need work (all complete), we're done
	if completedLayers + cachedLayers == totalLayers {
		utils.LogDebug("All layers complete or cached, returning 100%")
		return 100
	}

	// Use byte-based calculation only if we have meaningful byte data from multiple layers
	// or if the byte data represents a significant portion of the work
	shouldUseByteBased := (totalDownloadBytes > 0 || totalExtractBytes > 0) &&
						  (layersWithMeaningfulBytes >= 2 ||
						   float64(layersWithMeaningfulBytes)/float64(actualWorkLayers) > 0.3)

	if !shouldUseByteBased {
		// Use layer-based progress calculation
		workCompletedLayers := completedLayers
		if actualWorkLayers > 0 {
			progress := float64(workCompletedLayers) / float64(actualWorkLayers) * 100
			utils.LogDebug(fmt.Sprintf("Layer-based progress: %d/%d work layers complete = %.1f%% (insufficient meaningful byte data)",
				workCompletedLayers, actualWorkLayers, progress))
			return progress
		}
		return 0
	}

	// Calculate weighted progress (download is 60%, extract is 40%)
	// Only for layers that actually have meaningful byte data
	var downloadProgress, extractProgress float64

	if totalDownloadBytes > 0 {
		downloadProgress = float64(currentDownloadBytes) / float64(totalDownloadBytes) * 60
	}

	if totalExtractBytes > 0 {
		extractProgress = float64(currentExtractBytes) / float64(totalExtractBytes) * 40
	}

	totalProgress := downloadProgress + extractProgress

	// For mixed scenarios (some layers with bytes, some without), we need to account for
	// completed layers that don't have byte data
	if layersWithMeaningfulBytes < actualWorkLayers {
		// Calculate how much progress the non-byte-tracked completed layers contribute
		layersWithoutBytes := actualWorkLayers - layersWithMeaningfulBytes
		completedLayersWithoutBytes := 0

		for _, layer := range p.layers {
			if (layer.Status == "Pull complete" || layer.Status == "Download complete") &&
			   layer.Status != "Already exists" &&
			   layer.DownloadTotal <= 1 && layer.ExtractTotal <= 1 {
				completedLayersWithoutBytes++
			}
		}

		if layersWithoutBytes > 0 {
			// Add progress contribution from completed layers without byte tracking
			nonByteProgress := float64(completedLayersWithoutBytes) / float64(layersWithoutBytes) * 100
			// Weight the byte-based progress by the portion of layers that have byte data
			byteWeight := float64(layersWithMeaningfulBytes) / float64(actualWorkLayers)
			nonByteWeight := float64(layersWithoutBytes) / float64(actualWorkLayers)

			totalProgress = (totalProgress * byteWeight) + (nonByteProgress * nonByteWeight)

			utils.LogDebug(fmt.Sprintf("Mixed progress - Byte layers: %d/%d, Non-byte layers: %d/%d, Weighted total: %.1f%%",
				layersWithMeaningfulBytes, actualWorkLayers, completedLayersWithoutBytes, layersWithoutBytes, totalProgress))
		}
	}

	// Ensure we don't exceed 100%
	if totalProgress > 100 {
		totalProgress = 100
	}

	utils.LogDebug(fmt.Sprintf("Byte-based progress - Download: %.1f%% (%.1fMB/%.1fMB), Extract: %.1f%% (%.1fMB/%.1fMB), Total: %.1f%%",
		downloadProgress, float64(currentDownloadBytes)/(1024*1024), float64(totalDownloadBytes)/(1024*1024),
		extractProgress, float64(currentExtractBytes)/(1024*1024), float64(totalExtractBytes)/(1024*1024),
		totalProgress))

	return totalProgress
}

// getOverallStatus returns a human-readable status message
func (p *PullProgress) getOverallStatus() string {
	downloadingCount := 0
	extractingCount := 0
	completeCount := 0
	cachedCount := 0
	preparingCount := 0

	for _, layer := range p.layers {
		switch layer.Status {
		case "Downloading":
			downloadingCount++
		case "Extracting":
			extractingCount++
		case "Pull complete", "Download complete":
			completeCount++
		case "Already exists":
			cachedCount++
		case "Preparing", "Waiting", "Pulling fs layer":
			preparingCount++
		}
	}

	totalLayers := len(p.layers)
	workLayers := totalLayers - cachedCount // Layers that need actual work
	workCompleted := completeCount

	// If everything is cached, show appropriate message
	if cachedCount == totalLayers && totalLayers > 0 {
		return "Image already available"
	}

	// Show downloading status
	if downloadingCount > 0 {
		return fmt.Sprintf("Downloading layers (%d/%d completed)", workCompleted, workLayers)
	}

	// Show extracting status
	if extractingCount > 0 {
		return fmt.Sprintf("Extracting layers (%d/%d completed)", workCompleted, workLayers)
	}

	// Show completion status
	if workCompleted == workLayers && workLayers > 0 {
		return "Pull complete"
	}

	// For preparing/waiting states, just show a simple initializing message
	if preparingCount > 0 {
		return "Initializing download..."
	}

	return "Starting download..."
}

// notifyCallbacks notifies all registered callbacks of progress update
func (p *PullProgress) notifyCallbacks(percentage float64, status string) {
	for _, callback := range p.callbacks {
		callback(percentage, status)
	}
}