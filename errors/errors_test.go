package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestCustomErrorTypes(t *testing.T) {
	t.Run("DockerError", func(t *testing.T) {
		baseErr := fmt.Errorf("command failed")
		dockerErr := NewDockerErrorWithContainer("stop", "container123", baseErr)

		// Test error message
		expected := "docker stop failed for container container123"
		if !strings.Contains(dockerErr.Error(), expected) {
			t.Errorf("Expected error to contain '%s', got: %s", expected, dockerErr.Error())
		}

		// Test unwrapping
		if !errors.Is(dockerErr, baseErr) {
			t.Errorf("Expected DockerError to wrap the base error")
		}

		// Test output addition
		dockerErr.WithOutput("exit code 1")
		if dockerErr.Output != "exit code 1" {
			t.Errorf("Expected output to be set, got: %s", dockerErr.Output)
		}
	})

	t.Run("FileError", func(t *testing.T) {
		baseErr := fmt.Errorf("permission denied")
		fileErr := NewFileError("read", "/tmp/test.txt", baseErr)

		expected := "file read failed for /tmp/test.txt"
		if !strings.Contains(fileErr.Error(), expected) {
			t.Errorf("Expected error to contain '%s', got: %s", expected, fileErr.Error())
		}

		if !errors.Is(fileErr, baseErr) {
			t.Errorf("Expected FileError to wrap the base error")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		validationErr := NewValidationError("containerID", "too short", "abc")

		expected := "validation failed for field containerID"
		if !strings.Contains(validationErr.Error(), expected) {
			t.Errorf("Expected error to contain '%s', got: %s", expected, validationErr.Error())
		}

		if !strings.Contains(validationErr.Error(), "abc") {
			t.Errorf("Expected error to contain the value 'abc'")
		}
	})

	t.Run("NetworkError", func(t *testing.T) {
		baseErr := fmt.Errorf("connection timeout")
		networkErr := NewNetworkErrorWithURL("connect", "http://localhost:8080", baseErr)

		expected := "network connect failed for http://localhost:8080"
		if !strings.Contains(networkErr.Error(), expected) {
			t.Errorf("Expected error to contain '%s', got: %s", expected, networkErr.Error())
		}
	})
}

func TestErrorWrapping(t *testing.T) {
	t.Run("WrapWithContext", func(t *testing.T) {
		baseErr := fmt.Errorf("original error")
		wrappedErr := WrapWithContext(baseErr, "operation failed for %s", "test")

		expected := "operation failed for test"
		if !strings.Contains(wrappedErr.Error(), expected) {
			t.Errorf("Expected wrapped error to contain '%s', got: %s", expected, wrappedErr.Error())
		}

		if !errors.Is(wrappedErr, baseErr) {
			t.Errorf("Expected wrapped error to preserve the original error")
		}
	})

	t.Run("WrapWithNilError", func(t *testing.T) {
		result := WrapWithContext(nil, "operation failed")
		if result != nil {
			t.Errorf("Expected nil when wrapping nil error, got: %v", result)
		}
	})
}

func TestErrorChecking(t *testing.T) {
	dockerErr := NewDockerError("test", fmt.Errorf("base"))
	fileErr := NewFileError("read", "/tmp", fmt.Errorf("base"))
	validationErr := NewValidationError("field", "reason", "value")
	networkErr := NewNetworkError("connect", fmt.Errorf("base"))

	tests := []struct {
		name     string
		err      error
		isDocker bool
		isFile   bool
		isValid  bool
		isNet    bool
	}{
		{"DockerError", dockerErr, true, false, false, false},
		{"FileError", fileErr, false, true, false, false},
		{"ValidationError", validationErr, false, false, true, false},
		{"NetworkError", networkErr, false, false, false, true},
		{"StandardError", fmt.Errorf("standard"), false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsDockerError(tt.err) != tt.isDocker {
				t.Errorf("IsDockerError() = %v, want %v", IsDockerError(tt.err), tt.isDocker)
			}
			if IsFileError(tt.err) != tt.isFile {
				t.Errorf("IsFileError() = %v, want %v", IsFileError(tt.err), tt.isFile)
			}
			if IsValidationError(tt.err) != tt.isValid {
				t.Errorf("IsValidationError() = %v, want %v", IsValidationError(tt.err), tt.isValid)
			}
			if IsNetworkError(tt.err) != tt.isNet {
				t.Errorf("IsNetworkError() = %v, want %v", IsNetworkError(tt.err), tt.isNet)
			}
		})
	}
}

