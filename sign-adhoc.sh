#!/bin/bash

# Ad-hoc signing script for macOS app (no certificate required)

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🔏 Ad-hoc Code Signing (No Certificate)${NC}"
echo -e "${BLUE}======================================${NC}"

if [ ! -d "build/bin/moodle-prototype-manager.app" ]; then
    echo "❌ App not found. Run ./build.sh first"
    exit 1
fi

# Ad-hoc sign (uses "-" as identity)
echo "Signing with ad-hoc signature..."
codesign --force --deep -s - "build/bin/moodle-prototype-manager.app"

# Verify
echo -e "\n${GREEN}✅ Ad-hoc signing complete${NC}"
codesign -v "build/bin/moodle-prototype-manager.app"

# Create distribution ZIP
echo -e "\nCreating distribution package..."
cd build/bin
zip -r ../moodle-prototype-manager-adhoc-signed.zip moodle-prototype-manager.app
cd ../..

echo -e "\n${GREEN}✅ Done!${NC}"
echo -e "${BLUE}Distribution file:${NC} build/moodle-prototype-manager-adhoc-signed.zip"
echo -e "\n${YELLOW}Note:${NC} Users will still need to right-click → Open on first launch"
echo "But it's slightly easier than completely unsigned apps."