//go:build !windows
// +build !windows

package utils

import (
	"os/exec"
)

// SetupCommandForPlatform configures the command for the current platform
// On non-Windows platforms, this is a no-op
func SetupCommandForPlatform(cmd *exec.Cmd) {
	// No special setup needed for non-Windows platforms
}