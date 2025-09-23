# Moodle Prototype Manager - Development Setup and Workflow Guide

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Development Environment Setup](#development-environment-setup)
3. [Project Structure](#project-structure)
4. [Development Workflow](#development-workflow)
5. [Testing Strategy](#testing-strategy)
6. [Code Quality and Standards](#code-quality-and-standards)
7. [Debugging and Troubleshooting](#debugging-and-troubleshooting)
8. [Performance Optimization](#performance-optimization)
9. [Cross-Platform Development](#cross-platform-development)
10. [Continuous Integration](#continuous-integration)

## Prerequisites

### System Requirements

**Operating Systems:**
- **Windows**: Windows 10/11 (64-bit)
- **macOS**: macOS 10.15+ (Intel or Apple Silicon)
- **Linux**: Ubuntu 18.04+ or equivalent (for development)

**Development Tools:**
- **Go**: Version 1.19+ (latest stable recommended)
- **Node.js**: Version 16+ (for frontend tooling, if needed)
- **Docker Desktop**: Latest version for container testing
- **Git**: Version 2.25+ for version control

**Code Editors (Recommended):**
- **VS Code**: With Go extension and Wails extension
- **GoLand**: Full IDE with excellent Go support
- **Vim/Neovim**: With gopls LSP integration

### Docker Requirements

**Docker Desktop Setup:**
```bash
# Verify Docker installation
docker --version
docker info

# Test Docker functionality
docker run hello-world
```

**Required Docker Permissions:**
- Access to Docker daemon socket
- Ability to pull images from public registries
- Port binding permissions (8080:8080)

### Go Environment Setup

**Go Installation Verification:**
```bash
# Check Go version
go version

# Verify GOPATH and GOROOT
go env GOPATH
go env GOROOT

# Ensure Go modules are enabled
go env GO111MODULE  # Should be 'on' or empty
```

**Required Go Tools:**
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install additional development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
```

## Development Environment Setup

### Initial Project Setup

**1. Clone Repository:**
```bash
git clone <repository-url>
cd moodle-prototype-manager
```

**2. Install Dependencies:**
```bash
# Install Go dependencies
go mod download
go mod verify

# Verify Wails installation
wails doctor
```

**3. Configure Docker Image:**
```bash
# Create image configuration file
echo "wenkhairu/moodle-prototype:502-stable" > image.docker

# Or use custom image for development
echo "your-custom/moodle-image:dev" > image.docker
```

**4. Verify Setup:**
```bash
# Run development setup script
./dev.sh --help

# Verify project can build
wails build
```

### IDE Configuration

#### VS Code Setup

**Required Extensions:**
- Go (golang.go)
- Wails (wails-app.wails)
- JavaScript/HTML/CSS support (built-in)
- Docker (ms-azuretools.vscode-docker)

**Settings (`.vscode/settings.json`):**
```json
{
    "go.toolsManagement.checkForUpdates": "local",
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.lintFlags": ["--fast"],
    "files.associations": {
        "*.go": "go"
    },
    "go.testFlags": ["-v", "-race"],
    "go.buildFlags": ["-race"]
}
```

**Launch Configuration (`.vscode/launch.json`):**
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Development",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "env": {
                "WAILS_ENV": "development"
            }
        },
        {
            "name": "Run Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}",
            "args": ["-v", "./..."]
        }
    ]
}
```

#### GoLand Configuration

**Project Settings:**
- Go Modules: Enabled
- GOROOT: Automatically detected
- Build Tags: `development` for dev builds
- Code Style: Go formatting with goimports

**Run Configurations:**
```
Name: Development Server
Type: Go Build
Package: .
Output directory: build/bin/
Working directory: $ProjectFileDir$
Environment: WAILS_ENV=development
```

### Development Scripts

#### Primary Development Script (`dev.sh`)

**Usage Examples:**
```bash
# Start development server
./dev.sh

# Clean and start with specific image
./dev.sh --clean --image wenkhairu/moodle-prototype:dev

# Show help
./dev.sh --help
```

**Script Features:**
- Automatic frontend file preparation
- Docker image configuration
- Wails development server startup
- Live reload for Go code changes
- Frontend file synchronization

#### Additional Scripts

**Build Script (`build.sh`):**
```bash
#!/bin/bash
# Cross-platform build script

# Clean previous builds
rm -rf build/

# Build for multiple platforms
wails build -platform darwin/amd64,darwin/arm64,windows/amd64

# Create distribution packages
# (Implementation depends on distribution strategy)
```

**Test Script:**
```bash
#!/bin/bash
# Comprehensive testing script

# Run Go tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Run linting
golangci-lint run ./...

# Run static analysis
staticcheck ./...
```

## Project Structure

### Directory Layout

```
moodle-prototype-manager/
├── cmd/                        # Command-line tools (if any)
├── docs/                       # Documentation
│   ├── api-documentation.md
│   ├── frontend-documentation.md
│   ├── docker-integration.md
│   ├── development-guide.md
│   └── ...
├── docker/                     # Docker operations
│   ├── manager.go             # Core Docker manager
│   ├── health.go              # Health check functionality
│   ├── logs.go                # Log parsing
│   ├── progress.go            # Pull progress tracking
│   ├── path.go                # Cross-platform paths
│   └── *_test.go              # Unit tests
├── errors/                     # Error handling
│   ├── errors.go              # Error types and functions
│   └── errors_test.go         # Error handling tests
├── frontend/                   # Frontend source files
│   ├── index.html             # Main HTML structure
│   ├── css/
│   │   └── styles.css         # Application styles
│   ├── js/
│   │   ├── app.js             # Main application logic
│   │   ├── ui.js              # UI manipulation
│   │   └── events.js          # Event handling
│   ├── assets/
│   │   └── images/            # Static assets
│   └── dist/                  # Built frontend (auto-generated)
├── storage/                    # File storage operations
│   ├── files.go               # File I/O operations
│   ├── credentials.go         # Credential management
│   └── *_test.go              # Storage tests
├── utils/                      # Utility functions
│   ├── logger.go              # Logging functionality
│   ├── platform_windows.go   # Windows-specific utilities
│   ├── platform_other.go     # Unix-like utilities
│   └── *_test.go              # Utility tests
├── main.go                    # Application entry point
├── app.go                     # Wails application context
├── wails.json                 # Wails configuration
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── dev.sh                     # Development script
├── build.sh                   # Build script
├── image.docker               # Docker image configuration
├── README.md                  # Project overview
└── .gitignore                # Git ignore rules
```

### Code Organization Principles

**Package Structure:**
- **Single Responsibility**: Each package has a clear, focused purpose
- **Dependency Direction**: Dependencies flow toward core business logic
- **Interface Segregation**: Small, focused interfaces
- **Testability**: All packages designed for easy unit testing

**Naming Conventions:**
- **Packages**: Short, lowercase, single words when possible
- **Files**: Descriptive names with Go-style snake_case for tests
- **Functions**: CamelCase with verb-noun structure
- **Variables**: Clear, descriptive names avoiding abbreviations

## Development Workflow

### Daily Development Cycle

**1. Start Development Environment:**
```bash
# Pull latest changes
git pull origin main

# Ensure dependencies are current
go mod tidy

# Start development server
./dev.sh
```

**2. Development Process:**
```bash
# Create feature branch
git checkout -b feature/new-functionality

# Make changes to Go code
# Wails automatically reloads on save

# Make changes to frontend
# Re-run dev.sh to sync frontend files
./dev.sh
```

**3. Testing and Validation:**
```bash
# Run specific tests
go test -v ./docker/...

# Run all tests
go test -v ./...

# Check code quality
golangci-lint run ./...

# Build to verify
wails build
```

**4. Commit and Push:**
```bash
# Add changes
git add .

# Commit with descriptive message
git commit -m "Add: Docker health check improvements

- Enhanced error handling for Docker daemon detection
- Added timeout configuration for health checks
- Improved cross-platform compatibility"

# Push to remote
git push origin feature/new-functionality
```

### Feature Development Process

#### 1. Planning Phase

**Requirements Analysis:**
- Define user story and acceptance criteria
- Identify affected components (frontend, backend, Docker integration)
- Plan API changes and data structures
- Consider cross-platform implications

**Design Decisions:**
- Update architecture documentation if needed
- Plan error handling strategy
- Consider performance implications
- Plan testing approach

#### 2. Implementation Phase

**Backend Development:**
```go
// Follow established patterns
func (m *Manager) NewFeature() error {
    // 1. Input validation
    if err := validateInput(); err != nil {
        return errors.WrapWithContext(err, "invalid input for new feature")
    }

    // 2. Core logic
    result, err := m.performOperation()
    if err != nil {
        return errors.NewOperationError("new_feature", err)
    }

    // 3. State updates
    m.updateState(result)

    return nil
}
```

**Frontend Development:**
```javascript
// Follow modular approach
export async function handleNewFeature() {
    try {
        // 1. Update UI state
        setButtonLoading('feature-btn', true, 'Processing...');

        // 2. Call backend
        const result = await window.go.main.App.NewFeature();

        // 3. Handle success
        showNotification('Feature completed successfully', 'success');
        updateUIState(result);

    } catch (error) {
        // 4. Handle errors
        console.error('Feature failed:', error);
        showNotification(`Feature failed: ${error.message}`, 'error');
    } finally {
        // 5. Cleanup
        setButtonLoading('feature-btn', false);
    }
}
```

#### 3. Testing Phase

**Unit Tests:**
```go
func TestNewFeature(t *testing.T) {
    // Arrange
    manager := NewManager()
    manager.SetImageName("test-image")

    // Act
    err := manager.NewFeature()

    // Assert
    assert.NoError(t, err)
    assert.True(t, manager.featureCompleted)
}

func TestNewFeatureWithError(t *testing.T) {
    // Test error conditions
    manager := NewManager()
    // Don't set image name to trigger validation error

    err := manager.NewFeature()

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid input")
}
```

**Integration Tests:**
```go
func TestFeatureEndToEnd(t *testing.T) {
    // Skip if Docker not available
    if !isDockerAvailable() {
        t.Skip("Docker not available for integration tests")
    }

    // Full workflow test
    app := NewApp()
    // ... test complete feature workflow
}
```

### Git Workflow

**Branch Strategy:**
- `main`: Stable, deployable code
- `develop`: Integration branch for features
- `feature/name`: Feature development branches
- `hotfix/name`: Critical bug fixes

**Commit Message Format:**
```
Type: Brief description (50 chars max)

Longer explanation of the change, including:
- Why the change was made
- What was changed
- Any breaking changes
- Reference to issues (#123)

Examples:
- Add: New Docker health check system
- Fix: Container restart logic on macOS
- Update: Frontend styling for better UX
- Refactor: Error handling consolidation
```

**Pull Request Process:**
1. Create feature branch from `develop`
2. Implement feature with tests
3. Update documentation if needed
4. Create pull request with description
5. Code review and approval
6. Merge to `develop`
7. Regular integration to `main`

## Testing Strategy

### Test Types and Coverage

**Unit Tests (80% of tests):**
- Individual function testing
- Mock external dependencies
- Test error conditions
- Validate edge cases

**Integration Tests (15% of tests):**
- Component interaction testing
- Docker integration scenarios
- File system operations
- Cross-platform compatibility

**End-to-End Tests (5% of tests):**
- Complete user workflows
- GUI automation (if needed)
- Performance testing
- Deployment validation

### Test Organization

**Test File Structure:**
```go
// docker/manager_test.go
func TestManager_CheckImageExists(t *testing.T) { ... }
func TestManager_CheckImageExists_WithError(t *testing.T) { ... }
func TestManager_PullImage(t *testing.T) { ... }

// Table-driven tests
func TestValidateContainerID(t *testing.T) {
    tests := []struct {
        name        string
        containerID string
        wantError   bool
    }{
        {"valid ID", "abc123def456", false},
        {"empty ID", "", true},
        {"short ID", "abc", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateContainerID(tt.containerID)
            if (err != nil) != tt.wantError {
                t.Errorf("ValidateContainerID() error = %v, wantError %v",
                        err, tt.wantError)
            }
        })
    }
}
```

### Mock and Stub Strategies

**Docker Command Mocking:**
```go
type MockDockerManager struct {
    imageExists bool
    pullError   error
    runResult   string
}

func (m *MockDockerManager) CheckImageExists() (bool, error) {
    return m.imageExists, nil
}

func (m *MockDockerManager) PullImage() error {
    return m.pullError
}
```

**File System Mocking:**
```go
type MockFileManager struct {
    files map[string]string
    errors map[string]error
}

func (m *MockFileManager) LoadContainerID() (string, error) {
    if err, exists := m.errors["container.id"]; exists {
        return "", err
    }
    return m.files["container.id"], nil
}
```

### Test Execution

**Running Tests:**
```bash
# All tests
go test ./...

# Specific package
go test ./docker/...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Race condition detection
go test -race ./...

# Verbose output
go test -v ./...

# Benchmark tests
go test -bench=. ./...
```

**Continuous Testing:**
```bash
# Watch mode (using external tool like entr)
find . -name "*.go" | entr -r go test ./...

# Or use Wails development mode which includes test running
wails dev
```

## Code Quality and Standards

### Go Code Standards

**Formatting and Style:**
```bash
# Format code
gofmt -w .

# Import organization
goimports -w .

# Linting
golangci-lint run ./...

# Static analysis
staticcheck ./...
```

**Code Review Checklist:**
- [ ] Proper error handling with context
- [ ] Input validation for all public functions
- [ ] Thread safety considerations
- [ ] Memory leak prevention
- [ ] Cross-platform compatibility
- [ ] Comprehensive test coverage
- [ ] Documentation for public APIs
- [ ] Consistent naming conventions

### Frontend Code Standards

**JavaScript Style:**
- ES6+ syntax with modules
- Consistent indentation (2 spaces)
- Clear function and variable names
- Error handling for all async operations
- Comments for complex logic

**CSS Standards:**
- BEM methodology for class naming
- Consistent color palette usage
- Mobile-first responsive design
- Performance-optimized selectors
- Documentation for design system components

### Documentation Standards

**Go Documentation:**
```go
// Manager handles Docker container operations for Moodle prototypes.
// It provides methods for container lifecycle management, health checking,
// and image operations with cross-platform compatibility.
type Manager struct {
    imageName string
}

// CheckImageExists verifies if the configured Docker image exists locally.
// It executes 'docker images' command and searches for the image name.
//
// Returns:
//   - bool: true if image exists, false otherwise
//   - error: any error encountered during the check
//
// Example:
//   manager := NewManager()
//   manager.SetImageName("nginx:latest")
//   exists, err := manager.CheckImageExists()
//   if err != nil {
//       log.Fatal("Failed to check image:", err)
//   }
func (m *Manager) CheckImageExists() (bool, error) { ... }
```

**README and Markdown Standards:**
- Clear headings and structure
- Code examples with proper syntax highlighting
- Step-by-step instructions
- Cross-references to related documentation
- Regular updates with changes

## Debugging and Troubleshooting

### Development Debugging

**Go Debugging with Delve:**
```bash
# Install delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug specific test
dlv test ./docker/ -- -test.run TestManagerCheckImageExists

# Debug main application
dlv debug . -- --dev
```

**VS Code Debugging:**
- Use integrated debugger with launch.json configuration
- Set breakpoints in Go code
- Inspect variable values and call stack
- Step through code execution

### Logging and Diagnostics

**Application Logging:**
```go
// Enhanced logging for development
func (a *App) RunMoodle() error {
    utils.LogInfo("RunMoodle called by frontend")

    // Log important state
    utils.LogDebug(fmt.Sprintf("Current image: %s", a.dockerManager.GetImageName()))
    utils.LogDebug(fmt.Sprintf("Container exists: %v", a.fileManager.ContainerIDExists()))

    // Log operation results
    if err := operation(); err != nil {
        utils.LogError("Operation failed", err)
        return err
    }

    utils.LogInfo("RunMoodle completed successfully")
    return nil
}
```

**Frontend Debugging:**
```javascript
// Enhanced console logging
console.group('Container Start Operation');
console.log('AppState:', AppState);
console.log('Health Status:', { docker: AppState.dockerStatus, internet: AppState.internetStatus });

try {
    const result = await window.go.main.App.RunMoodle();
    console.log('Backend result:', result);
} catch (error) {
    console.error('Backend error:', error);
    console.trace(); // Show call stack
} finally {
    console.groupEnd();
}
```

### Common Development Issues

**Wails Build Issues:**
```bash
# Clean and rebuild
wails clean
wails build

# Check Wails doctor
wails doctor

# Update Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**Go Module Issues:**
```bash
# Clean module cache
go clean -modcache
go mod download

# Verify module integrity
go mod verify

# Update dependencies
go mod tidy
```

**Docker Integration Issues:**
```bash
# Test Docker connectivity
docker version
docker info

# Check image availability
docker images | grep moodle-prototype

# Test container creation
docker run -d -p 8080:8080 wenkhairu/moodle-prototype:502-stable
```

## Performance Optimization

### Go Performance

**Profiling:**
```go
import _ "net/http/pprof"

// Add to main() for development builds
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

**Memory Optimization:**
```go
// Efficient string building
var builder strings.Builder
builder.WriteString("prefix")
builder.WriteString(variable)
result := builder.String()

// Pool expensive objects
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}
```

**Concurrency Optimization:**
```go
// Use bounded goroutines
semaphore := make(chan struct{}, 10) // Max 10 concurrent operations

for _, item := range items {
    semaphore <- struct{}{} // Acquire
    go func(item string) {
        defer func() { <-semaphore }() // Release
        processItem(item)
    }(item)
}
```

### Frontend Performance

**DOM Optimization:**
```javascript
// Cache DOM queries
const elements = {
    runButton: document.getElementById('run-moodle-btn'),
    statusText: document.getElementById('status-text'),
    credentialsDisplay: document.getElementById('credentials-display')
};

// Batch DOM updates
function updateUI(state) {
    // Update all elements at once to minimize reflow
    elements.runButton.textContent = state.buttonText;
    elements.runButton.disabled = state.buttonDisabled;
    elements.statusText.textContent = state.statusText;
}
```

**Memory Management:**
```javascript
// Cleanup event listeners
function cleanup() {
    window.removeEventListener('resize', handleResize);
    clearInterval(healthCheckInterval);
}

// Efficient event handling
const handleButtonClick = debounce((event) => {
    // Handle click
}, 300);
```

## Cross-Platform Development

### Platform-Specific Code

**Build Constraints:**
```go
//go:build windows
// +build windows

package utils

import "syscall"

func SetupCommandForPlatform(cmd *exec.Cmd) {
    cmd.SysProcAttr = &syscall.SysProcAttr{
        HideWindow: true,
    }
}
```

```go
//go:build !windows
// +build !windows

package utils

func SetupCommandForPlatform(cmd *exec.Cmd) {
    // No special setup needed for Unix-like systems
}
```

### Testing on Multiple Platforms

**GitHub Actions Matrix:**
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
    go-version: [1.19, 1.20]
```

**Local Testing:**
```bash
# Cross-compile for testing
GOOS=windows GOARCH=amd64 go build -o build/windows/
GOOS=darwin GOARCH=amd64 go build -o build/macos-intel/
GOOS=darwin GOARCH=arm64 go build -o build/macos-apple/
```

## Continuous Integration

### GitHub Actions Workflow

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go-version: [1.19, 1.20]

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v3

    - name: Lint
      uses: golangci/golangci-lint-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.20

    - name: Install Wails
      run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

    - name: Build application
      run: wails build -platform darwin/amd64,darwin/arm64,windows/amd64
```

### Quality Gates

**Pre-commit Hooks:**
```bash
# Install pre-commit
pip install pre-commit

# Configure hooks in .pre-commit-config.yaml
go install github.com/pre-commit/pre-commit-hooks
```

**Quality Metrics:**
- Test coverage > 80%
- No linting errors
- Successful build on all platforms
- No known security vulnerabilities
- Performance benchmarks within limits

This development guide provides comprehensive coverage of all aspects of developing, testing, and maintaining the Moodle Prototype Manager application.