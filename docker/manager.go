package docker

import (
	"fmt"
	"strings"
	"time"

	"moodle-prototype-manager/errors"
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
		return false, errors.NewValidationError("imageName", "no image name set in Docker manager", "")
	}

	// Validate image name format
	if err := errors.ValidateImageName(m.imageName); err != nil {
		return false, errors.WrapWithContext(err, "invalid image name in Docker manager")
	}

	cmd := GetDockerCommand("images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithImage("check", m.imageName, err).WithOutput(string(output))
		return false, errors.WrapWithContext(dockerErr, "failed to execute docker images command")
	}

	exists := strings.Contains(string(output), m.imageName)
	utils.LogDebug(fmt.Sprintf("Image check - looking for: %s, exists: %v", m.imageName, exists))
	return exists, nil
}

// PullImage downloads the Moodle Docker image
func (m *Manager) PullImage() error {
	if m.imageName == "" {
		return errors.NewValidationError("imageName", "no image name set in Docker manager", "")
	}

	// Validate image name format
	if err := errors.ValidateImageName(m.imageName); err != nil {
		return errors.WrapWithContext(err, "invalid image name for pull operation")
	}

	utils.LogInfo(fmt.Sprintf("Pulling Docker image: %s", m.imageName))
	cmd := GetDockerCommand("pull", m.imageName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithImage("pull", m.imageName, err).WithOutput(string(output))
		return errors.WrapWithContext(dockerErr, "failed to pull Docker image")
	}
	return nil
}

// PullImageWithProgress downloads the Docker image with progress tracking
func (m *Manager) PullImageWithProgress(progressCallback func(float64, string)) error {
	if m.imageName == "" {
		return errors.NewValidationError("imageName", "no image name set in Docker manager", "")
	}

	// Validate image name format
	if err := errors.ValidateImageName(m.imageName); err != nil {
		return errors.WrapWithContext(err, "invalid image name for pull with progress operation")
	}

	utils.LogInfo(fmt.Sprintf("Pulling Docker image with progress: %s", m.imageName))

	// Create command but don't run it yet
	cmd := GetDockerCommand("pull", m.imageName)

	// Get stdout pipe for reading progress
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithImage("pull_setup", m.imageName, err)
		return errors.WrapWithContext(dockerErr, "failed to create stdout pipe for Docker pull")
	}

	// Get stderr pipe as Docker may output to stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithImage("pull_setup", m.imageName, err)
		return errors.WrapWithContext(dockerErr, "failed to create stderr pipe for Docker pull")
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		dockerErr := errors.NewDockerErrorWithImage("pull_start", m.imageName, err)
		return errors.WrapWithContext(dockerErr, "failed to start docker pull command")
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
		dockerErr := errors.NewDockerErrorWithImage("pull", m.imageName, cmdErr)
		return errors.WrapWithContext(dockerErr, "docker pull command failed")
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
		return "", errors.NewValidationError("imageName", "no image name set in Docker manager", "")
	}

	// Validate image name format
	if err := errors.ValidateImageName(m.imageName); err != nil {
		return "", errors.WrapWithContext(err, "invalid image name for run container operation")
	}

	utils.LogInfo(fmt.Sprintf("Running container from image: %s", m.imageName))
	cmd := GetDockerCommand("run", "-d", "-p", ContainerPort, m.imageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithImage("run", m.imageName, err).WithOutput(string(output))
		return "", errors.WrapWithContext(dockerErr, "failed to run new container")
	}

	containerID := strings.TrimSpace(string(output))

	// Validate the returned container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return "", errors.WrapWithContext(err, "Docker returned invalid container ID: %s", containerID)
	}

	utils.LogInfo(fmt.Sprintf("Container started with ID: %s", containerID))
	return containerID, nil
}

// StartContainer starts an existing container
func (m *Manager) StartContainer(containerID string) error {
	// Validate container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return errors.WrapWithContext(err, "invalid container ID provided to StartContainer")
	}

	cmd := GetDockerCommand("start", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("start", containerID, err).WithOutput(string(output))
		utils.LogError("Docker start command failed", dockerErr)
		return errors.WrapWithContext(dockerErr, "failed to start existing container")
	}
	return nil
}

