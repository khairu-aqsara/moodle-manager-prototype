//go:build windows
// +build windows

package utils

import (
	"os/exec"
	"syscall"
)

// SetupCommandForPlatform configures the command for the current platform
// On Windows, this hides the console window to prevent flashing
func SetupCommandForPlatform(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}