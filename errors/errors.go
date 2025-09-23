package errors

import (
	"errors"
	"fmt"
)

// Error types for different failure categories
var (
	// Docker-related errors
	ErrDockerNotAvailable   = errors.New("docker is not available")
	ErrDockerPermission     = errors.New("docker permission denied")
	ErrImageNotFound        = errors.New("docker image not found")
	ErrContainerNotFound    = errors.New("container not found")
	ErrContainerRunning     = errors.New("container is already running")
	ErrContainerNotRunning  = errors.New("container is not running")
	ErrPortConflict         = errors.New("port conflict detected")

	// File operation errors
	ErrFileNotFound         = errors.New("file not found")
	ErrFilePermission       = errors.New("file permission denied")
	ErrDirectoryNotFound    = errors.New("directory not found")
	ErrFileCorrupted        = errors.New("file is corrupted or invalid")
	ErrConfigInvalid        = errors.New("configuration file is invalid")

	// Validation errors
	ErrInvalidInput         = errors.New("invalid input provided")
	ErrMissingRequired      = errors.New("required field is missing")
	ErrInvalidFormat        = errors.New("invalid format")
	ErrInvalidContainerID   = errors.New("invalid container ID")
	ErrInvalidImageName     = errors.New("invalid image name")

	// Network errors
	ErrNetworkUnavailable   = errors.New("network is unavailable")
	ErrConnectionTimeout    = errors.New("connection timeout")
	ErrServiceUnavailable   = errors.New("service is unavailable")

	// Application state errors
	ErrAppNotInitialized    = errors.New("application not properly initialized")
	ErrOperationInProgress  = errors.New("operation already in progress")
	ErrInvalidState         = errors.New("invalid application state")
)

// Custom error types for enhanced context

// DockerError represents Docker-related errors with additional context
type DockerError struct {
	Operation   string // e.g., "pull", "run", "stop"
	ImageName   string
	ContainerID string
	Command     string
	Output      string
	Underlying  error
}

func (e *DockerError) Error() string {
	if e.ContainerID != "" {
		return fmt.Sprintf("docker %s failed for container %s: %v", e.Operation, e.ContainerID, e.Underlying)
	}
	if e.ImageName != "" {
		return fmt.Sprintf("docker %s failed for image %s: %v", e.Operation, e.ImageName, e.Underlying)
	}
	return fmt.Sprintf("docker %s failed: %v", e.Operation, e.Underlying)
}

func (e *DockerError) Unwrap() error {
	return e.Underlying
}

func (e *DockerError) WithOutput(output string) *DockerError {
	e.Output = output
	return e
}

// FileError represents file operation errors with path context
type FileError struct {
	Operation string // e.g., "read", "write", "delete", "create"
	Path      string
	Underlying error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file %s failed for %s: %v", e.Operation, e.Path, e.Underlying)
}

func (e *FileError) Unwrap() error {
	return e.Underlying
}

// ValidationError represents validation errors with field context
type ValidationError struct {
	Field   string
	Value   interface{}
	Reason  string
	Underlying error
}

func (e *ValidationError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("validation failed for field %s (value: %v): %s", e.Field, e.Value, e.Reason)
	}
	return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Reason)
}

func (e *ValidationError) Unwrap() error {
	return e.Underlying
}

// NetworkError represents network-related errors
type NetworkError struct {
	Operation string // e.g., "connect", "download", "health_check"
	URL       string
	Underlying error
}

func (e *NetworkError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("network %s failed for %s: %v", e.Operation, e.URL, e.Underlying)
	}
	return fmt.Sprintf("network %s failed: %v", e.Operation, e.Underlying)
}

func (e *NetworkError) Unwrap() error {
	return e.Underlying
}

// Error creation utilities

// NewDockerError creates a new DockerError with context
func NewDockerError(operation string, err error) *DockerError {
	return &DockerError{
		Operation:  operation,
		Underlying: err,
	}
}

// NewDockerErrorWithImage creates a DockerError with image context
func NewDockerErrorWithImage(operation, imageName string, err error) *DockerError {
	return &DockerError{
		Operation:  operation,
		ImageName:  imageName,
		Underlying: err,
	}
}

// NewDockerErrorWithContainer creates a DockerError with container context
func NewDockerErrorWithContainer(operation, containerID string, err error) *DockerError {
	return &DockerError{
		Operation:   operation,
		ContainerID: containerID,
		Underlying:  err,
	}
}

// NewFileError creates a new FileError with path context
func NewFileError(operation, path string, err error) *FileError {
	return &FileError{
		Operation:  operation,
		Path:       path,
		Underlying: err,
	}
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, reason string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:  field,
		Value:  value,
		Reason: reason,
	}
}

// NewValidationErrorWithCause creates a ValidationError with underlying cause
func NewValidationErrorWithCause(field, reason string, value interface{}, err error) *ValidationError {
	return &ValidationError{
		Field:      field,
		Value:      value,
		Reason:     reason,
		Underlying: err,
	}
}

