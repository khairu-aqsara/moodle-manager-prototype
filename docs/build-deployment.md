# Moodle Prototype Manager - Build and Deployment Guide

## Table of Contents

1. [Overview](#overview)
2. [Build Prerequisites](#build-prerequisites)
3. [Development Builds](#development-builds)
4. [Production Builds](#production-builds)
5. [Cross-Platform Compilation](#cross-platform-compilation)
6. [Code Signing](#code-signing)
7. [Distribution Packaging](#distribution-packaging)
8. [Release Management](#release-management)
9. [Deployment Strategies](#deployment-strategies)
10. [Continuous Integration](#continuous-integration)
11. [Troubleshooting Build Issues](#troubleshooting-build-issues)

## Overview

This guide covers the complete build and deployment process for the Moodle Prototype Manager application. The application uses Wails v2 framework for creating cross-platform desktop applications with Go backend and web frontend.

### Build Architecture

```
Source Code → Frontend Preparation → Go Build → Asset Embedding → Platform Packaging → Distribution
     ↓              ↓                  ↓           ↓                ↓                 ↓
Go + HTML/CSS/JS → Dist Directory → Binary → Embedded Assets → Platform Package → Release Artifact
```

### Supported Platforms

- **Windows**: x64 (AMD64) architecture
- **macOS**: Intel (AMD64) and Apple Silicon (ARM64)
- **Linux**: x64 (for development and testing)

## Build Prerequisites

### System Requirements

**Development Machine:**
- **Go**: Version 1.19+ (1.20+ recommended)
- **Node.js**: Version 16+ (for frontend tooling if needed)
- **Wails CLI**: Latest version
- **Platform-specific tools**: See platform sections below

### Wails Installation and Setup

**Install Wails CLI:**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify installation
wails version

# Check system prerequisites
wails doctor
```

**System Dependencies:**

**Windows:**
- Windows 10/11 SDK
- WebView2 runtime (automatically included)
- Visual Studio Build Tools or Visual Studio Community

**macOS:**
- Xcode Command Line Tools
- macOS SDK (included with Xcode)

**Linux (Development):**
- WebKitGTK development libraries
- pkg-config
- Build essentials (gcc, make)

```bash
# Ubuntu/Debian
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev pkg-config build-essential

# Fedora/RHEL
sudo dnf install gtk3-devel webkit2gtk3-devel pkgconf-pkg-config gcc-c++
```

### Project Dependencies

**Go Dependencies:**
```bash
# Install and verify dependencies
go mod download
go mod verify

# Check for updates
go list -u -m all

# Clean up unused dependencies
go mod tidy
```

**Development Tools:**
```bash
# Code quality tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest
```

## Development Builds

### Quick Development Build

**Using Development Script:**
```bash
# Start development server with live reload
./dev.sh

# Build for quick testing
./dev.sh --clean
```

**Manual Development Build:**
```bash
# Prepare frontend files
mkdir -p frontend/dist
cp -r frontend/* frontend/dist/

# Build for current platform
wails build -debug

# Output will be in build/bin/ directory
```

### Development Build Configuration

**Debug Features:**
- Enhanced logging output
- Developer tools enabled
- Unminified assets
- Source maps included
- Hot reload capability

**Environment Detection:**
```go
// Build-time flags
var (
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
)

// Runtime environment detection
func isDevelopment() bool {
    return Version == "dev" ||
           strings.Contains(os.Args[0], "dev") ||
           fileExists("go.mod")
}
```

### Testing Development Builds

**Automated Testing:**
```bash
# Run tests before building
go test -v -race ./...

# Build and test
wails build -debug
./build/bin/moodle-prototype-manager --test

# Integration testing with Docker
docker --version && ./build/bin/moodle-prototype-manager --health-check
```

## Production Builds

### Frontend Preparation

**Production Frontend Build:**
```bash
# Clean previous builds
rm -rf frontend/dist

# Create distribution directory
mkdir -p frontend/dist/{js,css,assets/images}

# Copy and optimize frontend files
cp frontend/index.html frontend/dist/
cp frontend/js/*.js frontend/dist/js/
cp frontend/css/*.css frontend/dist/css/
cp -r frontend/assets/* frontend/dist/assets/

# Optional: Minify assets (if using build tools)
# npx terser frontend/js/*.js --compress --mangle -o frontend/dist/js/app.min.js
```

### Production Build Process

**Single Platform Build:**
```bash
# Clean build environment
wails clean

# Build for production
wails build -clean -s -trimpath

# Build with custom flags
wails build \
    -clean \
    -s \
    -trimpath \
    -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -tags "production"
```

**Build Flags Explanation:**
- `-clean`: Clean build cache before building
- `-s`: Strip symbol table and debug information
- `-trimpath`: Remove absolute file paths from executable
- `-ldflags`: Pass flags to the Go linker
- `-tags`: Build tags for conditional compilation

### Production Configuration

**Build-time Variables:**
```go
// Set via ldflags during build
var (
    Version   = "1.0.0"
    BuildTime = "2024-01-15T10:30:00Z"
    GitCommit = "abc123def456"
    BuildType = "production"
)

func init() {
    // Configure production settings
    if BuildType == "production" {
        setupProductionLogging()
        disableDevFeatures()
    }
}
```

**Production Build Verification:**
```bash
# Verify binary
file build/bin/moodle-prototype-manager*

# Check embedded assets
strings build/bin/moodle-prototype-manager | grep -i version

# Test basic functionality
./build/bin/moodle-prototype-manager --version
./build/bin/moodle-prototype-manager --health-check
```

## Cross-Platform Compilation

### Multi-Platform Build Script

**Build Script (`build.sh`):**
```bash
#!/bin/bash
set -e

# Configuration
APP_NAME="Moodle Prototype Manager"
VERSION=${VERSION:-$(git describe --tags --always)}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse HEAD)

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -s -w"

echo "Building ${APP_NAME} version ${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "Git commit: ${GIT_COMMIT}"

# Clean previous builds
rm -rf build/
mkdir -p build/dist

# Prepare frontend
echo "Preparing frontend..."
./dev.sh --clean > /dev/null 2>&1

# Build for multiple platforms
echo "Building for multiple platforms..."
wails build \
    -platform darwin/amd64,darwin/arm64,windows/amd64 \
    -clean \
    -s \
    -trimpath \
    -ldflags "${LDFLAGS}"

echo "Build completed successfully!"
echo "Artifacts:"
find build/ -name "${APP_NAME}*" -type f | sed 's/^/  /'
```

### Platform-Specific Builds

**macOS Builds:**
```bash
# Intel Macs
wails build -platform darwin/amd64 -clean

# Apple Silicon Macs
wails build -platform darwin/arm64 -clean

# Universal Binary (both architectures)
wails build -platform darwin/universal -clean
```

**Windows Builds:**
```bash
# 64-bit Windows
wails build -platform windows/amd64 -clean

# With custom icon and manifest
wails build \
    -platform windows/amd64 \
    -clean \
    -windowsconsole
```

**Linux Builds (Development/Testing):**
```bash
# 64-bit Linux
wails build -platform linux/amd64 -clean

# ARM64 Linux
wails build -platform linux/arm64 -clean
```

### Build Output Structure

**Generated Files:**
```
build/
├── bin/
│   ├── Moodle Prototype Manager           # macOS Intel
│   ├── Moodle Prototype Manager (ARM64)   # macOS Apple Silicon
│   └── Moodle Prototype Manager.exe       # Windows x64
└── darwin/
    └── Moodle Prototype Manager.app/      # macOS app bundle
        ├── Contents/
        │   ├── Info.plist
        │   ├── MacOS/
        │   │   └── Moodle Prototype Manager
        │   └── Resources/
        └── ...
```

## Code Signing

### macOS Code Signing

**Requirements:**
- Apple Developer Account
- Developer ID Application Certificate
- macOS development environment

**Signing Process:**
```bash
# Install certificate in Keychain
# (Done through Xcode or Keychain Access)

# Sign the application
codesign --force --options runtime \
    --sign "Developer ID Application: Your Name (TEAM_ID)" \
    "build/darwin/Moodle Prototype Manager.app"

# Verify signature
codesign -v -v "build/darwin/Moodle Prototype Manager.app"
spctl -a -v "build/darwin/Moodle Prototype Manager.app"
```

**Automated Signing Script (`sign-mac.sh`):**
```bash
#!/bin/bash
set -e

APP_PATH="build/darwin/Moodle Prototype Manager.app"
CERT_NAME="Developer ID Application: Your Name (TEAM_ID)"

if [ ! -d "$APP_PATH" ]; then
    echo "Error: App bundle not found at $APP_PATH"
    exit 1
fi

echo "Signing macOS application..."

# Sign all nested frameworks and binaries first
find "$APP_PATH" -name "*.dylib" -exec codesign --force --options runtime --sign "$CERT_NAME" {} \;
find "$APP_PATH" -name "*.framework" -exec codesign --force --options runtime --sign "$CERT_NAME" {} \;

# Sign the main application
codesign --force --options runtime --sign "$CERT_NAME" "$APP_PATH"

# Verify
codesign -v -v "$APP_PATH"
echo "✓ Application signed successfully"

# Create signed DMG (optional)
create-dmg \
    --volname "Moodle Prototype Manager" \
    --background "assets/dmg-background.png" \
    --window-pos 200 120 \
    --window-size 600 400 \
    --icon-size 100 \
    --icon "Moodle Prototype Manager.app" 200 190 \
    --hide-extension "Moodle Prototype Manager.app" \
    --app-drop-link 400 190 \
    "build/Moodle Prototype Manager.dmg" \
    "$APP_PATH"

# Sign the DMG
codesign --force --sign "$CERT_NAME" "build/Moodle Prototype Manager.dmg"
```

### Windows Code Signing

**Requirements:**
- Code signing certificate (EV or Standard)
- SignTool utility (Windows SDK)

**Signing Process:**
```powershell
# Using SignTool
signtool sign /f "certificate.p12" /p "password" /t "http://timestamp.digicert.com" "Moodle Prototype Manager.exe"

# Verify signature
signtool verify /pa "Moodle Prototype Manager.exe"
```

**Automated Signing Script (`sign-windows.ps1`):**
```powershell
param(
    [Parameter(Mandatory=$true)]
    [string]$CertificatePath,

    [Parameter(Mandatory=$true)]
    [string]$CertificatePassword
)

$AppPath = "build\bin\Moodle Prototype Manager.exe"
$TimestampUrl = "http://timestamp.digicert.com"

if (-Not (Test-Path $AppPath)) {
    Write-Error "Application not found at $AppPath"
    exit 1
}

Write-Host "Signing Windows application..."

# Sign the executable
& signtool sign /f $CertificatePath /p $CertificatePassword /t $TimestampUrl $AppPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Application signed successfully"

    # Verify signature
    & signtool verify /pa $AppPath
} else {
    Write-Error "Signing failed"
    exit 1
}
```

## Distribution Packaging

### macOS Distribution

**DMG Creation:**
```bash
# Create DMG with create-dmg tool
create-dmg \
    --volname "Moodle Prototype Manager v1.0.0" \
    --background "assets/dmg-background.png" \
    --window-pos 200 120 \
    --window-size 600 400 \
    --icon-size 100 \
    --icon "Moodle Prototype Manager.app" 200 190 \
    --hide-extension "Moodle Prototype Manager.app" \
    --app-drop-link 400 190 \
    "dist/Moodle Prototype Manager v1.0.0.dmg" \
    "build/darwin/Moodle Prototype Manager.app"
```

**App Store Package (if needed):**
```bash
# Create App Store package
productbuild --component "build/darwin/Moodle Prototype Manager.app" /Applications \
    "dist/Moodle Prototype Manager.pkg"
```

### Windows Distribution

**Portable Executable:**
```bash
# Simple ZIP distribution
cd build/bin
zip -r "../../dist/Moodle Prototype Manager v1.0.0 Windows x64.zip" \
    "Moodle Prototype Manager.exe"
```

**Installer Creation with NSIS:**

**NSIS Script (`installer.nsi`):**
```nsis
!define APPNAME "Moodle Prototype Manager"
!define APPVERSION "1.0.0"
!define DESCRIPTION "Desktop application for managing Moodle prototype Docker containers"

Name "${APPNAME}"
OutFile "dist\Moodle Prototype Manager Setup v${APPVERSION}.exe"
InstallDir "$PROGRAMFILES64\${APPNAME}"

Page directory
Page instfiles

Section "MainSection" SEC01
    SetOutPath "$INSTDIR"
    File "build\bin\Moodle Prototype Manager.exe"

    # Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"

    # Create start menu shortcuts
    CreateDirectory "$SMPROGRAMS\${APPNAME}"
    CreateShortCut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\Moodle Prototype Manager.exe"
    CreateShortCut "$SMPROGRAMS\${APPNAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe"

    # Create desktop shortcut
    CreateShortCut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\Moodle Prototype Manager.exe"
SectionEnd

Section "Uninstall"
    Delete "$INSTDIR\Moodle Prototype Manager.exe"
    Delete "$INSTDIR\Uninstall.exe"

    Delete "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk"
    Delete "$SMPROGRAMS\${APPNAME}\Uninstall.lnk"
    RMDir "$SMPROGRAMS\${APPNAME}"

    Delete "$DESKTOP\${APPNAME}.lnk"
    RMDir "$INSTDIR"
SectionEnd
```

**Build Installer:**
```bash
# Install NSIS (if not already installed)
# Build installer
makensis installer.nsi
```

### Linux Distribution (Development)

**AppImage Creation:**
```bash
# Create AppDir structure
mkdir -p dist/appdir/usr/bin
mkdir -p dist/appdir/usr/share/applications
mkdir -p dist/appdir/usr/share/icons

# Copy application
cp build/bin/moodle-prototype-manager dist/appdir/usr/bin/

# Create desktop file
cat > dist/appdir/usr/share/applications/moodle-prototype-manager.desktop << EOF
[Desktop Entry]
Type=Application
Name=Moodle Prototype Manager
Exec=moodle-prototype-manager
Icon=moodle-prototype-manager
Categories=Development;
EOF

# Copy icon
cp assets/icon.png dist/appdir/usr/share/icons/moodle-prototype-manager.png

# Create AppImage
appimagetool dist/appdir dist/Moodle\ Prototype\ Manager-x86_64.AppImage
```

## Release Management

### Version Management

**Semantic Versioning:**
- **Major (1.0.0)**: Breaking changes or major new features
- **Minor (1.1.0)**: New features, backwards compatible
- **Patch (1.1.1)**: Bug fixes, backwards compatible

**Version Configuration:**
```go
// version.go
package main

var (
    // Set via ldflags during build
    Version   = "dev"
    BuildTime = "unknown"
    GitCommit = "unknown"
    BuildType = "development"
)

func GetVersionInfo() VersionInfo {
    return VersionInfo{
        Version:   Version,
        BuildTime: BuildTime,
        GitCommit: GitCommit,
        BuildType: BuildType,
        GoVersion: runtime.Version(),
        Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
    }
}
```

### Release Process

**Pre-Release Checklist:**
- [ ] All tests pass on all platforms
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated
- [ ] Version numbers are bumped
- [ ] Code is signed (macOS/Windows)
- [ ] Distribution packages created
- [ ] Release notes prepared

**Release Script (`release.sh`):**
```bash
#!/bin/bash
set -e

# Configuration
VERSION=${1:-$(git describe --tags --always)}
RELEASE_DIR="release/v${VERSION}"

echo "Creating release v${VERSION}"

# Create release directory
mkdir -p "${RELEASE_DIR}"

# Build for all platforms
./build.sh

# Sign applications
if [[ "$OSTYPE" == "darwin"* ]]; then
    ./sign-mac.sh
fi

# Create distribution packages
echo "Creating distribution packages..."

# macOS DMG
create-dmg \
    --volname "Moodle Prototype Manager v${VERSION}" \
    "dist/Moodle Prototype Manager v${VERSION}.dmg" \
    "build/darwin/Moodle Prototype Manager.app"

# Windows ZIP
cd build/bin
zip -r "../../${RELEASE_DIR}/Moodle Prototype Manager v${VERSION} Windows x64.zip" \
    "Moodle Prototype Manager.exe"
cd ../..

# Copy macOS DMG
cp "dist/Moodle Prototype Manager v${VERSION}.dmg" "${RELEASE_DIR}/"

# Generate checksums
cd "${RELEASE_DIR}"
for file in *; do
    if [[ -f "$file" ]]; then
        shasum -a 256 "$file" >> "checksums.txt"
    fi
done
cd ../..

echo "Release artifacts created in ${RELEASE_DIR}/"
ls -la "${RELEASE_DIR}/"
```

### GitHub Releases

**Release Creation:**
```bash
# Using GitHub CLI
gh release create v1.0.0 \
    release/v1.0.0/* \
    --title "Moodle Prototype Manager v1.0.0" \
    --notes-file RELEASE_NOTES.md \
    --draft
```

**Release Notes Template:**
```markdown
# Moodle Prototype Manager v1.0.0

## New Features
- Added support for custom Docker images
- Improved progress tracking for image downloads
- Enhanced error handling and user feedback

## Bug Fixes
- Fixed container restart issues on macOS
- Resolved file permission issues in production builds
- Improved Docker daemon connectivity detection

## Breaking Changes
- Configuration file format updated (automatic migration provided)

## Platform Support
- Windows 10/11 (x64)
- macOS 10.15+ (Intel and Apple Silicon)
- Docker Desktop 4.0+ required

## Installation
Download the appropriate package for your platform:
- **macOS**: Download the .dmg file and drag to Applications
- **Windows**: Download and run the .exe installer or use the portable .zip

## Checksums
See `checksums.txt` for SHA-256 verification hashes.
```

## Deployment Strategies

### Direct Download Distribution

**Website Deployment:**
- Host distribution files on CDN or web server
- Provide platform-specific download links
- Include checksums for verification
- Maintain download statistics

**Auto-Update System (Future Enhancement):**
```go
type UpdateChecker struct {
    currentVersion string
    updateURL     string
}

func (uc *UpdateChecker) CheckForUpdates() (*UpdateInfo, error) {
    // Implementation for checking updates
    resp, err := http.Get(uc.updateURL + "/latest")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse response and compare versions
    // Return update information if available
}
```

### Enterprise Distribution

**Internal Distribution:**
- Host on internal servers
- Use configuration management tools
- Integrate with IT deployment systems
- Provide silent installation options

### App Store Distribution (Future)

**macOS App Store:**
- Sandbox requirements
- App Store review process
- Different build configuration
- In-app purchase integration (if needed)

## Continuous Integration

### GitHub Actions Pipeline

**.github/workflows/build.yml:**
```yaml
name: Build and Release

on:
  push:
    tags: ['v*']
  pull_request:
    branches: [main]

jobs:
  build:
    strategy:
      matrix:
        platform: [macos-latest, windows-latest, ubuntu-latest]

    runs-on: ${{ matrix.platform }}

    steps:
    - uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.20'

    - name: Install Wails
      run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v ./...

    - name: Build application
      run: wails build -clean

    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: app-${{ matrix.platform }}
        path: build/bin/

  release:
    if: startsWith(github.ref, 'refs/tags/v')
    needs: build
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Download artifacts
      uses: actions/download-artifact@v3

    - name: Create release
      uses: softprops/action-gh-release@v1
      with:
        files: |
          app-*/Moodle Prototype Manager*
        draft: true
        generate_release_notes: true
```

### Build Automation

**Automated Build Triggers:**
- Git tag creation (`v*` pattern)
- Pull request validation
- Scheduled nightly builds
- Manual workflow dispatch

**Build Matrix:**
- Multiple Go versions
- Multiple platforms
- Different build configurations
- Integration tests with Docker

## Troubleshooting Build Issues

### Common Build Problems

**Wails Build Failures:**
```bash
# Clear Wails cache
wails clean

# Check system dependencies
wails doctor

# Update Wails to latest version
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build with verbose output
wails build -debug -v 2
```

**Frontend Asset Issues:**
```bash
# Verify frontend structure
ls -la frontend/dist/

# Check asset paths in HTML
grep -r "src=" frontend/dist/
grep -r "href=" frontend/dist/

# Rebuild frontend
./dev.sh --clean
```

**Cross-Platform Build Issues:**
```bash
# Check available platforms
go tool dist list

# Build for specific platform
GOOS=windows GOARCH=amd64 go build .

# Check CGO dependencies
CGO_ENABLED=0 go build .
```

### Platform-Specific Issues

**macOS Issues:**
```bash
# Check Xcode command line tools
xcode-select --install

# Verify developer certificate
security find-identity -v -p codesigning

# Check notarization status
xcrun altool --notarization-history 0 -u "your-apple-id"
```

**Windows Issues:**
```powershell
# Check Windows SDK
where signtool.exe

# Verify Visual Studio Build Tools
where cl.exe

# Test WebView2 availability
reg query "HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"
```

### Debug Build Output

**Build Logging:**
```bash
# Enable verbose build output
export WAILS_LOG_LEVEL=debug

# Build with detailed output
wails build -debug -v 2 2>&1 | tee build.log

# Analyze build log
grep -i error build.log
grep -i warning build.log
```

**Binary Analysis:**
```bash
# Check binary dependencies (macOS/Linux)
otool -L "Moodle Prototype Manager"  # macOS
ldd moodle-prototype-manager          # Linux

# Check binary size
ls -lh build/bin/

# Analyze binary symbols
nm "Moodle Prototype Manager" | grep -i version
```

This build and deployment guide provides comprehensive coverage of all aspects of building, packaging, and distributing the Moodle Prototype Manager application across multiple platforms.