# Moodle Prototype Manager - API Documentation

## Table of Contents

1. [Overview](#overview)
2. [Go Backend API](#go-backend-api)
3. [Application Context](#application-context)
4. [Docker Management](#docker-management)
5. [Storage Management](#storage-management)
6. [Error Handling](#error-handling)
7. [Logging and Utilities](#logging-and-utilities)
8. [Frontend JavaScript API](#frontend-javascript-api)

## Overview

This document provides comprehensive API documentation for the Moodle Prototype Manager application. The application follows a clean architecture pattern with clear separation between the Go backend (Wails) and JavaScript frontend.

### Architecture Pattern

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Frontend (JS)  │◄──►│  Go Backend     │◄──►│  Docker Engine  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  UI Components  │    │ Storage Layer   │    │  File System    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Go Backend API

### Application Context

#### `App` Struct

**File:** `app.go`

The main application context that coordinates all operations.

```go
type App struct {
    ctx               context.Context
    dockerManager     *docker.Manager
    credentialManager *storage.CredentialManager
    fileManager       *storage.FileManager
    logParser         *docker.LogParser
}
```

**Constructor:**
```go
func NewApp() *App
```

**Lifecycle Methods:**

##### `OnStartup(ctx context.Context)`
Called when the Wails application starts up.

**Purpose:** Initialize application context and load configuration.

**Process:**
1. Sets the application context
2. Loads Docker image configuration from `image.docker` file
3. Configures Docker manager with image name
4. Logs initialization status

**Error Handling:**
- Falls back to default image if configuration file is missing
- Logs warning messages for fallback scenarios

##### `OnShutdown(ctx context.Context)`
Called when the application is shutting down.

**Purpose:** Perform graceful cleanup of running containers.

**Process:**
1. Checks if a container ID file exists
2. Loads container ID and validates container status
3. Stops running container if found
4. Logs shutdown process

**Error Handling:**
- Attempts failsafe container stop on errors
- Logs all shutdown operations

### Main API Methods

#### `HealthCheck() map[string]bool`
**Export:** Frontend-callable via Wails

**Purpose:** Check Docker and Internet connectivity status.

**Returns:**
```go
map[string]bool{
    "docker":   bool, // Docker daemon accessibility
    "internet": bool, // Internet connectivity status
}
```

**Process:**
1. Calls `docker.PerformHealthChecks()`
2. Returns status map for frontend consumption

#### `RunMoodle() error`
**Export:** Frontend-callable via Wails

**Purpose:** Start the Moodle container (handles both new and existing containers).

**Process for Existing Container:**
1. Check for existing `container.id` file
2. Validate container ID and check container status
3. Start existing container using `docker.StartContainer()`
4. Wait for container readiness using HTTP checks
5. Update credentials with existing password

**Process for New Container:**
1. Validate Docker image name configuration
2. Check if image exists locally
3. Pull image with progress tracking if needed
4. Clear old credentials
5. Run new container using `docker.RunContainer()`
6. Save container ID to file
7. Extract credentials from container logs
8. Save credentials to file

**Error Handling:**
- Returns validation errors for invalid configurations
- Wraps Docker operation errors with context
- Provides detailed error messages for troubleshooting

#### `StopMoodle() error`
**Export:** Frontend-callable via Wails

**Purpose:** Stop the running Moodle container gracefully.

**Process:**
1. Validate container ID file exists
2. Load and validate container ID
3. Check container running status
4. Attempt graceful stop using `docker.StopContainer()`
5. Fallback to force stop if graceful stop fails

**Error Handling:**
- Validates container existence before stopping
- Attempts multiple stop strategies
- Logs all stop attempts and results

#### `GetCredentials() map[string]string`
**Export:** Frontend-callable via Wails

**Purpose:** Retrieve stored Moodle credentials.

**Returns:**
```go
map[string]string{
    "username": string, // Always "admin"
    "password": string, // Extracted admin password
    "url":      string, // Moodle URL (http://localhost:8080)
}
```

**Error Handling:**
- Returns default credentials on file read errors
- Validates credentials before returning
- Logs validation issues

#### `IsContainerReady() bool`
**Export:** Frontend-callable via Wails

**Purpose:** Check if the Moodle container is ready to accept connections.

**Logic:**
1. For existing containers: Test HTTP connectivity to port 8080
2. For new containers: Check if credentials file exists

**HTTP Test:**
- Makes HTTP GET request to `http://localhost:8080`
- 5-second timeout
- Any HTTP response indicates readiness

#### `OpenBrowser() error`
**Export:** Frontend-callable via Wails

**Purpose:** Open the default browser to the Moodle URL.

**Platform Support:**
- **macOS:** `open` command
- **Windows:** `rundll32 url.dll,FileProtocolHandler` command
- **Linux:** `xdg-open` command

**Process:**
1. Load credentials to get URL
2. Execute platform-specific browser open command
3. Set up command for platform (Windows requires special handling)

#### `GetImageName() string`
**Export:** Frontend-callable via Wails

**Purpose:** Return the current Docker image name for frontend display.

### Docker Management

#### `docker.Manager` Struct

**File:** `docker/manager.go`

Handles all Docker container operations.

```go
type Manager struct {
    imageName string
}
```

##### Core Methods

**`SetImageName(imageName string)`**
- Sets the Docker image name to use for operations
- Called during application startup

**`GetImageName() string`**
- Returns the current Docker image name
- Used for frontend display and validation

**`CheckImageExists() (bool, error)`**
- Executes: `docker images --format "{{.Repository}}:{{.Tag}}"`
- Searches for the configured image name in output
- Returns boolean indicating image existence

**`PullImage() error`**
- Executes: `docker pull <imageName>`
- Simple pull without progress tracking
- Used for basic image downloads

**`PullImageWithProgress(progressCallback func(float64, string)) error`**
- Advanced pull with real-time progress tracking
- Creates stdout and stderr pipes for progress monitoring
- Calls progress callback with percentage and status updates
- Used for user-facing image downloads

**`RunContainer() (string, error)`**
- Executes: `docker run -d -p 8080:8080 <imageName>`
- Returns container ID
- Validates returned container ID format
- Used for launching new containers

**`StartContainer(containerID string) error`**
- Executes: `docker start <containerID>`
- Used for restarting existing containers
- Validates container ID before execution

**`StopContainer(containerID string) error`**
- Executes: `docker stop <containerID>`
- Graceful container shutdown
- Used for normal container stopping

**`ForceStopContainer(containerID string) error`**
- Executes: `docker kill <containerID>`
- Forceful container termination
- Used as last resort when graceful stop fails

**`IsContainerRunning(containerID string) (bool, error)`**
- Executes: `docker inspect --format={{.State.Running}} <containerID>`
- Returns boolean indicating if container is running
- Used for container status checking

**`GetContainerLogs(containerID string) (string, error)`**
- Executes: `docker logs <containerID>`
- Returns complete container log output
- Uses CombinedOutput to capture both stdout and stderr

**`ValidateContainerID(containerID string) error`**
- Validates container ID format
- Checks if container exists using docker inspect
- Used before container operations

#### Docker Health Checks

**File:** `docker/health.go`

**`PerformHealthChecks() HealthStatus`**

```go
type HealthStatus struct {
    Docker   bool
    Internet bool
}
```

**Docker Check Process:**
1. Execute `docker --version`
2. Check command exit code
3. Verify Docker daemon accessibility

**Internet Check Process:**
1. HTTP GET request to `http://www.google.com`
2. 5-second timeout
3. Any successful response indicates connectivity

#### Log Parsing

**File:** `docker/logs.go`

**`LogParser` Methods:**

**`ExtractCredentials(logs string) Credentials`**

Parses container logs to extract Moodle credentials.

**Search Patterns:**
- Admin Password: `Generated admin password: <password>`
- Moodle URL: `Moodle is available at: http://localhost:8080`

**Returns:**
```go
type Credentials struct {
    Password string
    URL      string
}
```

#### Progress Tracking

**File:** `docker/progress.go`

**`PullProgress` Methods:**

**`NewPullProgress() *PullProgress`**
- Creates new progress tracker
- Initializes callback system

**`AddCallback(callback func(float64, string))`**
- Adds progress callback function
- Supports multiple callbacks

**`ProcessStream(reader io.Reader) error`**
- Processes Docker pull output stream
- Parses progress information
- Calls registered callbacks with updates

### Storage Management

#### File Management

**File:** `storage/files.go`

**`FileManager` Methods:**

**`getBaseDir() string`**
- Determines appropriate directory for file storage
- **Development Mode:** Uses working directory (detects `go.mod`)
- **Production Mode:** Uses `~/.moodle-prototype-manager`
- Critical fix for production builds on macOS/Windows

**`getUserDataDir() string`**
- Returns platform-specific user data directory
- **All Platforms:** `~/.moodle-prototype-manager`
- Fallback to current directory if home directory unavailable

**`SaveContainerID(containerID string) error`**
- Saves container ID to `container.id` file
- Single line, plain text format
- Creates directory if needed

**`LoadContainerID() (string, error)`**
- Loads container ID from `container.id` file
- Validates container ID format
- Returns cleaned, trimmed container ID

**`ContainerIDExists() bool`**
- Checks if `container.id` file exists
- Used for container state detection

**`SaveCredentials(password, url string) error`**
- Saves credentials to `moodle.txt` file
- Key-value format: `password=<value>\nurl=<value>`
- Creates directory if needed

**`LoadCredentials() (map[string]string, error)`**
- Loads credentials from `moodle.txt` file
- Parses key-value format
- Handles malformed lines gracefully

**`CredentialsExist() bool`**
- Checks if `moodle.txt` file exists
- Used for credentials state detection

**`LoadImageName() (string, error)`**
- Loads Docker image name from `image.docker` file
- Searches multiple potential paths:
  1. Base directory (getBaseDir result)
  2. Current directory
  3. Parent directory
  4. Working directory
  5. Executable directory
- Validates image name format

**`CleanupFiles() error`**
- Removes all storage files
- Deletes both `container.id` and `moodle.txt`
- Uses MultiError for comprehensive error reporting

#### Credential Management

**File:** `storage/credentials.go`

**`CredentialManager` Methods:**

**`Update(password, url string) error`**
- Updates stored credentials
- Validates input parameters
- Calls FileManager.SaveCredentials

**`Load() (Credentials, error)`**
- Loads credentials from storage
- Returns structured Credentials object
- Validates loaded data

**`Clear() error`**
- Removes credentials file
- Used when starting new containers

**`Exists() bool`**
- Checks if credentials file exists
- Used for state checking

```go
type Credentials struct {
    Username string // Always "admin"
    Password string // Extracted from logs
    URL      string // Default: "http://localhost:8080"
}

func (c Credentials) IsValid() bool
func (c Credentials) IsComplete() bool
func (c Credentials) ToMap() map[string]string
```

### Error Handling

**File:** `errors/errors.go`

The application uses a comprehensive error handling system with structured error types.

#### Error Types

**`DockerError`**
- Docker-specific operation errors
- Contains operation, image/container info, and output

**`FileError`**
- File system operation errors
- Contains operation, file path, and underlying error

**`NetworkError`**
- Network connectivity errors
- Contains operation and underlying error

**`ValidationError`**
- Input validation errors
- Contains field name, message, and invalid value

**`MultiError`**
- Collection of multiple errors
- Used for operations that can have multiple failure points

#### Validation Functions

**`ValidateContainerID(containerID string) error`**
- Validates Docker container ID format
- Ensures non-empty, trimmed string of appropriate length

**`ValidateImageName(imageName string) error`**
- Validates Docker image name format
- Checks for proper repository:tag format

**`ValidateNotEmpty(fieldName, value string) error`**
- Generic non-empty validation
- Used for required string fields

### Logging and Utilities

**File:** `utils/logger.go`

**Logging Functions:**

**`InitLogger()`**
- Initializes application logging
- Called during app startup

**`LogInfo(message string)`**
- Logs informational messages
- Used for normal operation tracking

**`LogError(message string, err error)`**
- Logs error messages with context
- Used for error reporting and debugging

**`LogWarning(message string)`**
- Logs warning messages
- Used for non-critical issues

**`LogDebug(message string)`**
- Logs debug messages
- Used for detailed operation tracking

#### Platform Utilities

**Files:** `utils/platform_windows.go`, `utils/platform_other.go`

**`SetupCommandForPlatform(cmd *exec.Cmd)`**
- Sets up OS-specific command execution
- **Windows:** Configures window hiding for commands
- **Other platforms:** No special configuration needed

## Frontend JavaScript API

### Application State Management

**File:** `frontend/js/app.js`

#### `AppState` Object

Central state management for the frontend application.

```javascript
export const AppState = {
    dockerStatus: false,      // Docker connectivity status
    internetStatus: false,    // Internet connectivity status
    containerRunning: false,  // Container running state
    credentials: {
        username: 'admin',
        password: '',
        url: 'http://localhost:8080'
    }
}
```

#### Core Functions

**`updateHealthCheckResults()`**
- Updates UI based on AppState health check results
- Controls button enable/disable state
- Updates status indicator colors
- Updates status text

**`updateStatusText(text)`**
- Updates footer status text
- Used for operation status communication

### UI Management

**File:** `frontend/js/ui.js`

**Modal Management:**

**`showDownloadModal()`**
- Displays Docker image download progress modal
- Initializes progress bar at 0%

**`hideDownloadModal()`**
- Closes download progress modal
- Resets progress bar state

**`updateDownloadProgress(percentage, status)`**
- Updates download progress bar and status text
- Called by progress event handlers

**`showStartupModal()`**
- Displays container startup waiting modal
- Shows loading spinner

**`hideStartupModal()`**
- Closes startup waiting modal

**`showBrowserDialog()`**
- Displays browser confirmation dialog
- Offers to open Moodle in browser

**`hideBrowserDialog()`**
- Closes browser confirmation dialog

**Button State Management:**

**`setButtonRunning()`**
- Changes button text to "Stop Moodle"
- Changes button color to red/stop state

**`setButtonStopped()`**
- Changes button text to "Run Moodle"
- Changes button color to green/run state

**Credential Display:**

**`showCredentials(credentials)`**
- Displays extracted Moodle credentials
- Updates password field with copy button
- Shows credential table

**`hideCredentials()`**
- Hides credential display table

### Event Handling

**File:** `frontend/js/events.js`

**Wails Integration:**

**Event Listeners:**
- `docker:pull:progress` - Docker image download progress
- Container status updates
- Error message display

**Backend Method Calls:**
```javascript
// Health checks
window.wails.Go.main.App.HealthCheck()

// Container operations
window.wails.Go.main.App.RunMoodle()
window.wails.Go.main.App.StopMoodle()

// Data retrieval
window.wails.Go.main.App.GetCredentials()
window.wails.Go.main.App.IsContainerReady()

// Browser operations
window.wails.Go.main.App.OpenBrowser()
```

**Container Management Functions:**

**`startMoodleContainer()`**
- Initiates container start process
- Shows appropriate modals during operation
- Handles progress events and error states
- Updates UI based on operation results

**`stopMoodleContainer()`**
- Initiates container stop process
- Updates button states during operation
- Handles error conditions

**Health Check Management:**

**`performHealthChecks()`**
- Executes backend health check calls
- Updates AppState with results
- Triggers UI updates

**`startHealthCheckLoop()`**
- Runs periodic health checks
- Updates UI indicators continuously

### Browser Integration

**Copy to Clipboard:**
- Password copy button functionality
- Uses navigator.clipboard API where available
- Fallback for older browsers

**URL Handling:**
- Click handler for URL links
- Option to open in browser or copy URL

## Error Handling Patterns

### Backend Error Handling

1. **Validation Errors**: Input parameter validation with descriptive messages
2. **Docker Errors**: Operation-specific errors with Docker output context
3. **File Errors**: File system operation errors with path information
4. **Network Errors**: Connectivity errors with timeout information
5. **Multi Errors**: Collection of related errors for complex operations

### Frontend Error Handling

1. **Promise Rejection**: Backend call failures are caught and displayed
2. **UI State Management**: Error states update UI appropriately
3. **User Notification**: Error messages are shown in status text
4. **Graceful Degradation**: Partial functionality when possible

### Error Logging Strategy

1. **Structured Logging**: Consistent log format with context
2. **Error Wrapping**: Preserve error chain with additional context
3. **Debug Information**: Detailed logging for troubleshooting
4. **User-Friendly Messages**: Clear error communication to users

## Integration Points

### Wails Framework Integration

- **Context Binding**: Go methods exposed to JavaScript
- **Event System**: Real-time progress updates
- **Asset Embedding**: Frontend assets bundled with executable

### Docker Integration

- **Command Execution**: Cross-platform Docker command handling
- **Stream Processing**: Real-time output parsing
- **Error Handling**: Docker-specific error interpretation

### File System Integration

- **Platform Awareness**: OS-specific directory handling
- **Permission Management**: Proper file access permissions
- **Production Fix**: User data directory for production builds

This API documentation provides comprehensive coverage of all public interfaces and internal architecture patterns used in the Moodle Prototype Manager application.