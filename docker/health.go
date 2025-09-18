package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"moodle-prototype-manager/utils"
)

// HealthStatus represents the health check results
type HealthStatus struct {
	Docker   bool `json:"docker"`
	Internet bool `json:"internet"`
}

// CheckDockerHealth verifies Docker is installed and available
func CheckDockerHealth() bool {
	utils.LogDebug("Starting Docker health check...")
	
	// Log environment info for debugging
	pathEnv := os.Getenv("PATH")
	utils.LogDebug(fmt.Sprintf("Current PATH: %s", pathEnv))
	utils.LogDebug(fmt.Sprintf("Platform: %s", runtime.GOOS))
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Use our centralized Docker path detection
	dockerPath, err := FindDockerPath()
	if err != nil {
		utils.LogError("Docker path detection failed", err)
		utils.LogDebug("Docker may not be installed or not accessible from this application")
		return false
	}
	
	utils.LogDebug(fmt.Sprintf("Found Docker at: %s", dockerPath))
	
	// Test the Docker executable
	cmd := exec.CommandContext(ctx, dockerPath, "--version")
	utils.SetupCommandForPlatform(cmd)
	err = cmd.Run()
	
	if err != nil {
		utils.LogError(fmt.Sprintf("Docker health check failed using %s", dockerPath), err)
		return false
	}
	
	utils.LogDebug(fmt.Sprintf("Docker health check passed using: %s", dockerPath))
	return true
}

// CheckInternetHealth verifies internet connectivity using ping
func CheckInternetHealth() bool {
	utils.LogDebug("Starting Internet health check...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Try multiple methods to check internet connectivity
	targets := []string{"8.8.8.8", "1.1.1.1"} // Google DNS and Cloudflare DNS
	
	for _, target := range targets {
		if checkPingConnectivity(ctx, target) {
			utils.LogDebug(fmt.Sprintf("Internet health check passed using: %s", target))
			return true
		}
	}
	
	// On Windows, try alternative method if ping fails
	if runtime.GOOS == "windows" {
		utils.LogDebug("Trying alternative Windows connectivity check...")
		if checkWindowsConnectivity(ctx) {
			utils.LogDebug("Internet health check passed using Windows alternative method")
			return true
		}
	}
	
	utils.LogError("Internet health check failed - no connectivity to any target", nil)
	return false
}

// checkPingConnectivity tries to ping a specific target with platform-specific commands
func checkPingConnectivity(ctx context.Context, target string) bool {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "windows":
		// Windows ping: ping -n 1 -w 5000 target (5 second timeout in milliseconds)
		cmd = exec.CommandContext(ctx, "ping", "-n", "1", "-w", "5000", target)
	case "darwin":
		// macOS ping: ping -c 1 -W 5000 target (5 second timeout in milliseconds)
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", "5000", target)
	case "linux":
		// Linux ping: ping -c 1 -w 5 target (5 second timeout in seconds)
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-w", "5", target)
	default:
		// Fallback for other platforms
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", target)
	}
	
	utils.LogDebug(fmt.Sprintf("Trying ping: %s", strings.Join(cmd.Args, " ")))
	utils.SetupCommandForPlatform(cmd)
	err := cmd.Run()
	
	if err != nil {
		utils.LogDebug(fmt.Sprintf("Ping failed for %s: %v", target, err))
		return false
	}
	
	return true
}

// checkWindowsConnectivity tries alternative connectivity methods on Windows
func checkWindowsConnectivity(ctx context.Context) bool {
	// Try using nslookup as an alternative to ping on Windows
	targets := []string{"google.com", "cloudflare.com"}
	
	for _, target := range targets {
		utils.LogDebug(fmt.Sprintf("Trying nslookup: %s", target))
		cmd := exec.CommandContext(ctx, "nslookup", target)
		utils.SetupCommandForPlatform(cmd)
		err := cmd.Run()
		
		if err == nil {
			utils.LogDebug(fmt.Sprintf("nslookup successful for: %s", target))
			return true
		}
		utils.LogDebug(fmt.Sprintf("nslookup failed for %s: %v", target, err))
	}
	
	// Try using telnet as a last resort
	utils.LogDebug("Trying telnet connectivity check...")
	cmd := exec.CommandContext(ctx, "telnet", "8.8.8.8", "53")
	utils.SetupCommandForPlatform(cmd)
	err := cmd.Run()
	
	if err == nil {
		utils.LogDebug("Telnet connectivity check passed")
		return true
	}
	
	utils.LogDebug(fmt.Sprintf("Telnet connectivity check failed: %v", err))
	return false
}

// PerformHealthChecks runs all health checks
func PerformHealthChecks() HealthStatus {
	utils.LogInfo("=== Starting Health Checks ===")
	
	dockerHealth := CheckDockerHealth()
	internetHealth := CheckInternetHealth()
	
	status := HealthStatus{
		Docker:   dockerHealth,
		Internet: internetHealth,
	}
	
	utils.LogInfo(fmt.Sprintf("Health check results: Docker=%t, Internet=%t", dockerHealth, internetHealth))
	return status
}