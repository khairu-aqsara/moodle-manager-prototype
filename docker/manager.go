package docker

import (
	"fmt"
	"strings"
	"time"
	
	"moodle-prototype-manager/utils"
)

const (
	ContainerPort = "8080:8080"
)

// Manager handles Docker container operations
type Manager struct{
	imageName string
}

// NewManager creates a new Docker manager
func NewManager() *Manager {
	return &Manager{}
}

// SetImageName sets the Docker image name to use
func (m *Manager) SetImageName(imageName string) {
	m.imageName = imageName
}

// GetImageName returns the current Docker image name
func (m *Manager) GetImageName() string {
	return m.imageName
}

// CheckImageExists verifies if the Moodle image exists locally
func (m *Manager) CheckImageExists() (bool, error) {
	if m.imageName == "" {
		return false, fmt.Errorf("no image name set")
	}
	
	cmd := GetDockerCommand("images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check images: %w (output: %s)", err, string(output))
	}
	
	exists := strings.Contains(string(output), m.imageName)
	utils.LogDebug(fmt.Sprintf("Image check - looking for: %s, exists: %v", m.imageName, exists))
	return exists, nil
}

// PullImage downloads the Moodle Docker image
func (m *Manager) PullImage() error {
	if m.imageName == "" {
		return fmt.Errorf("no image name set")
	}

	utils.LogInfo(fmt.Sprintf("Pulling Docker image: %s", m.imageName))
	cmd := GetDockerCommand("pull", m.imageName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w (output: %s)", m.imageName, err, string(output))
	}
	return nil
}

// PullImageWithProgress downloads the Docker image with progress tracking
func (m *Manager) PullImageWithProgress(progressCallback func(float64, string)) error {
	if m.imageName == "" {
		return fmt.Errorf("no image name set")
	}

	utils.LogInfo(fmt.Sprintf("Pulling Docker image with progress: %s", m.imageName))

	// Create command but don't run it yet
	cmd := GetDockerCommand("pull", m.imageName)

	// Get stdout pipe for reading progress
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Get stderr pipe as Docker may output to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start docker pull: %w", err)
	}

	// Create progress tracker
	progress := NewPullProgress()
	if progressCallback != nil {
		progress.AddCallback(progressCallback)
	}

	// Process output in separate goroutines
	errChan := make(chan error, 2)

	// Process stdout
	go func() {
		if err := progress.ProcessStream(stdout); err != nil {
			errChan <- fmt.Errorf("error processing stdout: %w", err)
		} else {
			errChan <- nil
		}
	}()

	// Process stderr (Docker sometimes outputs progress here)
	go func() {
		if err := progress.ProcessStream(stderr); err != nil {
			errChan <- fmt.Errorf("error processing stderr: %w", err)
		} else {
			errChan <- nil
		}
	}()

	// Wait for command to complete
	cmdErr := cmd.Wait()

	// Wait for both stream processors to complete
	streamErr1 := <-errChan
	streamErr2 := <-errChan

	// Check for errors
	if cmdErr != nil {
		return fmt.Errorf("docker pull failed: %w", cmdErr)
	}

	if streamErr1 != nil {
		utils.LogWarning(fmt.Sprintf("Stream processing warning: %v", streamErr1))
	}

	if streamErr2 != nil {
		utils.LogWarning(fmt.Sprintf("Stream processing warning: %v", streamErr2))
	}

	utils.LogInfo("Docker image pulled successfully with progress tracking")
	return nil
}

// RunContainer starts a new Moodle container
func (m *Manager) RunContainer() (string, error) {
	if m.imageName == "" {
		return "", fmt.Errorf("no image name set")
	}
	
	utils.LogInfo(fmt.Sprintf("Running container from image: %s", m.imageName))
	cmd := GetDockerCommand("run", "-d", "-p", ContainerPort, m.imageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run container: %w (output: %s)", err, string(output))
	}
	
	containerID := strings.TrimSpace(string(output))
	utils.LogInfo(fmt.Sprintf("Container started with ID: %s", containerID))
	return containerID, nil
}

// StartContainer starts an existing container
func (m *Manager) StartContainer(containerID string) error {
	cmd := GetDockerCommand("start", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Docker start command failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return fmt.Errorf("failed to start container %s: %w (output: %s)", containerID, err, string(output))
	}
	return nil
}

// StopContainer stops a running container gracefully
func (m *Manager) StopContainer(containerID string) error {
	cmd := GetDockerCommand("stop", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Docker stop command failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return fmt.Errorf("failed to stop container %s: %w (output: %s)", containerID, err, string(output))
	}
	return nil
}

// IsContainerRunning checks if a container is currently running
func (m *Manager) IsContainerRunning(containerID string) (bool, error) {
	cmd := GetDockerCommand("inspect", "--format={{.State.Running}}", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Docker inspect command failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return false, fmt.Errorf("failed to inspect container %s: %w (output: %s)", containerID, err, string(output))
	}
	
	return strings.TrimSpace(string(output)) == "true", nil
}

// GetContainerLogs retrieves logs from a container
func (m *Manager) GetContainerLogs(containerID string) (string, error) {
	cmd := GetDockerCommand("logs", containerID)
	
	// Use CombinedOutput to capture both stdout and stderr
	// Docker logs may output to stderr on some platforms, especially Windows
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Docker logs command failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	
	return string(output), nil
}

// GetContainerLogsSince retrieves logs from a container since a specific time
func (m *Manager) GetContainerLogsSince(containerID string, since time.Time) (string, error) {
	// Format time for docker logs --since flag
	// Docker accepts RFC3339 format
	sinceStr := since.Format(time.RFC3339)
	
	cmd := GetDockerCommand("logs", "--since", sinceStr, containerID)
	
	// Use CombinedOutput to capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Docker logs --since command failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return "", fmt.Errorf("failed to get container logs since %s: %w", sinceStr, err)
	}
	
	return string(output), nil
}





// ValidateContainerID checks if a container ID is valid and exists
func (m *Manager) ValidateContainerID(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is empty")
	}
	
	// Check if container exists by trying to inspect it
	cmd := GetDockerCommand("inspect", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Container validation failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return fmt.Errorf("container %s does not exist or cannot be accessed: %w", containerID, err)
	}
	
	return nil
}

// ForceStopContainer forcefully stops a container (used as last resort)
func (m *Manager) ForceStopContainer(containerID string) error {
	cmd := GetDockerCommand("kill", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogError("Docker kill command failed", fmt.Errorf("%v (output: %s)", err, string(output)))
		return fmt.Errorf("failed to force stop container %s: %w (output: %s)", containerID, err, string(output))
	}
	return nil
}