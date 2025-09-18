#!/bin/bash

# Moodle Prototype Manager - Development Script
# This script prepares frontend files and runs the development server

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

echo -e "${BLUE}🚀 Moodle Prototype Manager - Development Mode${NC}"
echo -e "${BLUE}===============================================${NC}"

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

# Prepare frontend files for development
prepare_dev_frontend() {
    print_info "Preparing frontend files for development..."
    
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
    
    # Copy existing Wails bindings if they exist
    if [ -f "frontend/wailsjs/go/main/App.js" ]; then
        cp frontend/wailsjs/go/main/App.js frontend/dist/wailsjs/go/main/
        print_status "Copied existing Wails bindings"
    fi
    
    if [ -f "frontend/wailsjs/go/main/App.d.ts" ]; then
        cp frontend/wailsjs/go/main/App.d.ts frontend/dist/wailsjs/go/main/
        print_status "Copied TypeScript definitions"
    fi
    
    if [ -d "frontend/wailsjs/runtime" ]; then
        cp -r frontend/wailsjs/runtime/* frontend/dist/wailsjs/runtime/ 2>/dev/null || true
        print_status "Copied runtime files"
    fi
}

# Check image.docker file and set image if provided
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

# Run development server
run_dev_server() {
    print_info "Starting Wails development server..."
    export PATH=$PATH:$(go env GOPATH)/bin
    
    # Check if Wails CLI is available
    if ! command -v wails &> /dev/null; then
        print_error "Wails CLI not found. Installing..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
    fi
    
    echo
    print_info "Development server will start shortly..."
    print_info "The app will open automatically, or visit the browser URL shown below"
    echo
    
    # Start development server
    wails dev
}

# Show help
show_help() {
    echo -e "${BLUE}Moodle Prototype Manager - Development Script${NC}"
    echo
    echo -e "${YELLOW}Usage:${NC}"
    echo "  ./dev.sh [options]"
    echo
    echo -e "${YELLOW}Options:${NC}"
    echo "  -i, --image <image>  Set Docker image name (overrides image.docker file)"
    echo "  -c, --clean         Clean and prepare fresh frontend files"
    echo "  -h, --help          Show this help message"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  ./dev.sh                                        # Start dev server"
    echo "  ./dev.sh --clean                                # Clean and start dev server"
    echo "  ./dev.sh -i wenkhairu/moodle-prototype:503-beta # Start with specific image"
    echo "  ./dev.sh --image moodle:latest --clean          # Clean and start with custom image"
    echo
    echo -e "${YELLOW}What this script does:${NC}"
    echo "  1. 📁 Copies frontend source files to dist directory"
    echo "  2. 📄 Sets Docker image configuration"
    echo "  3. 🚀 Starts Wails development server with live reload"
    echo
    echo -e "${YELLOW}Development Tips:${NC}"
    echo "  • Edit files in frontend/ directory (not frontend/dist/)"
    echo "  • Run this script again after making changes to see them in the browser"
    echo "  • The development server will auto-reload when Go files change"
    echo "  • For frontend changes, you may need to restart this script"
    echo
    echo -e "${YELLOW}File Watching:${NC}"
    echo "  • Go files: Auto-reload (handled by Wails)"
    echo "  • Frontend files: Manual sync (run this script again)"
    echo
}

# Parse command line arguments
parse_arguments() {
    DOCKER_IMAGE=""
    CLEAN_MODE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--image)
                DOCKER_IMAGE="$2"
                shift 2
                ;;
            -c|--clean)
                CLEAN_MODE=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_warning "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Handle clean mode
handle_clean() {
    if [ "$CLEAN_MODE" = true ]; then
        print_info "Cleaning frontend dist directory..."
        if [ -d "frontend/dist" ]; then
            rm -rf frontend/dist
            print_status "Cleaned frontend/dist directory"
        fi
    fi
}

# Main execution
main() {
    echo -e "${BLUE}Starting development mode...${NC}"
    echo
    
    handle_clean
    check_image_config
    prepare_dev_frontend
    
    echo
    run_dev_server
}

# Parse arguments and run main function
parse_arguments "$@"
main