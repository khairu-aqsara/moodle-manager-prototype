package main

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"moodle-prototype-manager/docker"
	"moodle-prototype-manager/storage"
	"moodle-prototype-manager/utils"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx              context.Context
	dockerManager    *docker.Manager
	credentialManager *storage.CredentialManager
	fileManager      *storage.FileManager
	logParser        *docker.LogParser
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Initialize logging
	utils.InitLogger()
	utils.LogInfo("Initializing Moodle Prototype Manager")
	
	return &App{
		dockerManager:    docker.NewManager(),
		credentialManager: storage.NewCredentialManager(),
		fileManager:      storage.NewFileManager(),
		logParser:        docker.NewLogParser(),
	}
}

// OnStartup is called when the app starts
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	
	// Load image configuration
	imageName, err := a.fileManager.LoadImageName()
	if err != nil {
		utils.LogError("Failed to load image configuration", err)
		// Return error rather than using potentially wrong fallback
		// This will make the issue visible to users so they can fix it
		utils.LogError("Cannot start application without valid image configuration", fmt.Errorf("image.docker file missing or unreadable"))
		// Use a safe fallback but log it prominently
		imageName = "wenkhairu/moodle-prototype:502-stable"
		utils.LogWarning(fmt.Sprintf("FALLBACK: Using default image '%s' - please create image.docker file with correct image name", imageName))
	}
	
	// Set the image name in Docker manager
	a.dockerManager.SetImageName(imageName)
	utils.LogInfo(fmt.Sprintf("Using Docker image: %s", imageName))
	
	utils.LogInfo("Application startup completed")
}

// OnShutdown is called when the app is shutting down
func (a *App) OnShutdown(ctx context.Context) {
	utils.LogInfo("Application shutdown initiated")
	
	// Check if container is running and stop it gracefully
	if !a.fileManager.ContainerIDExists() {
		utils.LogInfo("No container ID file found during shutdown")
		return
	}

	containerID, err := a.fileManager.LoadContainerID()
	if err != nil {
		utils.LogError("Failed to load container ID during shutdown", err)
		return
	}

	utils.LogInfo(fmt.Sprintf("Checking container status during shutdown: %s", containerID))
	running, err := a.dockerManager.IsContainerRunning(containerID)
	if err != nil {
		utils.LogError("Failed to check container status during shutdown", err)
		// Still try to stop it anyway as a failsafe
		utils.LogWarning("Attempting failsafe container stop during shutdown")
		if stopErr := a.dockerManager.StopContainer(containerID); stopErr != nil {
			utils.LogError("Failsafe container stop failed during shutdown", stopErr)
		}
		return
	}

	if running {
		utils.LogInfo("Stopping running container on app shutdown...")
		err := a.dockerManager.StopContainer(containerID)
		if err != nil {
			utils.LogError("Failed to stop container during shutdown", err)
		} else {
			utils.LogInfo("Container stopped successfully during shutdown")
		}
	} else {
		utils.LogInfo("Container is already stopped during shutdown")
	}
}


// HealthCheck performs Docker and Internet connectivity checks
func (a *App) HealthCheck() map[string]bool {
	utils.LogInfo("Frontend requested health check")
	
	healthStatus := docker.PerformHealthChecks()
	
	result := map[string]bool{
		"docker":   healthStatus.Docker,
		"internet": healthStatus.Internet,
	}
	
	utils.LogInfo(fmt.Sprintf("Returning health status to frontend: %+v", result))
	return result
}

