# Moodle Prototype Manager - Installation Guide

## For macOS Users

Since this app is signed with a development certificate (not notarized), macOS will show a security warning on first launch. This is normal and the app is safe to use.

### Installation Steps:

1. **Download and Extract**
   - Download the `moodle-prototype-manager-signed.zip` file
   - Double-click to extract the app
   - Move `moodle-prototype-manager.app` to your Applications folder (optional)

2. **First Launch - Override Security Warning**
   
   **Method 1: System Settings (Recommended)**
   - Double-click the app - you'll see "cannot be opened because it is from an unidentified developer"
   - Click "OK" to dismiss the dialog
   - Open **System Settings** → **Privacy & Security**
   - Scroll down to see "moodle-prototype-manager.app was blocked"
   - Click **"Open Anyway"**
   - Enter your password if prompted
   - Click **"Open"** in the final dialog

   **Method 2: Right-Click Method**
   - Right-click (or Control-click) on the app
   - Select **"Open"** from the menu
   - Click **"Open"** in the dialog that appears
   - Enter your password if prompted

3. **Subsequent Launches**
   - After the first successful launch, the app will open normally without warnings

### Requirements:
- macOS 11.0 or later
- Docker Desktop installed and running
- Internet connection for downloading Moodle Docker image

### Troubleshooting:

**If you see "App is damaged and can't be opened":**
```bash
# Open Terminal and run:
xattr -cr /Applications/moodle-prototype-manager.app
```

**If the app won't launch:**
1. Make sure Docker Desktop is installed and running
2. Check that port 8080 is not in use by another application

### Security Note:
This app is signed by the developer but not notarized by Apple. The security warning is a standard macOS protection. The app is safe and only manages a local Docker container for Moodle development.