func TestValidationFunctions(t *testing.T) {
	t.Run("ValidateNotEmpty", func(t *testing.T) {
		tests := []struct {
			field    string
			value    string
			hasError bool
		}{
			{"test", "valid", false},
			{"test", "", true},
			{"test", "  ", false}, // spaces are considered valid
		}

		for _, tt := range tests {
			err := ValidateNotEmpty(tt.field, tt.value)
			if (err != nil) != tt.hasError {
				t.Errorf("ValidateNotEmpty(%q, %q) error = %v, wantError = %v",
					tt.field, tt.value, err, tt.hasError)
			}
		}
	})

	t.Run("ValidateContainerID", func(t *testing.T) {
		tests := []struct {
			containerID string
			hasError    bool
		}{
			{"abcdef123456", false},
			{"abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", false},
			{"", true},
			{"short", true},
			{"abc", true},
		}

		for _, tt := range tests {
			err := ValidateContainerID(tt.containerID)
			if (err != nil) != tt.hasError {
				t.Errorf("ValidateContainerID(%q) error = %v, wantError = %v",
					tt.containerID, err, tt.hasError)
			}
		}
	})

	t.Run("ValidateImageName", func(t *testing.T) {
		tests := []struct {
			imageName string
			hasError  bool
		}{
			{"nginx:latest", false},
			{"registry.example.com/user/image:tag", false},
			{"", true},
			{"ab", true},
		}

		for _, tt := range tests {
			err := ValidateImageName(tt.imageName)
			if (err != nil) != tt.hasError {
				t.Errorf("ValidateImageName(%q) error = %v, wantError = %v",
					tt.imageName, err, tt.hasError)
			}
		}
	})
}

func TestMultiError(t *testing.T) {
	t.Run("EmptyMultiError", func(t *testing.T) {
		multiErr := NewMultiError("test")
		if multiErr.HasErrors() {
			t.Errorf("Expected empty MultiError to have no errors")
		}
		if multiErr.ToError() != nil {
			t.Errorf("Expected empty MultiError.ToError() to return nil")
		}
	})

	t.Run("MultiErrorWithErrors", func(t *testing.T) {
		multiErr := NewMultiError("operation failed")
		multiErr.Add(fmt.Errorf("error 1"))
		multiErr.Add(fmt.Errorf("error 2"))
		multiErr.Add(nil) // Should be ignored

		if !multiErr.HasErrors() {
			t.Errorf("Expected MultiError to have errors")
		}

		if len(multiErr.Errors) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(multiErr.Errors))
		}

		err := multiErr.ToError()
		if err == nil {
			t.Errorf("Expected ToError() to return error")
		}

		errorStr := err.Error()
		if !strings.Contains(errorStr, "operation failed") {
			t.Errorf("Expected error string to contain context")
		}
		if !strings.Contains(errorStr, "error 1") {
			t.Errorf("Expected error string to contain first error")
		}
		if !strings.Contains(errorStr, "error 2") {
			t.Errorf("Expected error string to contain second error")
		}
	})

	t.Run("SingleMultiError", func(t *testing.T) {
		multiErr := NewMultiError("test")
		multiErr.Add(fmt.Errorf("single error"))

		errorStr := multiErr.Error()
		if !strings.Contains(errorStr, "test: single error") {
			t.Errorf("Expected single error format, got: %s", errorStr)
		}
	})
}

func TestSpecificErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		match  bool
	}{
		{"DockerNotAvailable", ErrDockerNotAvailable, ErrDockerNotAvailable, true},
		{"ImageNotFound", ErrImageNotFound, ErrImageNotFound, true},
		{"Different errors", ErrDockerNotAvailable, ErrImageNotFound, false},
		{"Wrapped error", fmt.Errorf("wrapped: %w", ErrDockerNotAvailable), ErrDockerNotAvailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSpecificError(tt.err, tt.target)
			if result != tt.match {
				t.Errorf("IsSpecificError() = %v, want %v", result, tt.match)
			}
		})
	}
}