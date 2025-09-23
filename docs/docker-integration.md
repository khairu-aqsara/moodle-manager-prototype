# Moodle Prototype Manager - Docker Integration Documentation

## Table of Contents

1. [Overview](#overview)
2. [Docker Architecture](#docker-architecture)
3. [Container Lifecycle Management](#container-lifecycle-management)
4. [Image Management](#image-management)
5. [Health Check System](#health-check-system)
6. [Log Processing](#log-processing)
7. [Progress Tracking](#progress-tracking)
8. [Cross-Platform Considerations](#cross-platform-considerations)
9. [Error Handling](#error-handling)
10. [Security Considerations](#security-considerations)
11. [Performance Optimization](#performance-optimization)

## Overview

The Moodle Prototype Manager application provides a sophisticated Docker integration layer that manages the complete lifecycle of Moodle prototype containers. The integration is designed to handle cross-platform Docker operations, real-time progress tracking, and robust error handling.

### Key Features

- **Container Lifecycle Management**: Create, start, stop, and monitor containers
- **Image Management**: Pull images with progress tracking
- **Health Monitoring**: Docker daemon and container status checks
- **Log Processing**: Extract credentials and monitor container startup
- **Cross-Platform Support**: Windows, macOS, and Linux compatibility
- **Error Recovery**: Graceful handling of Docker operation failures

### Target Container

The application is specifically designed to work with the `wenkhairu/moodle-prototype:502-stable` Docker image, though the architecture supports configurable image names.

## Docker Architecture

### Component Structure

```
Docker Integration Layer
├── docker.Manager          # Core Docker operations
├── docker.HealthChecker     # Health check operations
├── docker.LogParser         # Container log processing
├── docker.ProgressTracker   # Pull progress monitoring
└── docker.PathManager      # Cross-platform command handling
```

### Manager Pattern

The `docker.Manager` struct serves as the central orchestrator for all Docker operations:

```go
type Manager struct {
    imageName string  // Configurable Docker image name
}
```

**Key Responsibilities:**
1. Image existence verification and pulling
2. Container creation and lifecycle management
3. Container status monitoring and health checks
4. Log retrieval and processing coordination

### Command Execution Strategy

**Cross-Platform Command Generation:**
```go
func GetDockerCommand(args ...string) *exec.Cmd {
    cmd := exec.Command("docker", args...)
    utils.SetupCommandForPlatform(cmd)
    return cmd
}
```

**Platform-Specific Setup:**
- **Windows**: Configures process creation flags to hide console windows
- **Unix-like systems**: Uses default process execution

## Container Lifecycle Management

### Container States

The application recognizes and manages several container states:

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   No Image  │───▶│ Image Ready │───▶│  Container  │───▶│ Container   │
│             │    │             │    │  Created    │    │  Running    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
       │                  │                  │                  │
       ▼                  ▼                  ▼                  ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Pull Image  │    │ Run Command │    │ Start/Stop  │    │ Log Monitor │
│ (Progress)  │    │ Execution   │    │ Operations  │    │ & Extract   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### First-Time Container Creation

**Process Flow:**
1. **Image Verification**: Check if target image exists locally
2. **Image Pull**: Download image with progress tracking if needed
3. **Container Creation**: Run new container with port mapping
4. **ID Persistence**: Save container ID for future operations
5. **Startup Monitoring**: Wait for container readiness
6. **Credential Extraction**: Parse logs for admin credentials

**Implementation:**
```go
func (m *Manager) RunContainer() (string, error) {
    // Validate image name
    if err := errors.ValidateImageName(m.imageName); err != nil {
        return "", err
    }

    // Execute docker run command
    cmd := GetDockerCommand("run", "-d", "-p", ContainerPort, m.imageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", errors.NewDockerErrorWithImage("run", m.imageName, err)
    }

    // Extract and validate container ID
    containerID := strings.TrimSpace(string(output))
    if err := errors.ValidateContainerID(containerID); err != nil {
        return "", err
    }

    return containerID, nil
}
```

### Container Restart Operations

**Process Flow:**
1. **ID Validation**: Verify container ID exists and is valid
2. **Status Check**: Determine current container state
3. **Start Operation**: Execute docker start command
4. **Health Check**: Verify container is responding
5. **Credential Update**: Refresh stored credentials if needed

**Implementation:**
```go
func (m *Manager) StartContainer(containerID string) error {
    // Validate container ID format
    if err := errors.ValidateContainerID(containerID); err != nil {
        return err
    }

    // Execute start command
    cmd := GetDockerCommand("start", containerID)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return errors.NewDockerErrorWithContainer("start", containerID, err).
               WithOutput(string(output))
    }

    return nil
}
```

### Container Shutdown Process

**Graceful Stop Strategy:**
1. **Status Verification**: Confirm container is running
2. **Graceful Stop**: Execute `docker stop` command
3. **Fallback Force Stop**: Use `docker kill` if graceful stop fails
4. **Cleanup**: Remove temporary files if needed

**Implementation:**
```go
func (m *Manager) StopContainer(containerID string) error {
    cmd := GetDockerCommand("stop", containerID)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return errors.NewDockerErrorWithContainer("stop", containerID, err).
               WithOutput(string(output))
    }
    return nil
}

func (m *Manager) ForceStopContainer(containerID string) error {
    cmd := GetDockerCommand("kill", containerID)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return errors.NewDockerErrorWithContainer("kill", containerID, err).
               WithOutput(string(output))
    }
    return nil
}
```

## Image Management

### Image Discovery Process

**Local Image Check:**
```go
func (m *Manager) CheckImageExists() (bool, error) {
    cmd := GetDockerCommand("images", "--format", "{{.Repository}}:{{.Tag}}")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return false, errors.NewDockerErrorWithImage("check", m.imageName, err)
    }

    exists := strings.Contains(string(output), m.imageName)
    return exists, nil
}
```

**Command Execution:**
```bash
docker images --format "{{.Repository}}:{{.Tag}}"
```

**Output Processing:**
- Searches for exact image name match in output
- Handles multi-line responses with different image entries
- Returns boolean indicating image presence

### Image Pull Operations

#### Simple Pull

**Basic image download without progress tracking:**
```go
func (m *Manager) PullImage() error {
    cmd := GetDockerCommand("pull", m.imageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return errors.NewDockerErrorWithImage("pull", m.imageName, err).
               WithOutput(string(output))
    }
    return nil
}
```

#### Advanced Pull with Progress

**Progress-tracked image download:**
```go
func (m *Manager) PullImageWithProgress(progressCallback func(float64, string)) error {
    cmd := GetDockerCommand("pull", m.imageName)

    // Create pipes for output processing
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return errors.NewDockerErrorWithImage("pull_setup", m.imageName, err)
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return errors.NewDockerErrorWithImage("pull_setup", m.imageName, err)
    }

    // Start command execution
    if err := cmd.Start(); err != nil {
        return errors.NewDockerErrorWithImage("pull_start", m.imageName, err)
    }

    // Process streams concurrently
    progress := NewPullProgress()
    progress.AddCallback(progressCallback)

    errChan := make(chan error, 2)

    go func() {
        errChan <- progress.ProcessStream(stdout)
    }()

    go func() {
        errChan <- progress.ProcessStream(stderr)
    }()

    // Wait for completion
    cmdErr := cmd.Wait()
    <-errChan  // Wait for stdout processing
    <-errChan  // Wait for stderr processing

    return cmdErr
}
```

## Health Check System

### Docker Daemon Health

**Health Check Process:**
```go
func CheckDockerHealth() bool {
    cmd := GetDockerCommand("--version")

    // Set short timeout for responsiveness
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    cmd = cmd.WithContext(ctx)

    err := cmd.Run()
    return err == nil
}
```

**What This Checks:**
- Docker CLI availability
- Docker daemon accessibility
- Basic permission levels
- Command execution capability

### Internet Connectivity

**Network Check Implementation:**
```go
func CheckInternetConnectivity() bool {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get("http://www.google.com")
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode >= 200 && resp.StatusCode < 400
}
```

**Features:**
- Short timeout for responsiveness
- Uses reliable target (Google)
- Accepts any successful HTTP response
- Handles network timeouts gracefully

### Container Health Monitoring

**Container Status Check:**
```go
func (m *Manager) IsContainerRunning(containerID string) (bool, error) {
    cmd := GetDockerCommand("inspect", "--format={{.State.Running}}", containerID)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return false, errors.NewDockerErrorWithContainer("inspect", containerID, err)
    }

    return strings.TrimSpace(string(output)) == "true", nil
}
```

**HTTP Health Check:**
```go
func testMoodleHTTP() bool {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get("http://localhost:8080")
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    // Any HTTP response indicates server is running
    return resp.StatusCode > 0
}
```

## Log Processing

### Container Log Retrieval

**Standard Log Retrieval:**
```go
func (m *Manager) GetContainerLogs(containerID string) (string, error) {
    cmd := GetDockerCommand("logs", containerID)

    // Use CombinedOutput to capture both stdout and stderr
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", errors.NewDockerErrorWithContainer("logs", containerID, err)
    }

    return string(output), nil
}
```

**Time-Based Log Retrieval:**
```go
func (m *Manager) GetContainerLogsSince(containerID string, since time.Time) (string, error) {
    sinceStr := since.Format(time.RFC3339)
    cmd := GetDockerCommand("logs", "--since", sinceStr, containerID)

    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", errors.NewDockerErrorWithContainer("logs_since", containerID, err)
    }

    return string(output), nil
}
```

### Credential Extraction

**Log Parsing Patterns:**
```go
type LogParser struct {
    passwordRegex *regexp.Regexp
    urlRegex      *regexp.Regexp
}

func NewLogParser() *LogParser {
    return &LogParser{
        passwordRegex: regexp.MustCompile(`Generated admin password: (.+)`),
        urlRegex:      regexp.MustCompile(`Moodle is available at: (http://[^\s]+)`),
    }
}
```

**Pattern Matching:**
```go
func (lp *LogParser) ExtractCredentials(logs string) Credentials {
    var creds Credentials

    // Extract admin password
    if matches := lp.passwordRegex.FindStringSubmatch(logs); len(matches) > 1 {
        creds.Password = strings.TrimSpace(matches[1])
    }

    // Extract Moodle URL
    if matches := lp.urlRegex.FindStringSubmatch(logs); len(matches) > 1 {
        creds.URL = strings.TrimSpace(matches[1])
    }

    return creds
}
```

**Target Log Patterns:**
- **Admin Password**: `Generated admin password: <password_value>`
- **Moodle URL**: `Moodle is available at: http://localhost:8080`

### Polling Strategy

**Continuous Log Monitoring:**
```go
func (a *App) waitForContainerAndExtractCredentialsSince(containerID string, since time.Time) {
    logErrorCount := 0
    maxLogErrors := 5

    for {
        logs, err := a.dockerManager.GetContainerLogs(containerID)
        if err != nil {
            logErrorCount++
            if logErrorCount > maxLogErrors {
                time.Sleep(5 * time.Second)  // Longer sleep after errors
            } else {
                time.Sleep(2 * time.Second)  // Normal poll interval
            }
            continue
        }

        logErrorCount = 0  // Reset on success

        creds := a.logParser.ExtractCredentials(logs)
        if creds.IsComplete() {
            // Save credentials and exit
            return
        }

        time.Sleep(2 * time.Second)
    }
}
```

**Features:**
- Exponential backoff on errors
- Error count limiting to prevent spam
- Continuous polling until credentials found
- No timeout for first-time installations (can take 20+ minutes on Windows)

## Progress Tracking

### Pull Progress Architecture

**Progress Tracker Structure:**
```go
type PullProgress struct {
    callbacks []func(float64, string)
    layers    map[string]*LayerProgress
    mutex     sync.Mutex
}

type LayerProgress struct {
    ID       string
    Status   string
    Current  int64
    Total    int64
}
```

### Docker Output Processing

**Stream Processing:**
```go
func (pp *PullProgress) ProcessStream(reader io.Reader) error {
    scanner := bufio.NewScanner(reader)

    for scanner.Scan() {
        line := scanner.Text()

        // Parse Docker JSON output
        var dockerMsg DockerMessage
        if err := json.Unmarshal([]byte(line), &dockerMsg); err != nil {
            continue  // Skip non-JSON lines
        }

        pp.processDockerMessage(dockerMsg)
    }

    return scanner.Err()
}
```

**Progress Calculation:**
```go
func (pp *PullProgress) calculateOverallProgress() (float64, string) {
    pp.mutex.Lock()
    defer pp.mutex.Unlock()

    if len(pp.layers) == 0 {
        return 0.0, "Starting download..."
    }

    var totalCurrent, totalSize int64
    var status string

    for _, layer := range pp.layers {
        totalCurrent += layer.Current
        totalSize += layer.Total

        // Use most recent status message
        if layer.Status != "" {
            status = layer.Status
        }
    }

    if totalSize == 0 {
        return 0.0, status
    }

    percentage := (float64(totalCurrent) / float64(totalSize)) * 100
    return percentage, status
}
```

### Progress Event Emission

**Frontend Integration:**
```go
func (a *App) RunMoodle() error {
    // ... image pull with progress
    err := a.dockerManager.PullImageWithProgress(func(percentage float64, status string) {
        progressData := map[string]any{
            "percentage": percentage,
            "status":     status,
        }
        wailsruntime.EventsEmit(a.ctx, "docker:pull:progress", progressData)
    })
}
```

**Frontend Handling:**
```javascript
// Listen for progress events
window.wails.Events.On('docker:pull:progress', (data) => {
    const { percentage, status } = data;
    updateDownloadProgress(percentage, status);
});
```

## Cross-Platform Considerations

### Command Execution Differences

**Windows-Specific Handling:**
```go
// utils/platform_windows.go
func SetupCommandForPlatform(cmd *exec.Cmd) {
    cmd.SysProcAttr = &syscall.SysProcAttr{
        HideWindow:    true,
        CreationFlags: CREATE_NO_WINDOW,
    }
}
```

**Unix-Like Systems:**
```go
// utils/platform_other.go
func SetupCommandForPlatform(cmd *exec.Cmd) {
    // No special configuration needed
}
```

### Path and File Handling

**Directory Management:**
- **Development**: Uses project directory (detected by `go.mod`)
- **Production**: Uses user data directories (`~/.moodle-prototype-manager`)

**Cross-Platform Paths:**
```go
func (fm *FileManager) getUserDataDir() string {
    appName := ".moodle-prototype-manager"

    if home, err := os.UserHomeDir(); err == nil {
        return filepath.Join(home, appName)
    }

    // Fallback to current directory
    if wd, err := os.Getwd(); err == nil {
        return wd
    }
    return "."
}
```

### Docker Daemon Differences

**Platform-Specific Considerations:**
- **Windows**: Docker Desktop must be running, may require administrator privileges
- **macOS**: Docker Desktop provides Docker daemon, socket permissions vary
- **Linux**: Native Docker daemon, user group membership important

## Error Handling

### Docker Error Classification

**Error Types and Context:**
```go
type DockerError struct {
    Operation string      // What operation failed (pull, run, start, etc.)
    Image     string      // Target image (if applicable)
    Container string      // Container ID (if applicable)
    Output    string      // Docker command output
    Cause     error       // Underlying error
}
```

**Error Creation Patterns:**
```go
// Image-specific errors
func NewDockerErrorWithImage(op, image string, err error) *DockerError {
    return &DockerError{
        Operation: op,
        Image:     image,
        Cause:     err,
    }
}

// Container-specific errors
func NewDockerErrorWithContainer(op, containerID string, err error) *DockerError {
    return &DockerError{
        Operation: op,
        Container: containerID,
        Cause:     err,
    }
}
```

### Error Recovery Strategies

**Graceful Degradation:**
1. **Image Pull Failures**: Retry with exponential backoff
2. **Container Start Failures**: Fall back to new container creation
3. **Health Check Failures**: Continue with warnings
4. **Log Parsing Failures**: Continue polling with error limiting

**Error Context Preservation:**
```go
func (m *Manager) StopContainer(containerID string) error {
    err := m.gracefulStop(containerID)
    if err != nil {
        // Try force stop as fallback
        forceErr := m.ForceStopContainer(containerID)
        if forceErr != nil {
            return fmt.Errorf("failed to stop container (graceful: %v, force: %v)",
                             err, forceErr)
        }
    }
    return nil
}
```

### User-Friendly Error Messages

**Error Translation:**
```go
func TranslateDockerError(err error) string {
    if dockerErr, ok := err.(*DockerError); ok {
        switch dockerErr.Operation {
        case "pull":
            return fmt.Sprintf("Failed to download Moodle image: %s", dockerErr.Image)
        case "run":
            return "Failed to start new Moodle container"
        case "start":
            return "Failed to restart existing Moodle container"
        default:
            return fmt.Sprintf("Docker operation failed: %s", dockerErr.Operation)
        }
    }
    return err.Error()
}
```

## Security Considerations

### Docker Socket Access

**Security Implications:**
- Docker socket access provides root-equivalent privileges
- Container operations can affect host system
- Image pulls can consume significant disk space and bandwidth

**Mitigation Strategies:**
- Validate all container IDs and image names
- Use read-only operations where possible
- Limit container capabilities (current implementation uses defaults)
- Implement proper input sanitization

### Container Security

**Port Binding:**
```go
const ContainerPort = "8080:8080"  // Bind to localhost only
```

**Security Features:**
- Containers run with default Docker security
- Port binding limited to localhost interface
- No privileged container operations
- Container isolation provided by Docker

### Credential Handling

**Local Storage:**
- Credentials stored in plain text (appropriate for development tool)
- Files stored in user-accessible directory
- No network transmission of credentials
- Temporary credential display in UI

## Performance Optimization

### Concurrent Operations

**Stream Processing:**
```go
// Process stdout and stderr concurrently
go func() {
    errChan <- progress.ProcessStream(stdout)
}()

go func() {
    errChan <- progress.ProcessStream(stderr)
}()
```

### Efficient Log Processing

**Polling Optimization:**
```go
// Adaptive polling intervals
if logErrorCount > maxLogErrors {
    time.Sleep(5 * time.Second)  // Reduced frequency on errors
} else {
    time.Sleep(2 * time.Second)  // Normal polling
}
```

### Resource Management

**Memory Management:**
- Stream processing to avoid loading large logs into memory
- Proper cleanup of command processes
- Efficient regex compilation (done once at startup)

**CPU Optimization:**
- Concurrent stream processing
- Non-blocking UI updates via events
- Efficient string parsing algorithms

### Caching Strategies

**Container State Caching:**
- Container ID persistence across application restarts
- Credential caching to avoid repeated log parsing
- Image existence caching (with validation)

**Health Check Optimization:**
- Periodic health checks with reasonable intervals
- Cached health status with timeout-based invalidation
- Efficient network connectivity testing

## Integration Testing

### Docker Environment Testing

**Test Scenarios:**
1. **No Docker**: Application behavior when Docker not installed
2. **Docker Not Running**: Daemon unavailable scenarios
3. **Permission Issues**: User lacks Docker access
4. **Network Issues**: Registry access problems
5. **Image Corruption**: Handling of corrupted or missing images

### Container Lifecycle Testing

**Test Cases:**
1. **First-Time Launch**: Complete image pull and container creation
2. **Subsequent Launches**: Container restart operations
3. **Multiple Restarts**: Container state consistency
4. **Forced Termination**: Recovery from abnormal container shutdown
5. **Long-Running Containers**: Extended operation stability

### Cross-Platform Testing

**Platform-Specific Tests:**
- Windows Docker Desktop integration
- macOS permission handling
- Linux native Docker daemon
- Path separator and file system differences
- Command execution and output parsing

This Docker integration documentation provides comprehensive coverage of all Docker-related operations and considerations in the Moodle Prototype Manager application.