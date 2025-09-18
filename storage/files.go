package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ContainerIDFile = "container.id"
	CredentialsFile = "moodle.txt"
	ImageConfigFile = "image.docker"
)

// FileManager handles file I/O operations
type FileManager struct{}

// NewFileManager creates a new file manager
func NewFileManager() *FileManager {
	return &FileManager{}
}

// getBaseDir returns the appropriate base directory for file operations
// In production, files should be next to the executable
// In development, files are in the working directory
func (fm *FileManager) getBaseDir() string {
	var baseDir string

	// First, check current working directory for go.mod (development detection)
	if wd, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			// Development environment - use working directory
			baseDir = wd
			fmt.Printf("[DEBUG] getBaseDir: Development mode detected (go.mod in working dir), using: %s\n", baseDir)
			return baseDir
		}
	}

	// Try to get executable directory (for production builds)
	if executable, err := os.Executable(); err == nil {
		execDir := filepath.Dir(executable)
		// Check if we're in a development environment by looking for go.mod in exec dir
		if _, err := os.Stat(filepath.Join(execDir, "go.mod")); err == nil {
			// Development environment - but go.mod is in exec dir, use that
			baseDir = execDir
			fmt.Printf("[DEBUG] getBaseDir: Development mode (go.mod in exec dir), using: %s\n", baseDir)
		} else {
			// Production environment - use executable directory
			baseDir = execDir
			fmt.Printf("[DEBUG] getBaseDir: Production mode, using exec dir: %s\n", baseDir)
		}
	} else {
		// Fallback to working directory
		if wd, err := os.Getwd(); err == nil {
			baseDir = wd
			fmt.Printf("[DEBUG] getBaseDir: Failed to get executable, using working dir: %s\n", baseDir)
		} else {
			// Last resort - current directory
			baseDir = "."
			fmt.Printf("[DEBUG] getBaseDir: All methods failed, using current dir: %s\n", baseDir)
		}
	}

	return baseDir
}

// getFilePath returns the full path for a given filename
func (fm *FileManager) getFilePath(filename string) string {
	return filepath.Join(fm.getBaseDir(), filename)
}

// SaveContainerID saves the container ID to file
func (fm *FileManager) SaveContainerID(containerID string) error {
	return os.WriteFile(fm.getFilePath(ContainerIDFile), []byte(containerID), 0644)
}

// LoadContainerID loads the container ID from file
func (fm *FileManager) LoadContainerID() (string, error) {
	data, err := os.ReadFile(fm.getFilePath(ContainerIDFile))
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(data)), nil
}

// ContainerIDExists checks if container ID file exists
func (fm *FileManager) ContainerIDExists() bool {
	_, err := os.Stat(fm.getFilePath(ContainerIDFile))
	return err == nil
}

// SaveCredentials saves credentials to file in key=value format
func (fm *FileManager) SaveCredentials(password, url string) error {
	content := fmt.Sprintf("password=%s\nurl=%s\n", password, url)
	return os.WriteFile(fm.getFilePath(CredentialsFile), []byte(content), 0644)
}

// LoadCredentials loads credentials from file
func (fm *FileManager) LoadCredentials() (map[string]string, error) {
	data, err := os.ReadFile(fm.getFilePath(CredentialsFile))
	if err != nil {
		return nil, err
	}
	
	credentials := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			credentials[key] = value
		}
	}
	
	return credentials, nil
}

// CredentialsExist checks if credentials file exists
func (fm *FileManager) CredentialsExist() bool {
	_, err := os.Stat(fm.getFilePath(CredentialsFile))
	return err == nil
}

// DeleteContainerID removes the container ID file
func (fm *FileManager) DeleteContainerID() error {
	if fm.ContainerIDExists() {
		return os.Remove(fm.getFilePath(ContainerIDFile))
	}
	return nil
}

// DeleteCredentials removes the credentials file
func (fm *FileManager) DeleteCredentials() error {
	if fm.CredentialsExist() {
		return os.Remove(fm.getFilePath(CredentialsFile))
	}
	return nil
}

// LoadImageName loads the Docker image name from configuration file
func (fm *FileManager) LoadImageName() (string, error) {
	// Try multiple potential paths for the image configuration file
	searchPaths := []string{
		fm.getFilePath(ImageConfigFile),                    // Primary path (working or exec dir)
		filepath.Join(".", ImageConfigFile),                // Current directory
		filepath.Join("..", ImageConfigFile),               // Parent directory
	}

	// If we can get working directory, also try that explicitly
	if wd, err := os.Getwd(); err == nil {
		searchPaths = append(searchPaths, filepath.Join(wd, ImageConfigFile))
	}

	// If we can get executable directory, also try that explicitly
	if executable, err := os.Executable(); err == nil {
		execDir := filepath.Dir(executable)
		searchPaths = append(searchPaths, filepath.Join(execDir, ImageConfigFile))
	}

	var lastErr error
	for i, imagePath := range searchPaths {
		fmt.Printf("[DEBUG] LoadImageName attempt %d: trying path: %s\n", i+1, imagePath)

		data, err := os.ReadFile(imagePath)
		if err != nil {
			fmt.Printf("[DEBUG] LoadImageName attempt %d: failed to read %s: %v\n", i+1, imagePath, err)
			lastErr = err
			continue
		}

		imageName := strings.TrimSpace(string(data))
		if imageName == "" {
			fmt.Printf("[DEBUG] LoadImageName attempt %d: file is empty at %s\n", i+1, imagePath)
			lastErr = fmt.Errorf("image configuration file is empty")
			continue
		}

		fmt.Printf("[DEBUG] LoadImageName: Successfully loaded image name '%s' from: %s\n", imageName, imagePath)
		return imageName, nil
	}

	// If we get here, all paths failed
	return "", fmt.Errorf("failed to find image configuration file in any of the searched paths (tried %d locations), last error: %w", len(searchPaths), lastErr)
}

// ImageConfigExists checks if image configuration file exists
func (fm *FileManager) ImageConfigExists() bool {
	_, err := os.Stat(fm.getFilePath(ImageConfigFile))
	return err == nil
}

// CleanupFiles removes all storage files
func (fm *FileManager) CleanupFiles() error {
	var errors []string
	
	if err := fm.DeleteContainerID(); err != nil {
		errors = append(errors, fmt.Sprintf("container ID: %v", err))
	}
	
	if err := fm.DeleteCredentials(); err != nil {
		errors = append(errors, fmt.Sprintf("credentials: %v", err))
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errors, ", "))
	}
	
	return nil
}