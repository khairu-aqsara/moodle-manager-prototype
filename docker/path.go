package docker

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	
	"moodle-prototype-manager/utils"
)

var dockerPath string

// FindDockerPath attempts to locate the Docker executable
func FindDockerPath() (string, error) {
	// If we already found it, return cached path
	if dockerPath != "" {
		return dockerPath, nil
	}

	// First, try the standard 'docker' command (works if Docker is in PATH)
	if path, err := exec.LookPath("docker"); err == nil {
		dockerPath = path
		return dockerPath, nil
	}

	// If not found in PATH, try common installation locations
	var commonPaths []string

	switch runtime.GOOS {
	case "darwin": // macOS
		commonPaths = []string{
			"/usr/local/bin/docker",
			"/opt/homebrew/bin/docker", // Apple Silicon Homebrew
			"/Applications/Docker.app/Contents/Resources/bin/docker",
			"/usr/bin/docker",
		}
	case "windows":
		// Try common Windows locations including environment-based paths
		commonPaths = []string{
			"docker.exe", // Try with .exe extension first
			"C:\\Program Files\\Docker\\Docker\\resources\\bin\\docker.exe",
			"C:\\Program Files (x86)\\Docker\\Docker\\resources\\bin\\docker.exe",
			"C:\\ProgramData\\DockerDesktop\\version-bin\\docker.exe",
		}
		
		// Add paths based on environment variables (like the health check does)
		if programFiles := os.Getenv("PROGRAMFILES"); programFiles != "" {
			commonPaths = append(commonPaths, programFiles+"\\Docker\\Docker\\resources\\bin\\docker.exe")
		}
		if programFilesX86 := os.Getenv("PROGRAMFILES(X86)"); programFilesX86 != "" {
			commonPaths = append(commonPaths, programFilesX86+"\\Docker\\Docker\\resources\\bin\\docker.exe")
		}
	case "linux":
		commonPaths = []string{
			"/usr/bin/docker",
			"/usr/local/bin/docker",
			"/snap/bin/docker",
			"/opt/docker/bin/docker",
		}
	}

	// Check each common path
	for _, path := range commonPaths {
		if fileExists(path) {
			dockerPath = path
			return dockerPath, nil
		}
	}

	// Last resort: try to find docker in some additional paths by expanding PATH
	pathEnv := os.Getenv("PATH")
	if runtime.GOOS == "darwin" {
		// On macOS, GUI apps might not have the same PATH as terminal
		// Add common paths that might be missing
		additionalPaths := []string{
			"/usr/local/bin",
			"/opt/homebrew/bin",
			"/usr/bin",
		}
		
		for _, additionalPath := range additionalPaths {
			if !strings.Contains(pathEnv, additionalPath) {
				pathEnv = pathEnv + string(os.PathListSeparator) + additionalPath
			}
		}
		
		// Temporarily set PATH and try again
		oldPath := os.Getenv("PATH")
		os.Setenv("PATH", pathEnv)
		if path, err := exec.LookPath("docker"); err == nil {
			dockerPath = path
			os.Setenv("PATH", oldPath) // Restore original PATH
			return dockerPath, nil
		}
		os.Setenv("PATH", oldPath) // Restore original PATH
	}

	// Try Windows-specific exe extension if not already tried
	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath("docker.exe"); err == nil {
			dockerPath = path
			return dockerPath, nil
		}
	}

	return "", &DockerNotFoundError{
		Message: "Docker executable not found. Please ensure Docker is installed and accessible.",
		Suggestions: []string{
			"Make sure Docker Desktop is installed and running",
			"Verify Docker is in your system PATH",
			"Try restarting the application after installing Docker",
		},
	}
}

// DockerNotFoundError represents an error when Docker executable cannot be found
type DockerNotFoundError struct {
	Message     string
	Suggestions []string
}

func (e *DockerNotFoundError) Error() string {
	result := e.Message
	if len(e.Suggestions) > 0 {
		result += "\n\nSuggestions:\n"
		for _, suggestion := range e.Suggestions {
			result += "• " + suggestion + "\n"
		}
	}
	return result
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// GetDockerCommand returns a command configured with the correct Docker path
func GetDockerCommand(args ...string) *exec.Cmd {
	dockerBinary, err := FindDockerPath()
	if err != nil {
		// Fallback to "docker" and let it fail with a more specific error
		dockerBinary = "docker"
	}
	
	cmd := exec.Command(dockerBinary, args...)
	// Apply platform-specific configuration (Windows console hiding, etc.)
	utils.SetupCommandForPlatform(cmd)
	return cmd
}

// ResetDockerPath clears the cached Docker path (useful for testing)
func ResetDockerPath() {
	dockerPath = ""
}