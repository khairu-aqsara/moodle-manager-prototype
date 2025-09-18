#!/bin/bash

# Moodle Prototype Manager - macOS Code Signing Script
# This script signs the macOS application after building

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

echo -e "${BLUE}🔏 Moodle Prototype Manager - macOS Code Signing${NC}"
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

# Check if running on macOS
check_platform() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        print_error "This script must be run on macOS"
        exit 1
    fi
    print_status "Running on macOS"
}

# Check if app exists
check_app_exists() {
    if [ ! -d "build/bin/moodle-prototype-manager.app" ]; then
        print_error "Application not found at build/bin/moodle-prototype-manager.app"
        print_info "Please run ./build.sh first to build the application"
        exit 1
    fi
    print_status "Found application bundle"
}

# List available signing identities
list_identities() {
    print_info "Available signing identities:"
    echo
    security find-identity -v -p codesigning | grep -E "Developer ID Application|Apple Development|Mac Developer" || {
        print_warning "No valid signing identities found"
        echo
        echo -e "${YELLOW}To sign your application, you need:${NC}"
        echo "  1. An Apple Developer account"
        echo "  2. A valid signing certificate installed in Keychain"
        echo
        echo -e "${YELLOW}Options:${NC}"
        echo "  - For distribution outside App Store: 'Developer ID Application' certificate"
        echo "  - For development/testing: 'Apple Development' or 'Mac Developer' certificate"
        echo "  - For App Store: 'Apple Distribution' certificate"
        echo
        return 1
    }
    echo
}

# Sign the application
sign_app() {
    local identity="$1"
    local app_path="build/bin/moodle-prototype-manager.app"
    
    print_info "Signing application with identity: $identity"
    
    # Sign with hardened runtime, timestamp, and entitlements if needed
    codesign --force --deep --sign "$identity" \
        --options runtime \
        --timestamp \
        "$app_path" || {
        print_error "Failed to sign application"
        return 1
    }
    
    print_status "Application signed successfully"
}

# Verify signature
verify_signature() {
    local app_path="build/bin/moodle-prototype-manager.app"
    
    print_info "Verifying signature..."
    
    # Basic verification
    codesign --verify --verbose "$app_path" || {
        print_error "Signature verification failed"
        return 1
    }
    
    # Deep verification
    codesign --verify --deep --strict --verbose=2 "$app_path" || {
        print_warning "Deep verification found issues"
    }
    
    # Check signature details
    print_info "Signature details:"
    codesign --display --verbose=2 "$app_path"
    
    print_status "Signature verification complete"
}

# Create notarization package (optional)
create_notarization_zip() {
    print_info "Creating ZIP for notarization..."
    
    # Remove old zip if exists
    rm -f build/moodle-prototype-manager-signed.zip 2>/dev/null || true
    
    # Create zip
    cd build/bin
    zip -r ../moodle-prototype-manager-signed.zip moodle-prototype-manager.app
    cd - > /dev/null
    
    if [ -f "build/moodle-prototype-manager-signed.zip" ]; then
        local size=$(du -h "build/moodle-prototype-manager-signed.zip" | cut -f1)
        print_status "Created notarization package: build/moodle-prototype-manager-signed.zip ($size)"
    fi
}

# Main signing process
main() {
    # Parse arguments
    local identity=""
    local auto_sign=false
    local create_zip=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--identity)
                identity="$2"
                shift 2
                ;;
            -a|--auto)
                auto_sign=true
                shift
                ;;
            -z|--zip)
                create_zip=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Check prerequisites
    check_platform
    check_app_exists
    
    # List available identities
    echo
    if ! list_identities; then
        exit 1
    fi
    
    # Determine signing identity
    if [ -z "$identity" ]; then
        if [ "$auto_sign" = true ]; then
            # Try to find a valid identity automatically
            identity=$(security find-identity -v -p codesigning | grep -E "Developer ID Application|Apple Development" | head -1 | awk -F'"' '{print $2}')
            if [ -z "$identity" ]; then
                print_error "No valid signing identity found for automatic signing"
                exit 1
            fi
            print_info "Auto-selected identity: $identity"
        else
            # Prompt user for identity
            echo -e "${YELLOW}Enter the signing identity (copy from the list above):${NC}"
            read -r identity
            
            if [ -z "$identity" ]; then
                print_error "No identity provided"
                exit 1
            fi
        fi
    fi
    
    echo
    # Sign the application
    if sign_app "$identity"; then
        echo
        verify_signature
        
        if [ "$create_zip" = true ]; then
            echo
            create_notarization_zip
        fi
        
        echo
        print_status "Code signing completed successfully!"
        echo
        echo -e "${BLUE}📋 Next steps:${NC}"
        echo "  1. Test the signed application: open build/bin/moodle-prototype-manager.app"
        echo "  2. For distribution outside App Store, notarize the app:"
        echo "     xcrun notarytool submit build/moodle-prototype-manager-signed.zip --apple-id YOUR_APPLE_ID --wait"
        echo "  3. After notarization, staple the ticket:"
        echo "     xcrun stapler staple build/bin/moodle-prototype-manager.app"
        echo
    else
        exit 1
    fi
}

# Help function
show_help() {
    echo -e "${BLUE}Moodle Prototype Manager - macOS Code Signing Script${NC}"
    echo
    echo -e "${YELLOW}Usage:${NC}"
    echo "  ./sign-mac.sh [options]"
    echo
    echo -e "${YELLOW}Options:${NC}"
    echo "  -i, --identity <identity>  Specify signing identity (e.g., 'Developer ID Application: Your Name (TEAMID)')"
    echo "  -a, --auto                 Automatically select the first valid identity"
    echo "  -z, --zip                  Create ZIP file for notarization after signing"
    echo "  -h, --help                 Show this help message"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  ./sign-mac.sh                                      # Interactive mode"
    echo "  ./sign-mac.sh --auto                               # Auto-select identity"
    echo "  ./sign-mac.sh -i 'Developer ID Application: ...'   # Specify identity"
    echo "  ./sign-mac.sh --auto --zip                         # Auto-sign and create ZIP"
    echo
    echo -e "${YELLOW}Requirements:${NC}"
    echo "  - macOS with Xcode or Command Line Tools installed"
    echo "  - Valid Apple Developer signing certificate in Keychain"
    echo "  - Application built using ./build.sh"
    echo
    echo -e "${YELLOW}Certificate Types:${NC}"
    echo "  - Developer ID Application: For distribution outside App Store"
    echo "  - Apple Development: For development and testing"
    echo "  - Apple Distribution: For App Store submission"
}

# Run main function
main "$@"