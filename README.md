# Moodle Prototype Manager

A sleek, cross-platform desktop application built with Go and Wails that provides a modern graphical interface for managing Moodle prototype Docker containers. The application handles the complete lifecycle of Moodle prototypes including automated image management, container orchestration, credential extraction, and seamless browser integration.

![Moodle Prototype Manager](frontend/assets/images/moodle-logo.png)

## ✨ Features

- **🐳 Intelligent Container Management**: Automated Docker image pulling, container lifecycle management with state persistence
- **🔍 Real-time Health Monitoring**: Continuous Docker daemon and internet connectivity checks with visual indicators
- **🔑 Automatic Credential Extraction**: Smart parsing of container logs to extract admin credentials automatically
- **🌐 Seamless Browser Integration**: One-click browser launch with user consent dialogs
- **💾 Persistent State Management**: Remembers container state and credentials between application launches
- **🎨 Modern Compact UI**: Clean, responsive 400x400 interface with Moodle-branded styling and progress indicators
- **🔌 Graceful Lifecycle Management**: Automatic container cleanup on application exit
- **💻 True Cross-Platform**: Native support for macOS (Intel & Apple Silicon) and Windows
- **⚙️ Configurable Docker Images**: Easy Docker image switching via configuration file

## 📋 Prerequisites