// RunMoodle starts the Moodle container
func (a *App) RunMoodle() error {
	utils.LogInfo("RunMoodle called")
	
	// For existing containers, we'll preserve the password and only update after container is ready
	// For new containers, we'll clear to start fresh
	
	// Check if container already exists
	if a.fileManager.ContainerIDExists() {
		containerID, err := a.fileManager.LoadContainerID()
		if err == nil {
			utils.LogInfo(fmt.Sprintf("Found existing container ID: %s", containerID))
			
			// Try to start existing container
			running, err := a.dockerManager.IsContainerRunning(containerID)
			if err == nil {
				if running {
					utils.LogWarning("Container is already running")
					return fmt.Errorf("container is already running")
				}
				// Start existing container
				utils.LogInfo("Starting existing container")
				
				// Record the time before starting to only look for new logs
				startTime := time.Now()
				
				err := a.dockerManager.StartContainer(containerID)
				if err != nil {
					return fmt.Errorf("failed to start existing container: %w", err)
				}
				
				// Wait for existing container to be ready and extract credentials
				utils.LogInfo("Waiting for existing container to be ready...")
				go a.waitForContainerAndExtractCredentialsSince(containerID, startTime)
				
				return nil
			}
			utils.LogWarning(fmt.Sprintf("Error checking container status: %v", err))
		}
	}

	// First-time setup: check if image exists
	utils.LogInfo("Checking if Docker image exists")
	utils.LogInfo(fmt.Sprintf("Current image name: %s", a.dockerManager.GetImageName()))
	
	// Ensure we have an image name
	if a.dockerManager.GetImageName() == "" {
		utils.LogError("No Docker image name configured", nil)
		return fmt.Errorf("no Docker image name configured - please check image.docker file")
	}
	
	imageExists, err := a.dockerManager.CheckImageExists()
	if err != nil {
		utils.LogError("Failed to check image", err)
		return fmt.Errorf("failed to check Docker image: %w", err)
	}

	// Pull image if it doesn't exist
	if !imageExists {
		utils.LogInfo("Docker image not found, pulling with progress tracking...")

		// Use PullImageWithProgress to track download progress
		err := a.dockerManager.PullImageWithProgress(func(percentage float64, status string) {
			// Emit progress event to frontend
			progressData := map[string]interface{}{
				"percentage": percentage,
				"status":     status,
			}
			wailsruntime.EventsEmit(a.ctx, "docker:pull:progress", progressData)
			utils.LogDebug(fmt.Sprintf("Pull progress: %.1f%% - %s", percentage, status))
		})

		if err != nil {
			utils.LogError("Failed to pull image with progress", err)
			return fmt.Errorf("failed to pull image: %w", err)
		}
		utils.LogInfo("Docker image pulled successfully")
	} else {
		utils.LogInfo("Docker image already exists")
	}

	// Clear old credentials for new container
	utils.LogInfo("Clearing old credentials for new container")
	if err := a.credentialManager.Clear(); err != nil {
		utils.LogWarning(fmt.Sprintf("Failed to clear old credentials: %v", err))
	}

	// Run new container
	utils.LogInfo("Running new container")
	
	// Record the time before starting to only look for new logs
	startTime := time.Now()
	
	containerID, err := a.dockerManager.RunContainer()
	if err != nil {
		utils.LogError("Failed to run container", err)
		return fmt.Errorf("failed to run container: %w", err)
	}
	utils.LogInfo(fmt.Sprintf("Container started with ID: %s", containerID))

	// Save container ID
	if err := a.fileManager.SaveContainerID(containerID); err != nil {
		utils.LogError("Failed to save container ID", err)
		return fmt.Errorf("failed to save container ID: %w", err)
	}

	// Wait for container to be ready and extract credentials
	// Use the new method that only looks at logs since container start
	go a.waitForContainerAndExtractCredentialsSince(containerID, startTime)

	return nil
}

// StopMoodle stops the Moodle container
func (a *App) StopMoodle() error {
	utils.LogInfo("StopMoodle called")
	
	if !a.fileManager.ContainerIDExists() {
		utils.LogError("No container ID file found", nil)
		return fmt.Errorf("no container ID found")
	}

	containerID, err := a.fileManager.LoadContainerID()
	if err != nil {
		utils.LogError("Failed to load container ID", err)
		return fmt.Errorf("failed to load container ID: %w", err)
	}

	utils.LogInfo(fmt.Sprintf("Attempting to stop container: %s", containerID))
	
	// Validate container exists
	if err := a.dockerManager.ValidateContainerID(containerID); err != nil {
		utils.LogError("Container validation failed", err)
		return fmt.Errorf("container validation failed: %w", err)
	}
	
	// Check if container is actually running
	running, err := a.dockerManager.IsContainerRunning(containerID)
	if err != nil {
		utils.LogError("Failed to check container status", err)
		// Still try to stop it anyway
		utils.LogWarning("Attempting to stop container despite status check failure")
	} else if !running {
		utils.LogInfo("Container is already stopped")
		return nil
	}

	// Try graceful stop first
	err = a.dockerManager.StopContainer(containerID)
	if err != nil {
		utils.LogError("Graceful stop failed, attempting force stop", err)
		
		// Try force stop as fallback
		forceErr := a.dockerManager.ForceStopContainer(containerID)
		if forceErr != nil {
			utils.LogError("Force stop also failed", forceErr)
			return fmt.Errorf("failed to stop container (graceful: %v, force: %v)", err, forceErr)
		}
		
		utils.LogWarning("Container force stopped successfully")
		return nil
	}
	
	utils.LogInfo("Container stopped gracefully")
	return nil
}

// GetCredentials retrieves stored Moodle credentials
func (a *App) GetCredentials() map[string]string {
	creds, err := a.credentialManager.Load()
	if err != nil {
		// Return default credentials if loading fails
		return storage.DefaultCredentials().ToMap()
	}
	
	return creds.ToMap()
}

