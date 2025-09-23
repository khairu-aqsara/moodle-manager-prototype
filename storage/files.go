package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moodle-prototype-manager/errors"
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
// Strategy:
// 1. Development mode (go.mod found): Use working/executable directory
// 2. Production mode: Use ~/.moodle-prototype-manager for easy team sharing
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

	// Try to get executable directory and check for go.mod there too
	if executable, err := os.Executable(); err == nil {
		execDir := filepath.Dir(executable)
		if _, err := os.Stat(filepath.Join(execDir, "go.mod")); err == nil {
			// Development environment - go.mod is in exec dir, use that
			baseDir = execDir
			fmt.Printf("[DEBUG] getBaseDir: Development mode (go.mod in exec dir), using: %s\n", baseDir)
			return baseDir
		}
	}

	// Production environment - use ~/.moodle-prototype-manager
	baseDir = fm.getUserDataDir()
	fmt.Printf("[DEBUG] getBaseDir: Production mode, using user data dir: %s\n", baseDir)

	// Ensure directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Printf("[ERROR] getBaseDir: Failed to create user data directory %s: %v\n", baseDir, err)
		// Fallback to working directory
		if wd, err := os.Getwd(); err == nil {
			baseDir = wd
			fmt.Printf("[DEBUG] getBaseDir: Falling back to working dir: %s\n", baseDir)
		} else {
			// Last resort - current directory
			baseDir = "."
			fmt.Printf("[DEBUG] getBaseDir: All methods failed, using current dir: %s\n", baseDir)
		}
	}

	return baseDir
}

// getUserDataDir returns the user data directory
// Uses ~/.moodle-prototype-manager for all platforms for consistency and easy team sharing
func (fm *FileManager) getUserDataDir() string {
	appName := ".moodle-prototype-manager"

	// Use ~/.moodle-prototype-manager for all platforms
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, appName)
	}

	// Fallback to current directory if home directory cannot be determined
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}

// getFilePath returns the full path for a given filename
func (fm *FileManager) getFilePath(filename string) string {
	return filepath.Join(fm.getBaseDir(), filename)
}

// ensureDirectoryExists creates the directory if it doesn't exist
func (fm *FileManager) ensureDirectoryExists(dirPath string) error {
	if dirPath == "" {
		return errors.NewValidationError("dirPath", "directory path cannot be empty", dirPath)
	}

	// Check if directory already exists
	if info, err := os.Stat(dirPath); err == nil {
		if !info.IsDir() {
			return errors.NewFileError("validate", dirPath, fmt.Errorf("path exists but is not a directory"))
		}
		return nil // Directory already exists
	}

	// Create directory with proper permissions
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return errors.NewFileError("create", dirPath, err)
	}

	fmt.Printf("[DEBUG] ensureDirectoryExists: Created directory %s\n", dirPath)
	return nil
}

// SaveContainerID saves the container ID to file
func (fm *FileManager) SaveContainerID(containerID string) error {
	// Validate input
	if err := errors.ValidateContainerID(containerID); err != nil {
		return errors.WrapWithContext(err, "invalid container ID provided to SaveContainerID")
	}

	filePath := fm.getFilePath(ContainerIDFile)
	fmt.Printf("[DEBUG] SaveContainerID: Writing to %s\n", filePath)

	// Ensure directory exists
	if err := fm.ensureDirectoryExists(filepath.Dir(filePath)); err != nil {
		return errors.WrapWithContext(err, "failed to ensure directory exists for container ID file")
	}

	err := os.WriteFile(filePath, []byte(containerID), 0644)
	if err != nil {
		fmt.Printf("[ERROR] SaveContainerID: Failed to write to %s: %v\n", filePath, err)
		return errors.NewFileError("write", filePath, err)
	}

	fmt.Printf("[DEBUG] SaveContainerID: Successfully wrote container ID to %s\n", filePath)
	return nil
}

// LoadContainerID loads the container ID from file
func (fm *FileManager) LoadContainerID() (string, error) {
	filePath := fm.getFilePath(ContainerIDFile)
	fmt.Printf("[DEBUG] LoadContainerID: Reading from %s\n", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[ERROR] LoadContainerID: Failed to read from %s: %v\n", filePath, err)
		return "", errors.NewFileError("read", filePath, err)
	}

	containerID := strings.TrimSpace(string(data))
	if containerID == "" {
		return "", errors.NewFileError("parse", filePath, errors.ErrFileCorrupted)
	}

	// Validate the loaded container ID
	if err := errors.ValidateContainerID(containerID); err != nil {
		return "", errors.WrapWithContext(err, "loaded container ID from file %s is invalid", filePath)
	}

	fmt.Printf("[DEBUG] LoadContainerID: Successfully loaded container ID from %s\n", filePath)
	return containerID, nil
}

// ContainerIDExists checks if container ID file exists
func (fm *FileManager) ContainerIDExists() bool {
	_, err := os.Stat(fm.getFilePath(ContainerIDFile))
	return err == nil
}