### Required Dependencies
- **Docker Desktop**: Download and install from [docker.com](https://www.docker.com/products/docker-desktop)
  - Must be running before launching the application
  - Requires proper user permissions for Docker commands
- **Operating System**:
  - **macOS**: 10.14 Mojave or later (Intel and Apple Silicon supported)
  - **Windows**: Windows 10 version 1903 or later (64-bit)

### System Requirements
- **Memory**: 4GB RAM minimum (8GB recommended for optimal performance)
- **Storage**: 3GB free space (2GB for Docker image + 1GB for application and temporary files)
- **Network**: Stable internet connection for initial Docker image download (~500MB-1GB)
- **Ports**: Port 8080 must be available (used for Moodle container mapping)

## 🚀 Installation

### Option 1: Download Pre-built Release (Recommended)

1. **Visit Releases**: Go to the [GitHub Releases](../../releases) page
2. **Download for Your Platform**:
   - **macOS Intel**: `moodle-502-prototype-manager-darwin-amd64.zip`
   - **macOS Apple Silicon**: `moodle-502-prototype-manager-darwin-arm64.zip`
   - **Windows 64-bit**: `moodle-502-prototype-manager-windows-amd64.zip`
3. **Extract and Install**:
   - **macOS**: Extract `.zip` and drag `moodle-prototype-manager.app` to Applications folder
   - **Windows**: Extract `.zip` to desired location (contains `moodle-prototype-manager.exe` and `image.docker`)
4. **First Launch**:
   - **macOS**: Right-click app → Open (to bypass Gatekeeper warning on first run)
   - **Windows**: Run `moodle-prototype-manager.exe` (may require Windows Defender SmartScreen approval)

### Option 2: Build from Source

Perfect for developers who want to customize the application or contribute to the project.

#### Prerequisites for Building

- **Go**: Version 1.21 or higher ([Download Go](https://golang.org/dl/))
- **Wails CLI**: Latest version (auto-installed by build script)
- **Operating System**: macOS 10.14+ or Windows 10+ with development tools
- **Git**: For cloning the repository

#### Quick Build Process

```bash
# Clone the repository
git clone https://github.com/khairu-aqsara/moodle-prototype-manager.git
cd moodle-prototype-manager

# Build using the automated script (recommended)
./build.sh

# Or build for specific platforms
./build.sh darwin-arm64      # macOS Apple Silicon
./build.sh darwin-amd64      # macOS Intel
./build.sh windows-amd64     # Windows 64-bit
./build.sh all               # All platforms

# Build with custom Docker image
./build.sh --image wenkhairu/moodle-prototype:502-alpine
```

#### Development Mode

For rapid development with live reload:

```bash
# Start development server (recommended)
./dev.sh

# Clean frontend cache and restart
./dev.sh --clean

# Develop with different Docker image
./dev.sh --image wenkhairu/moodle-prototype:502-stable

# Show development options
./dev.sh --help
```

The development server provides:
- **Live Go Reload**: Automatic restart when backend Go code changes
- **Frontend Sync**: Manual sync of HTML/CSS/JS changes (re-run `./dev.sh`)
- **Debug Browser**: Optional browser debugging interface

## 🔐 Code Signing & Distribution

### macOS Code Signing

For distribution on macOS, applications should be properly signed to avoid security warnings.

#### Prerequisites
- **Apple Developer Account**: Required for signing certificates
- **Valid Signing Certificate**: One of the following installed in Keychain:
  - `Developer ID Application`: For distribution outside App Store
  - `Apple Development`: For development and testing
  - `Apple Distribution`: For App Store submission

#### Automatic Signing

```bash
# Build the application first
./build.sh darwin-arm64  # or darwin-amd64

# Auto-sign with first available certificate
./sign-mac.sh --auto

# Auto-sign and create notarization package
./sign-mac.sh --auto --zip
```

#### Manual Signing

```bash
# List available signing identities
./sign-mac.sh

# Sign with specific identity (interactive mode)
./sign-mac.sh -i "Developer ID Application: Your Name (TEAM123456)"

# Create ZIP for notarization
./sign-mac.sh --auto --zip
```

#### Notarization (Optional)

For distribution outside the App Store without security warnings:

```bash
# Sign and create ZIP
./sign-mac.sh --auto --zip

# Submit for notarization (requires Apple ID app-specific password)
xcrun notarytool submit build/moodle-prototype-manager-signed.zip \
  --apple-id your-apple-id@example.com \
  --password your-app-specific-password \
  --team-id YOUR_TEAM_ID \
  --wait

# After successful notarization, staple the ticket
xcrun stapler staple build/bin/moodle-prototype-manager.app
```

### Ad-hoc Signing (Development)

For local development without Apple Developer account:

```bash
# Use the ad-hoc signing script
./sign-adhoc.sh
```

This creates a locally-signed version that will run on your development machine.

## 🖥️ User Interface Guide

The application features a compact, modern 400x400 pixel interface for efficiency.

### Main Components

#### Header Section
- **Moodle Logo**: Resized to 120x72 pixels for compact display
- **Docker Image Name**: Dynamically displays current image from `image.docker` file
- **Version Information**: Application version displayed below logo

#### Central Action Button
- **"Run Moodle"** (Moodle Orange Gradient): Starts/restarts the container
- **"Stop Moodle"** (Red Gradient): Gracefully stops the running container
- **Disabled State**: Grayed out when health checks fail
- **Modern Design**: Pill-shaped with gradient background and smooth hover animations

#### Credentials Display (When Container Running)
```
Moodle Login Details
┌─────────────┬─────────────────────────────────┐
│ Username    │ admin                           │
│ Password    │ [randomly_generated_password]   │
│ URL         │ http://localhost:8080           │
└─────────────┴─────────────────────────────────┘
```
- **Copy-Friendly**: All text is selectable for easy copying
- **Monospace Font**: Password displayed in monospace for clarity
- **Compact Layout**: Optimized table design for 400x400 window

#### Footer Status Bar
- **Left Side**: Real-time status messages
  - "Checking..." → "Ready" → "Starting..." → "Running" → "Stopping"
  - Error messages when issues occur
- **Right Side**: Health indicators with colored circles
  - **🔴 Red**: Service unavailable
  - **🟡 Yellow**: Checking/transitioning
  - **🟢 Green**: Service available and healthy

### Modal Dialogs

#### Download Progress Modal
- **Real-time Progress**: Live parsing of Docker pull output
- **Percentage Display**: Current download percentage
- **Auto-dismiss**: Closes automatically when download completes

#### Startup Modal
- **Loading Animation**: Smooth CSS spinner animation
- **Status Message**: "Starting Moodle, please wait..."
- **Timeout Protection**: 5-minute maximum wait time
- **Auto-dismiss**: Closes when container logs show ready state

#### Browser Confirmation Dialog
- **Message**: "Would you like to open Moodle in your browser?"
- **Options**: "Yes" (Moodle orange) / "No" (gray) buttons
- **Remember Choice**: Option for future launches (if implemented)

## ⚙️ Configuration & Customization

### Docker Image Configuration

The application uses a flexible configuration system for Docker images:

#### Configuration File: `image.docker`
- **Location**: Application root directory (same as executable)
- **Format**: Single line containing Docker image name
- **Default**: `wenkhairu/moodle-prototype:502-alpine`

#### Changing Docker Image

**Edit Configuration File**
```bash
# Edit the image.docker file
echo "wenkhairu/moodle-prototype:502-stable" > image.docker
```

#### Available Image Variants
- `wenkhairu/moodle-prototype:502-amd64`: Debian-based (amd64)
- `wenkhairu/moodle-prototype:502-alpine`: Lightweight Alpine-based (arm64/Apple Silicon)

### Runtime Files

The application creates and manages these files:

#### `image.docker`
- **Purpose**: Specifies which Docker image to use
- **Format**: Plain text, single line
- **Example**: `wenkhairu/moodle-prototype:502-alpine`
- **Bundled**: Included in built applications

#### `container.id`
- **Purpose**: Stores the active container ID for state persistence
- **Format**: Plain text, single line Docker container ID
- **Example**: `a1b2c3d4e5f6...`
- **Lifecycle**: Created when container starts, deleted when manually stopped

#### `moodle.txt`
- **Purpose**: Stores extracted Moodle credentials
- **Format**: Key-value pairs
- **Example**:
  ```
  password=Kx9P2mL8qR5t
  url=http://localhost:8080
  ```
- **Security**: Contains sensitive data - protect access appropriately

### Network Configuration

- **Container Port**: 8080 (fixed)
- **Host Port**: 8080 (mapped from container)
- **Access URL**: `http://localhost:8080`
- **Network Mode**: Bridge (default Docker)
- **External Access**: Localhost only (no external network exposure)

## 🔄 Application Flow & Usage

### First-Time Startup Flow

1. **Application Launch**
   - Load configuration from `image.docker`
   - Initialize UI components
   - Set button to disabled state

2. **Health Check Sequence** (Automatic)
   - **Docker Check**: Verify Docker daemon accessibility (`docker --version`)
   - **Internet Check**: Verify connectivity to Docker registry
   - **UI Updates**: Color-coded status indicators update in real-time
   - **Button Activation**: "Run Moodle" button enabled when all checks pass

3. **First Container Launch** (User clicks "Run Moodle")
   - Check if Docker image exists locally
   - If not found: Show download modal with progress
   - Download Docker image with real-time progress updates
   - Create and start new container with port mapping `8080:8080`
   - Store container ID in `container.id` file
   - Show startup modal with loading animation

4. **Credential Extraction** (Automatic)
   - Monitor container logs every 2-3 seconds
   - Parse logs for admin password pattern: `Generated admin password: [password]`
   - Parse logs for URL pattern: `Moodle is available at: [url]`
   - Save credentials to `moodle.txt`
   - Display credentials in UI table
   - Dismiss startup modal

5. **Browser Integration** (Optional)
   - Show browser confirmation dialog
   - If user agrees: Launch default browser to `http://localhost:8080`
   - Button changes to "Stop Moodle" (red)

### Subsequent Startup Flow

1. **Quick Start Process**
   - Check for existing `container.id` file
   - Verify container still exists: `docker container [id] inspect`
   - If container exists: Start existing container
   - If container missing: Fall back to first-time flow

2. **Credential Loading**
   - Load credentials from `moodle.txt`
   - Display immediately (no extraction wait)
   - Container starts faster (no image download)

3. **State Management**
   - Button immediately shows correct state
   - Status bar shows appropriate messages
   - Health checks still performed for Docker/Internet

### Container Stop Flow

1. **User Clicks "Stop Moodle"**
   - Button shows loading state
   - Execute graceful container stop: `docker container [id] stop`
   - Hide credentials display
   - Button changes to "Run Moodle" (green)
   - Status bar updates to "Ready"

### Application Exit Flow

1. **Exit Initiated** (User closes app)
   - Check if container is running
   - If running: Execute graceful stop automatically
   - Clean up any temporary files
   - Save application state
   - Exit application

### Error Recovery Flows

#### Docker Not Running
1. Red circle for Docker status
2. Button remains disabled
3. Status message: "Docker not running - please start Docker Desktop"
4. Automatic retry every 5 seconds

#### Internet Connection Issues
1. Yellow circle during check, red if failed
2. Can still run if image exists locally
3. Download will fail with appropriate error message

#### Port 8080 Conflicts
1. Container start fails with port binding error
2. Error modal suggests stopping conflicting service
3. Provides command to find and stop conflicting process: `lsof -i :8080`

#### Container Startup Timeout
1. Startup modal shows for maximum 5 minutes
2. If timeout: Display error and stop container
3. Suggest checking Docker logs manually
4. Provide container ID for debugging

## 🛠️ Development & Contributing

### Development Environment Setup

```bash
# Prerequisites check
go version          # Should be 1.21+
docker --version    # Docker should be available

# Clone and setup
git clone https://github.com/khairu-aqsara/moodle-prototype-manager.git
cd moodle-prototype-manager
go mod tidy

# Start development
./dev.sh
```

### Project Architecture

```
moodle-prototype-manager/
├── main.go                 # Application entry point & window configuration
├── app.go                  # Main application logic & Wails bindings
├── wails.json             # Wails build configuration
├── go.mod/go.sum          # Go module dependencies
│
├── docker/                # Docker operations package
│   ├── manager.go         # Container lifecycle management
│   ├── health.go          # Docker/internet health checks
│   ├── logs.go            # Log parsing & credential extraction
│   ├── progress.go        # Download progress parsing
│   └── path.go            # Cross-platform path utilities
│
├── storage/               # File I/O and persistence
│   ├── files.go          # General file operations
│   └── credentials.go     # Credential management & parsing
│
├── utils/                 # Utility functions
│   ├── logger.go         # Logging utilities
│   ├── platform_*.go     # Platform-specific utilities
│
├── frontend/              # Web UI assets
│   ├── index.html        # Main UI structure
│   ├── css/styles.css    # Modern Moodle-themed styling
│   ├── js/
│   │   ├── app.js        # Core application logic
│   │   ├── ui.js         # UI manipulation functions
│   │   └── events.js     # Event handling & user interactions
│   ├── assets/images/    # Logo and image assets
│   └── dist/             # Built frontend files (generated)
│
├── build/                # Build output directory
│   └── bin/              # Final executables
│
├── scripts/              # Build and development scripts
│   ├── build.sh         # Automated build script
│   ├── dev.sh           # Development server script
│   ├── sign-mac.sh      # macOS code signing
│   └── sign-adhoc.sh    # Development signing
│
└── docs/                 # Documentation
    └── README.md         # This file
```

### Key Technologies

- **Backend**: Go 1.21+ with Wails v2 framework
- **Frontend**: Vanilla HTML5, CSS3, ES6+ JavaScript
- **Container Runtime**: Docker Engine via CLI commands
- **UI Framework**: Native WebView rendering
- **Build System**: Wails CLI with custom shell scripts
- **Logging**: Custom cross-platform logging utility

### API Architecture

#### Go Backend Functions (Exposed to Frontend)

```go
// Health monitoring
func (a *App) PerformHealthChecks() HealthCheckResult
func (a *App) CheckDockerHealth() bool
func (a *App) CheckInternetHealth() bool

// Container lifecycle
func (a *App) StartMoodleContainer() StartResult
func (a *App) StopMoodleContainer() StopResult
func (a *App) GetContainerStatus() StatusResult

// Browser integration
func (a *App) OpenInBrowser(url string) error

// Configuration
func (a *App) GetImageName() string
```

#### JavaScript Frontend Functions

```javascript
// Health monitoring
async function performHealthChecks()
async function updateHealthIndicators()

// Container management
async function startMoodleContainer()
async function stopMoodleContainer()

// UI state management
function showDownloadModal(progress)
function showStartupModal()
function showCredentials(data)
function hideCredentials()

// Modal management
function showModal(type, data)
function hideModal(type)
```

### Testing

```bash
# Run all Go tests
go test ./...

# Run tests with verbose output and coverage
go test -v -cover ./...

# Test specific packages
go test ./docker
go test ./storage

# Integration testing with Docker
docker --version && go test ./docker -integration
```

### Contributing Guidelines

1. **Fork & Clone**: Fork the repository and clone your fork
2. **Branch**: Create feature branch: `git checkout -b feature/your-feature-name`
3. **Code Standards**:
   - **Go**: Use `go fmt`, `go vet`, pass `go test ./...`
   - **JavaScript**: Follow ES6+ standards, use consistent indentation
   - **CSS**: Use semantic class names, maintain Moodle color scheme
4. **Testing**: Write tests for new functionality, ensure existing tests pass
5. **Documentation**: Update README and comments for significant changes
6. **Commit Messages**: Use clear, descriptive commit messages
7. **Pull Request**: Submit PR with detailed description of changes

## 🆘 Troubleshooting

### Common Issues & Solutions

#### "Docker status shows red circle"
**Symptoms**: Docker health indicator remains red, button disabled
**Causes**: Docker Desktop not running, permissions issues
**Solutions**:
```bash
# Check Docker installation
docker --version

# Start Docker Desktop (macOS)
open /Applications/Docker.app

# Start Docker Desktop (Windows)
# Use Start Menu or Desktop shortcut

# Check Docker daemon status
docker info

# Fix permissions (Linux/macOS)
sudo usermod -aG docker $USER
# Then log out and log back in
```

#### "Internet status shows red circle"
**Symptoms**: Internet health check fails, image downloads fail
**Causes**: Network connectivity, firewall, proxy issues
**Solutions**:
```bash
# Test Docker registry connectivity
curl -I https://registry-1.docker.io/v2/

# Test with proxy (if applicable)
curl -I -x http://proxy:port https://registry-1.docker.io/v2/

# Configure Docker proxy settings in Docker Desktop
# Settings → Resources → Proxies
```

#### "Port 8080 already in use"
**Symptoms**: Container fails to start with port binding error
**Causes**: Another service using port 8080
**Solutions**:
```bash
# Find process using port 8080
lsof -i :8080                    # macOS/Linux
netstat -ano | findstr :8080     # Windows

# Stop conflicting process
sudo kill -9 [PID]              # macOS/Linux
taskkill /PID [PID] /F          # Windows

# Or stop all Docker containers
docker stop $(docker ps -q)
```

#### "Container startup timeout"
**Symptoms**: Startup modal shows for 5+ minutes, then fails
**Causes**: Slow image download, insufficient resources, Docker issues
**Solutions**:
```bash
# Check Docker logs manually
docker logs [container_id]

# Verify system resources
docker system df
docker system info

# Pull image manually
docker pull wenkhairu/moodle-prototype:502-alpine

# Increase Docker memory allocation
# Docker Desktop → Settings → Resources → Memory → 4GB+
```

#### "Credentials not extracted"
**Symptoms**: Container starts but credentials don't appear
**Causes**: Log parsing issues, container initialization problems
**Solutions**:
```bash
# Check container logs manually
docker logs [container_id]

# Look for credential patterns
docker logs [container_id] | grep -i "password\|admin"

# Verify container is responding
curl -I http://localhost:8080

# Restart container through app
# Or manually: docker restart [container_id]
```

#### "App won't start" / "Code signature issues" (macOS)
**Symptoms**: App crashes on launch, security warnings
**Causes**: Gatekeeper restrictions, unsigned/corrupted app
**Solutions**:
```bash
# Bypass Gatekeeper (first launch)
# Right-click app → Open → Open

# Reset Gatekeeper (if issues persist)
sudo spctl --master-disable
sudo spctl --master-enable

# Check app signature
codesign -dv --verbose=4 /path/to/app

# Remove quarantine attribute
xattr -d com.apple.quarantine /path/to/app

# Re-download if corrupted
# Delete app and download fresh copy
```

### Performance Optimization

#### Slow Startup
```bash
# Optimize Docker Desktop settings
# → Settings → Resources → Advanced
# → Increase CPUs: 4+, Memory: 8GB+
# → Enable VirtioFS (macOS)

# Pre-pull Docker image
docker pull wenkhairu/moodle-prototype:502-alpine

# Clean Docker system
docker system prune -a
```

#### High Resource Usage
```bash
# Monitor resource usage
docker stats [container_id]

# Limit container resources (if needed)
# Note: This requires code modification
docker run --memory=1g --cpus=1 [image]

# Clean unused Docker resources
docker system prune
docker image prune
```

### Debug Mode & Logging

#### Enable Development Mode
```bash
# Run development version with browser debugging
./dev.sh

# Or run Wails with browser debugging
wails dev -browser
```

#### Access Application Logs

**macOS**:
```bash
# View Console.app logs
# Applications → Utilities → Console
# Filter for "moodle" or application name

# Command line access
log show --predicate 'process == "moodle-prototype-manager"'
```

**Windows**:
```bash
# Event Viewer logs
# Windows Key + R → eventvwr.msc
# Windows Logs → Application

# Or run from command prompt to see output
"C:\path\to\moodle-prototype-manager.exe"
```

#### Debug Container Issues
```bash
# Interactive container debugging
docker exec -it [container_id] /bin/bash

# View full container logs
docker logs --follow [container_id]

# Check container health
docker inspect [container_id]

# Test Moodle directly
curl -v http://localhost:8080
```

## 📊 System Requirements & Compatibility

### Supported Platforms

| Platform | Architecture | Minimum OS | Recommended OS |
|----------|-------------|-------------|----------------|
| macOS Intel | x86_64 | macOS 10.14 Mojave | macOS 12 Monterey+ |
| macOS Apple Silicon | arm64 | macOS 11 Big Sur | macOS 12 Monterey+ |
| Windows | x86_64 | Windows 10 1903 | Windows 11 |

### Resource Requirements

| Component | Minimum | Recommended | Notes |
|-----------|---------|-------------|-------|
| RAM | 4GB | 8GB | Docker needs 2-3GB, Moodle needs 1-2GB |
| Storage | 3GB free | 5GB free | Includes Docker image and temporary files |
| CPU | 2 cores | 4+ cores | Better performance with more cores |
| Network | Broadband | Broadband | Initial download ~500MB-1GB |

### Docker Compatibility

- **Docker Desktop**: Version 4.0+ recommended
- **Docker Engine**: Version 20.10+ required
- **WSL2**: Required on Windows (automatically configured by Docker Desktop)
- **Virtualization**: Must be enabled in BIOS/UEFI

## 🔗 Additional Resources

### Official Links
- **Moodle**: [moodle.org](https://moodle.org/) - Open source learning platform
- **Wails**: [wails.io](https://wails.io/) - Go + Web UI framework
- **Docker**: [docker.com](https://www.docker.com/) - Containerization platform

### Support Channels
- **GitHub Issues**: [Report bugs and request features](../../issues)
- **GitHub Discussions**: [Community support and questions](../../discussions)
- **Documentation**: [Additional docs and guides](docs/)

### Development Resources
- **Go Documentation**: [golang.org/doc](https://golang.org/doc/)
- **Wails Documentation**: [wails.io/docs](https://wails.io/docs/)
- **Docker Reference**: [docs.docker.com](https://docs.docker.com/)

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Wails Framework**: For providing excellent Go-to-Web tooling
- **Moodle Community**: For the open-source learning management system
- **Docker**: For containerization technology
- **Go Team**: For the robust programming language

---

**Version**: 1.0.1
**Last Updated**: January 2025
**Maintained by**: [Khairu Aqsara](mailto:wenkhairu@gmail.com)

For questions, issues, or contributions, please visit our [GitHub repository](../../).