// IsContainerReady checks if the container is ready
func (a *App) IsContainerReady() bool {
	utils.LogDebug("Frontend called IsContainerReady()")
	
	// If we have existing credentials, check if Moodle is responding
	if a.fileManager.ContainerIDExists() {
		utils.LogDebug("Container exists, testing HTTP availability")
		// For existing containers, test HTTP availability
		if a.testMoodleHTTP() {
			utils.LogDebug("HTTP test passed - container is ready")
			return true
		}
		utils.LogDebug("HTTP test failed - container not ready yet")
		return false
	}
	
	utils.LogDebug("No existing container, checking credentials file")
	// Fallback: check if credentials file exists (for first runs)
	result := a.credentialManager.Exists()
	utils.LogDebug(fmt.Sprintf("Credentials file exists: %v", result))
	return result
}

// OpenBrowser opens the default browser to the Moodle URL
func (a *App) OpenBrowser() error {
	creds, err := a.credentialManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}
	
	if creds.URL == "" {
		return fmt.Errorf("no URL available")
	}
	
	// Use different commands based on the platform
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", creds.URL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", creds.URL)
	case "linux":
		cmd = exec.Command("xdg-open", creds.URL)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	
	utils.SetupCommandForPlatform(cmd)
	return cmd.Start()
}


// waitForContainerAndExtractCredentialsSince waits for container startup and extracts credentials
func (a *App) waitForContainerAndExtractCredentialsSince(containerID string, since time.Time) {
	utils.LogInfo("Starting to wait for container and extract credentials")
	start := time.Now()

	// For subsequent runs, check if we already have credentials saved
	existingCreds, err := a.credentialManager.Load()
	hasExistingPassword := err == nil && existingCreds.Password != ""
	
	if hasExistingPassword {
		utils.LogInfo("Subsequent run - testing HTTP availability instead of parsing logs")
		// For subsequent runs, reasonable timeout since container should start quickly
		subsequentTimeout := 10 * time.Minute
		for time.Since(start) < subsequentTimeout {
			if a.testMoodleHTTP() {
				utils.LogInfo("Container is ready - Moodle is responding on HTTP")
				// Use existing password with default URL
				a.credentialManager.Update(existingCreds.Password, "http://localhost:8080")
				utils.LogInfo("Updated credentials with existing password")
				return
			}
			
			utils.LogDebug("Waiting for Moodle HTTP response...")
			time.Sleep(2 * time.Second)
		}
		
		utils.LogError("Timeout waiting for Moodle HTTP response", nil)
		return
	}

	// First run - extract credentials from logs
	utils.LogInfo("First run - extracting credentials from logs")
	// For first runs, we don't set a timeout limit because Windows installations can take 20-30+ minutes
	// The loop will continue indefinitely until credentials are found or the application is closed
	
	logErrorCount := 0
	maxLogErrors := 5 // Allow some log errors before increasing sleep time
	
	for {
		logs, err := a.dockerManager.GetContainerLogs(containerID)
		if err != nil {
			logErrorCount++
			utils.LogDebug(fmt.Sprintf("Error getting container logs (count: %d): %v", logErrorCount, err))
			
			// If we have many consecutive log errors, increase sleep time to reduce spam
			if logErrorCount > maxLogErrors {
				utils.LogWarning("Multiple log errors detected, increasing poll interval")
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(2 * time.Second)
			}
			continue
		}
		
		// Reset error count on successful log retrieval
		logErrorCount = 0

		// First run - extract both password and URL from logs
		creds := a.logParser.ExtractCredentials(logs)
		utils.LogDebug(fmt.Sprintf("Credentials extracted - Password: %s, URL: %s", 
			maskPassword(creds.Password), creds.URL))
		
		if creds.IsComplete() {
			a.credentialManager.Update(creds.Password, creds.URL)
			utils.LogInfo("Credentials extracted and saved successfully")
			return
		}

		time.Sleep(2 * time.Second)
	}

	// Note: This function now runs indefinitely for first runs until credentials are found
}

// testMoodleHTTP tests if Moodle is responding on port 8080
func (a *App) testMoodleHTTP() bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get("http://localhost:8080")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	// Any HTTP response (even 500) means the server is up
	return resp.StatusCode > 0
}


// GetImageName returns the current Docker image name for the frontend
func (a *App) GetImageName() string {
	return a.dockerManager.GetImageName()
}

// maskPassword masks password for logging
func maskPassword(password string) string {
	if len(password) > 4 {
		return password[:2] + "****" + password[len(password)-2:]
	}
	return "****"
}