// NewNetworkError creates a new NetworkError
func NewNetworkError(operation string, err error) *NetworkError {
	return &NetworkError{
		Operation:  operation,
		Underlying: err,
	}
}

// NewNetworkErrorWithURL creates a NetworkError with URL context
func NewNetworkErrorWithURL(operation, url string, err error) *NetworkError {
	return &NetworkError{
		Operation:  operation,
		URL:        url,
		Underlying: err,
	}
}

// Error wrapping utilities with enhanced context

// WrapWithContext wraps an error with additional context using fmt.Errorf with %w
func WrapWithContext(err error, context string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(context+": %w", append(args, err)...)
}

// WrapDockerError wraps a Docker operation error with consistent context
func WrapDockerError(operation, context string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("docker %s (%s): %w", operation, context, err)
}

// WrapFileError wraps a file operation error with path context
func WrapFileError(operation, path string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("file %s operation failed for %s: %w", operation, path, err)
}

// Error checking utilities

// IsDockerError checks if an error is Docker-related
func IsDockerError(err error) bool {
	var dockerErr *DockerError
	return errors.As(err, &dockerErr)
}

// IsFileError checks if an error is file-related
func IsFileError(err error) bool {
	var fileErr *FileError
	return errors.As(err, &fileErr)
}

// IsValidationError checks if an error is validation-related
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsNetworkError checks if an error is network-related
func IsNetworkError(err error) bool {
	var networkErr *NetworkError
	return errors.As(err, &networkErr)
}

// IsSpecificError checks if an error matches a specific error type
func IsSpecificError(err, target error) bool {
	return errors.Is(err, target)
}

// Error context extraction utilities

// GetDockerError extracts DockerError details if present
func GetDockerError(err error) (*DockerError, bool) {
	var dockerErr *DockerError
	if errors.As(err, &dockerErr) {
		return dockerErr, true
	}
	return nil, false
}

// GetFileError extracts FileError details if present
func GetFileError(err error) (*FileError, bool) {
	var fileErr *FileError
	if errors.As(err, &fileErr) {
		return fileErr, true
	}
	return nil, false
}

// GetValidationError extracts ValidationError details if present
func GetValidationError(err error) (*ValidationError, bool) {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return validationErr, true
	}
	return nil, false
}

// GetNetworkError extracts NetworkError details if present
func GetNetworkError(err error) (*NetworkError, bool) {
	var networkErr *NetworkError
	if errors.As(err, &networkErr) {
		return networkErr, true
	}
	return nil, false
}

// Validation utilities

// ValidateNotEmpty checks if a string is not empty
func ValidateNotEmpty(field, value string) error {
	if value == "" {
		return NewValidationError(field, "cannot be empty", value)
	}
	return nil
}

// ValidateContainerID validates a container ID format
func ValidateContainerID(containerID string) error {
	if err := ValidateNotEmpty("containerID", containerID); err != nil {
		return err
	}

	// Docker container IDs are typically 64-character hex strings, but can be shortened
	// Minimum practical length is 12 characters
	if len(containerID) < 12 {
		return NewValidationError("containerID", "too short (minimum 12 characters)", containerID)
	}

	return nil
}

// ValidateImageName validates a Docker image name format
func ValidateImageName(imageName string) error {
	if err := ValidateNotEmpty("imageName", imageName); err != nil {
		return err
	}

	// Basic validation - image names typically contain repository/image:tag format
	if len(imageName) < 3 {
		return NewValidationError("imageName", "too short", imageName)
	}

	return nil
}

// ValidateFilePath validates a file path
func ValidateFilePath(field, path string) error {
	if err := ValidateNotEmpty(field, path); err != nil {
		return err
	}

	// Additional path validation can be added here
	return nil
}

// Error aggregation utilities

// MultiError represents multiple errors
type MultiError struct {
	Errors []error
	Context string
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}

	if len(e.Errors) == 1 {
		if e.Context != "" {
			return fmt.Sprintf("%s: %v", e.Context, e.Errors[0])
		}
		return e.Errors[0].Error()
	}

	errorStr := fmt.Sprintf("multiple errors (%d)", len(e.Errors))
	if e.Context != "" {
		errorStr = fmt.Sprintf("%s - %s", e.Context, errorStr)
	}

	for i, err := range e.Errors {
		errorStr += fmt.Sprintf("\n  %d: %v", i+1, err)
	}

	return errorStr
}

func (e *MultiError) Unwrap() []error {
	return e.Errors
}

// NewMultiError creates a new MultiError
func NewMultiError(context string) *MultiError {
	return &MultiError{
		Context: context,
		Errors:  make([]error, 0),
	}
}

// Add adds an error to the MultiError
func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// HasErrors returns true if there are any errors
func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ToError returns the MultiError as an error if there are errors, nil otherwise
func (e *MultiError) ToError() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}