// SaveCredentials saves credentials to file in key=value format
func (fm *FileManager) SaveCredentials(password, url string) error {
	// Validate input
	if err := errors.ValidateNotEmpty("password", password); err != nil {
		return errors.WrapWithContext(err, "invalid password provided to SaveCredentials")
	}
	if err := errors.ValidateNotEmpty("url", url); err != nil {
		return errors.WrapWithContext(err, "invalid URL provided to SaveCredentials")
	}

	filePath := fm.getFilePath(CredentialsFile)
	fmt.Printf("[DEBUG] SaveCredentials: Writing to %s\n", filePath)

	// Ensure directory exists
	if err := fm.ensureDirectoryExists(filepath.Dir(filePath)); err != nil {
		return errors.WrapWithContext(err, "failed to ensure directory exists for credentials file")
	}

	content := fmt.Sprintf("password=%s\nurl=%s\n", password, url)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("[ERROR] SaveCredentials: Failed to write to %s: %v\n", filePath, err)
		return errors.NewFileError("write", filePath, err)
	}

	fmt.Printf("[DEBUG] SaveCredentials: Successfully wrote credentials to %s\n", filePath)
	return nil
}

// LoadCredentials loads credentials from file
func (fm *FileManager) LoadCredentials() (map[string]string, error) {
	filePath := fm.getFilePath(CredentialsFile)
	fmt.Printf("[DEBUG] LoadCredentials: Reading from %s\n", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[ERROR] LoadCredentials: Failed to read from %s: %v\n", filePath, err)
		return nil, errors.NewFileError("read", filePath, err)
	}

	if len(data) == 0 {
		return nil, errors.NewFileError("parse", filePath, errors.ErrFileCorrupted)
	}

	credentials := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	validLineCount := 0

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("[WARNING] LoadCredentials: Skipping malformed line %d in %s: %s\n", lineNum+1, filePath, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			fmt.Printf("[WARNING] LoadCredentials: Skipping line %d with empty key in %s\n", lineNum+1, filePath)
			continue
		}

		credentials[key] = value
		validLineCount++
	}

	if validLineCount == 0 {
		return nil, errors.WrapWithContext(errors.ErrFileCorrupted, "no valid credential entries found in file %s", filePath)
	}

	fmt.Printf("[DEBUG] LoadCredentials: Successfully loaded %d credential entries from %s\n", validLineCount, filePath)
	return credentials, nil
}

// CredentialsExist checks if credentials file exists
func (fm *FileManager) CredentialsExist() bool {
	_, err := os.Stat(fm.getFilePath(CredentialsFile))
	return err == nil
}

// DeleteContainerID removes the container ID file
func (fm *FileManager) DeleteContainerID() error {
	filePath := fm.getFilePath(ContainerIDFile)
	if !fm.ContainerIDExists() {
		// File doesn't exist, nothing to delete
		return nil
	}

	err := os.Remove(filePath)
	if err != nil {
		return errors.NewFileError("delete", filePath, err)
	}
	return nil
}

// DeleteCredentials removes the credentials file
func (fm *FileManager) DeleteCredentials() error {
	filePath := fm.getFilePath(CredentialsFile)
	if !fm.CredentialsExist() {
		// File doesn't exist, nothing to delete
		return nil
	}

	err := os.Remove(filePath)
	if err != nil {
		return errors.NewFileError("delete", filePath, err)
	}
	return nil
}

// LoadImageName loads the Docker image name from configuration file
func (fm *FileManager) LoadImageName() (string, error) {
	// Try multiple potential paths for the image configuration file
	searchPaths := []string{
		fm.getFilePath(ImageConfigFile),      // Primary path (working or exec dir)
		filepath.Join(".", ImageConfigFile),  // Current directory
		filepath.Join("..", ImageConfigFile), // Parent directory
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
			lastErr = errors.NewFileError("parse", imagePath, errors.ErrConfigInvalid)
			continue
		}

		// Validate image name format
		if err := errors.ValidateImageName(imageName); err != nil {
			fmt.Printf("[DEBUG] LoadImageName attempt %d: invalid image name in %s: %v\n", i+1, imagePath, err)
			lastErr = errors.WrapWithContext(err, "image name in file %s is invalid", imagePath)
			continue
		}

		fmt.Printf("[DEBUG] LoadImageName: Successfully loaded image name '%s' from: %s\n", imageName, imagePath)
		return imageName, nil
	}

	// If we get here, all paths failed
	return "", errors.WrapWithContext(lastErr, "failed to find image configuration file in any of %d searched paths", len(searchPaths))
}

// ImageConfigExists checks if image configuration file exists
func (fm *FileManager) ImageConfigExists() bool {
	_, err := os.Stat(fm.getFilePath(ImageConfigFile))
	return err == nil
}

// CleanupFiles removes all storage files
func (fm *FileManager) CleanupFiles() error {
	multiErr := errors.NewMultiError("file cleanup operation")

	if err := fm.DeleteContainerID(); err != nil {
		multiErr.Add(errors.WrapWithContext(err, "failed to delete container ID file"))
	}

	if err := fm.DeleteCredentials(); err != nil {
		multiErr.Add(errors.WrapWithContext(err, "failed to delete credentials file"))
	}

	return multiErr.ToError()
}
