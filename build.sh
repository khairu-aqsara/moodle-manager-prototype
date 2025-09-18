#!/bin/bash

# Moodle Prototype Manager - Build Script
# This script automates the build process and ensures all necessary files are copied

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}🏗️  Moodle Prototype Manager - Build Script${NC}"
echo -e "${BLUE}================================================${NC}"

# Function to print status
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if Wails CLI is available
check_wails() {
    print_info "Checking Wails CLI availability..."
    export PATH=$PATH:$(go env GOPATH)/bin
    
    if ! command -v wails &> /dev/null; then
        print_error "Wails CLI not found. Installing..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
        if ! command -v wails &> /dev/null; then
            print_error "Failed to install Wails CLI. Please install manually:"
            echo "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
            exit 1
        fi
    fi
    print_status "Wails CLI is available"
}

# Clean previous build
clean_build() {
    print_info "Cleaning previous build..."
    if [ -d "build" ]; then
        rm -rf build
        print_status "Cleaned build directory"
    fi
    
    if [ -d "frontend/dist" ]; then
        rm -rf frontend/dist
        print_status "Cleaned frontend/dist directory"
    fi
}

# Prepare frontend dist directory
prepare_frontend() {
    print_info "Preparing frontend files..."
    
    # Create dist directory structure
    mkdir -p frontend/dist
    mkdir -p frontend/dist/js
    mkdir -p frontend/dist/css
    mkdir -p frontend/dist/assets/images
    mkdir -p frontend/dist/wailsjs/go/main
    mkdir -p frontend/dist/wailsjs/runtime
    
    # Copy HTML files
    if [ -f "frontend/index.html" ]; then
        cp frontend/index.html frontend/dist/
        print_status "Copied index.html"
    else
        print_error "frontend/index.html not found"
        exit 1
    fi
    
    # Copy JavaScript files
    for js_file in frontend/js/*.js; do
        if [ -f "$js_file" ]; then
            cp "$js_file" frontend/dist/js/
            filename=$(basename "$js_file")
            print_status "Copied $filename"
        fi
    done
    
    # Copy CSS files
    for css_file in frontend/css/*.css; do
        if [ -f "$css_file" ]; then
            cp "$css_file" frontend/dist/css/
            filename=$(basename "$css_file")
            print_status "Copied $filename"
        fi
    done
    
    # Copy assets
    if [ -d "frontend/assets" ]; then
        cp -r frontend/assets/* frontend/dist/assets/ 2>/dev/null || true
        print_status "Copied assets"
    fi
    
    # Copy Wails bindings if they exist
    if [ -f "frontend/wailsjs/go/main/App.js" ]; then
        cp frontend/wailsjs/go/main/App.js frontend/dist/wailsjs/go/main/
        print_status "Copied Wails Go bindings"
    fi
    
    if [ -f "frontend/wailsjs/go/main/App.d.ts" ]; then
        cp frontend/wailsjs/go/main/App.d.ts frontend/dist/wailsjs/go/main/
        print_status "Copied Wails TypeScript definitions"
    fi
    
    if [ -d "frontend/wailsjs/runtime" ]; then
        cp -r frontend/wailsjs/runtime/* frontend/dist/wailsjs/runtime/ 2>/dev/null || true
        print_status "Copied Wails runtime files"
    fi
}

# Ensure image.docker file exists and set image if provided
check_image_config() {
    print_info "Checking Docker image configuration..."
    
    # Check if image name was provided as parameter
    if [ ! -z "$DOCKER_IMAGE" ]; then
        echo "$DOCKER_IMAGE" > image.docker
        print_status "Set Docker image to: $DOCKER_IMAGE"
    elif [ ! -f "image.docker" ]; then
        print_warning "image.docker not found, creating default..."
        echo "wenkhairu/moodle-prototype:502-stable" > image.docker
        print_status "Created default image.docker"
    else
        IMAGE_NAME=$(cat image.docker | tr -d '\n\r')
        print_status "Using existing Docker image: $IMAGE_NAME"
    fi
    
    # Display current image
    CURRENT_IMAGE=$(cat image.docker | tr -d '\n\r')
    echo -e "${BLUE}🐳 Docker image:${NC} ${GREEN}$CURRENT_IMAGE${NC}"
}

# Install Go dependencies
install_dependencies() {
    print_info "Installing Go dependencies..."
    go mod tidy
    print_status "Go dependencies installed"
}

# Generate Wails bindings
generate_bindings() {
    print_info "Generating Wails bindings..."
    export PATH=$PATH:$(go env GOPATH)/bin
    
    # Generate bindings
    wails generate module
    
    # Copy the generated bindings to dist
    if [ -f "frontend/wailsjs/go/main/App.js" ]; then
        cp frontend/wailsjs/go/main/App.js frontend/dist/wailsjs/go/main/
        print_status "Updated Wails bindings in dist"
    fi
    
    if [ -f "frontend/wailsjs/go/main/App.d.ts" ]; then
        cp frontend/wailsjs/go/main/App.d.ts frontend/dist/wailsjs/go/main/
        print_status "Updated TypeScript definitions in dist"
    fi
}

# Build application
build_app() {
    print_info "Building application..."
    export PATH=$PATH:$(go env GOPATH)/bin
    
    # Parse command line arguments for build options
    BUILD_ARGS=""
    
    # Check for platform argument
    if [ ! -z "$1" ]; then
        case "$1" in
            "darwin-amd64")
                BUILD_ARGS="-platform darwin/amd64"
                print_info "Building for macOS Intel (amd64)"
                ;;
            "darwin-arm64")
                BUILD_ARGS="-platform darwin/arm64"
                print_info "Building for macOS Apple Silicon (arm64)"
                ;;
            "windows-amd64")
                BUILD_ARGS="-platform windows/amd64"
                print_info "Building for Windows 64-bit"
                ;;
            "all")
                BUILD_ARGS="-platform darwin/amd64,darwin/arm64,windows/amd64"
                print_info "Building for all platforms"
                ;;
            "current")
                print_info "Building for current platform"
                ;;
            *)
                print_warning "Unknown platform: $1. Building for current platform."
                ;;
        esac
    fi
    
    # Run Wails build
    if [ -z "$BUILD_ARGS" ]; then
        wails build
    else
        wails build $BUILD_ARGS
    fi
    
    print_status "Application built successfully"
}

# Copy image.docker to build directory
copy_config_to_build() {
    print_info "Copying configuration files to build..."
    
    # For macOS app bundle
    if [ -d "build/bin/moodle-prototype-manager.app/Contents/MacOS" ]; then
        cp image.docker build/bin/moodle-prototype-manager.app/Contents/MacOS/
        print_status "Copied image.docker to macOS app bundle"
    fi
    
    # For Windows executable directory
    if [ -f "build/bin/Moodle Prototype Manager.exe" ]; then
        cp image.docker build/bin/
        print_status "Copied image.docker to Windows build directory"
    elif [ -f "build/bin/moodle-prototype-manager.exe" ]; then
        cp image.docker build/bin/
        print_status "Copied image.docker to Windows build directory"
    fi
    
    # For any other binaries
    if [ -d "build/bin" ]; then
        find build/bin -name "*moodle-prototype-manager*" -type f -perm +111 | while read -r binary; do
            binary_dir=$(dirname "$binary")
            if [ ! -f "$binary_dir/image.docker" ]; then
                cp image.docker "$binary_dir/"
                print_status "Copied image.docker to $(basename "$binary_dir")"
            fi
        done
    fi
}

# Create Windows distribution package
create_windows_distribution() {
    print_info "Creating Windows distribution package..."
    
    # Check if Windows executable exists (handle actual Wails output filename)
    WINDOWS_EXE=""
    if [ -f "build/bin/Moodle Prototype Manager.exe" ]; then
        WINDOWS_EXE="build/bin/Moodle Prototype Manager.exe"
    elif [ -f "build/bin/moodle-prototype-manager.exe" ]; then
        WINDOWS_EXE="build/bin/moodle-prototype-manager.exe"
    fi
    
    if [ ! -z "$WINDOWS_EXE" ]; then
        # Create temporary distribution directory
        TEMP_DIR="build/windows-dist-temp"
        rm -rf "$TEMP_DIR" 2>/dev/null || true
        mkdir -p "$TEMP_DIR"
        
        # Copy executable and configuration
        cp "$WINDOWS_EXE" "$TEMP_DIR/moodle-prototype-manager.exe"
        cp "image.docker" "$TEMP_DIR/"
        
        # Create zip file
        ZIP_FILE="build/moodle-prototype-manager-windows-amd64.zip"
        rm -f "$ZIP_FILE" 2>/dev/null || true
        
        # Check if zip command is available
        if command -v zip &> /dev/null; then
            cd "$TEMP_DIR"
            zip -q "../moodle-prototype-manager-windows-amd64.zip" moodle-prototype-manager.exe image.docker
            cd - > /dev/null
            
            if [ -f "$ZIP_FILE" ]; then
                print_status "Created Windows distribution: $ZIP_FILE"
            else
                print_error "Failed to create Windows distribution zip"
            fi
        else
            print_warning "zip command not available - Windows distribution not created"
            print_info "Windows executable and config are available in: build/bin/"
        fi
        
        # Clean up temporary directory
        rm -rf "$TEMP_DIR"
    else
        print_info "No Windows executable found - skipping Windows distribution"
    fi
}

# Display build results
show_results() {
    print_info "Build completed! Results:"
    echo
    
    if [ -d "build/bin" ]; then
        echo -e "${BLUE}📦 Built applications:${NC}"
        find build/bin -name "*moodle-prototype-manager*" | while read -r file; do
            if [ -f "$file" ]; then
                size=$(du -h "$file" | cut -f1)
                echo -e "  ${GREEN}▶${NC} $file (${size})"
            elif [ -d "$file" ]; then
                echo -e "  ${GREEN}▶${NC} $file (app bundle)"
            fi
        done
        echo
    fi
    
    # Show Windows distribution package if it exists
    if [ -f "build/moodle-prototype-manager-windows-amd64.zip" ]; then
        ZIP_SIZE=$(du -h "build/moodle-prototype-manager-windows-amd64.zip" | cut -f1)
        echo -e "${BLUE}📦 Windows distribution package:${NC}"
        echo -e "  ${GREEN}▶${NC} build/moodle-prototype-manager-windows-amd64.zip (${ZIP_SIZE})"
        echo -e "  ${BLUE}Contains:${NC} moodle-prototype-manager.exe + image.docker"
        echo
    fi
    
    # Show current Docker image
    if [ -f "image.docker" ]; then
        IMAGE_NAME=$(cat image.docker | tr -d '\n\r')
        echo -e "${BLUE}🐳 Docker image configured:${NC} ${GREEN}$IMAGE_NAME${NC}"
        echo
    fi
    
    echo -e "${GREEN}🎉 Build process completed successfully!${NC}"
    echo
    echo -e "${BLUE}📋 To run the application:${NC}"
    if [ -d "build/bin/moodle-prototype-manager.app" ]; then
        echo -e "  ${YELLOW}macOS:${NC} open build/bin/moodle-prototype-manager.app"
    fi
    if [ -f "build/moodle-prototype-manager-windows-amd64.zip" ]; then
        echo -e "  ${YELLOW}Windows:${NC} Extract build/moodle-prototype-manager-windows-amd64.zip and run moodle-prototype-manager.exe"
    elif [ -f "build/bin/Moodle Prototype Manager.exe" ]; then
        echo -e "  ${YELLOW}Windows:${NC} .\\build\\bin\\\"Moodle Prototype Manager.exe\""
    elif [ -f "build/bin/moodle-prototype-manager.exe" ]; then
        echo -e "  ${YELLOW}Windows:${NC} .\\build\\bin\\moodle-prototype-manager.exe"
    fi
    echo
    echo -e "${BLUE}🔧 To change Docker image:${NC}"
    echo -e "  1. Edit the ${YELLOW}image.docker${NC} file"
    echo -e "  2. Run this build script again"
    echo
}

# Parse command line arguments
parse_arguments() {
    PLATFORM=""
    DOCKER_IMAGE=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--image)
                DOCKER_IMAGE="$2"
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            darwin-amd64|darwin-arm64|windows-amd64|all|current)
                PLATFORM="$1"
                shift
                ;;
            *)
                if [ -z "$PLATFORM" ]; then
                    PLATFORM="$1"
                fi
                shift
                ;;
        esac
    done
}

# Main execution
main() {
    echo -e "${BLUE}Starting build process...${NC}"
    echo
    
    # Perform all build steps
    check_wails
    clean_build
    check_image_config
    install_dependencies
    prepare_frontend
    generate_bindings
    build_app "$PLATFORM"
    copy_config_to_build
    create_windows_distribution
    
    echo
    show_results
}

# Help function
show_help() {
    echo -e "${BLUE}Moodle Prototype Manager - Build Script${NC}"
    echo
    echo -e "${YELLOW}Usage:${NC}"
    echo "  ./build.sh [options] [platform]"
    echo
    echo -e "${YELLOW}Options:${NC}"
    echo "  -i, --image <image>  Set Docker image name (overrides image.docker file)"
    echo "  -h, --help          Show this help message"
    echo
    echo -e "${YELLOW}Platforms:${NC}"
    echo "  current         Build for current platform (default)"
    echo "  darwin-amd64    Build for macOS Intel"
    echo "  darwin-arm64    Build for macOS Apple Silicon"
    echo "  windows-amd64   Build for Windows 64-bit"
    echo "  all             Build for all platforms"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  ./build.sh                                        # Build for current platform"
    echo "  ./build.sh darwin-arm64                           # Build for macOS Apple Silicon"
    echo "  ./build.sh all                                    # Build for all platforms"
    echo "  ./build.sh -i wenkhairu/moodle-prototype:503-beta # Build with specific image"
    echo "  ./build.sh --image moodle:latest darwin-amd64     # Build with custom image for macOS Intel"
    echo
    echo -e "${YELLOW}What this script does:${NC}"
    echo "  1. ✅ Checks and installs Wails CLI if needed"
    echo "  2. 🧹 Cleans previous build artifacts"
    echo "  3. 📄 Sets Docker image configuration"
    echo "  4. 📦 Installs Go dependencies"
    echo "  5. 📁 Prepares frontend files in dist directory"
    echo "  6. 🔗 Generates Wails bindings"
    echo "  7. 🏗️  Builds the application"
    echo "  8. 📋 Copies configuration files to build"
    echo "  9. 📊 Shows build results"
}

# Handle help argument
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# Parse arguments and run main function
parse_arguments "$@"
main "$PLATFORM"