// StopContainer stops a running container gracefully
func (m *Manager) StopContainer(containerID string) error {
	// Validate container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return errors.WrapWithContext(err, "invalid container ID provided to StopContainer")
	}

	cmd := GetDockerCommand("stop", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("stop", containerID, err).WithOutput(string(output))
		utils.LogError("Docker stop command failed", dockerErr)
		return errors.WrapWithContext(dockerErr, "failed to stop container gracefully")
	}
	return nil
}

// IsContainerRunning checks if a container is currently running
func (m *Manager) IsContainerRunning(containerID string) (bool, error) {
	// Validate container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return false, errors.WrapWithContext(err, "invalid container ID provided to IsContainerRunning")
	}

	cmd := GetDockerCommand("inspect", "--format={{.State.Running}}", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("inspect", containerID, err).WithOutput(string(output))
		utils.LogError("Docker inspect command failed", dockerErr)
		return false, errors.WrapWithContext(dockerErr, "failed to inspect container status")
	}

	return strings.TrimSpace(string(output)) == "true", nil
}

// GetContainerLogs retrieves logs from a container
func (m *Manager) GetContainerLogs(containerID string) (string, error) {
	// Validate container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return "", errors.WrapWithContext(err, "invalid container ID provided to GetContainerLogs")
	}

	cmd := GetDockerCommand("logs", containerID)

	// Use CombinedOutput to capture both stdout and stderr
	// Docker logs may output to stderr on some platforms, especially Windows
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("logs", containerID, err).WithOutput(string(output))
		utils.LogError("Docker logs command failed", dockerErr)
		return "", errors.WrapWithContext(dockerErr, "failed to retrieve container logs")
	}

	return string(output), nil
}

// GetContainerLogsSince retrieves logs from a container since a specific time
func (m *Manager) GetContainerLogsSince(containerID string, since time.Time) (string, error) {
	// Validate container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return "", errors.WrapWithContext(err, "invalid container ID provided to GetContainerLogsSince")
	}

	// Validate time parameter
	if since.IsZero() {
		return "", errors.NewValidationError("since", "since time cannot be zero", since)
	}

	// Format time for docker logs --since flag
	// Docker accepts RFC3339 format
	sinceStr := since.Format(time.RFC3339)

	cmd := GetDockerCommand("logs", "--since", sinceStr, containerID)

	// Use CombinedOutput to capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("logs_since", containerID, err).WithOutput(string(output))
		utils.LogError("Docker logs --since command failed", dockerErr)
		return "", errors.WrapWithContext(dockerErr, "failed to retrieve container logs since %s", sinceStr)
	}

	return string(output), nil
}





// ValidateContainerID checks if a container ID is valid and exists
func (m *Manager) ValidateContainerID(containerID string) error {
	// Basic validation first
	if err := errors.ValidateContainerID(containerID); err != nil {
		return errors.WrapWithContext(err, "container ID format validation failed")
	}

	// Check if container exists by trying to inspect it
	cmd := GetDockerCommand("inspect", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("inspect", containerID, err).WithOutput(string(output))
		utils.LogError("Container validation failed", dockerErr)
		return errors.WrapWithContext(dockerErr, "container does not exist or cannot be accessed")
	}

	return nil
}

// ForceStopContainer forcefully stops a container (used as last resort)
func (m *Manager) ForceStopContainer(containerID string) error {
	// Validate container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return errors.WrapWithContext(err, "invalid container ID provided to ForceStopContainer")
	}

	cmd := GetDockerCommand("kill", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dockerErr := errors.NewDockerErrorWithContainer("kill", containerID, err).WithOutput(string(output))
		utils.LogError("Docker kill command failed", dockerErr)
		return errors.WrapWithContext(dockerErr, "failed to force stop container")
	}
	return nil
}