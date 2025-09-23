# Moodle Prototype Manager - User Guide

## Table of Contents

1. [Getting Started](#getting-started)
2. [System Requirements](#system-requirements)
3. [Installation](#installation)
4. [First Launch](#first-launch)
5. [Using the Application](#using-the-application)
6. [Understanding the Interface](#understanding-the-interface)
7. [Working with Containers](#working-with-containers)
8. [Accessing Moodle](#accessing-moodle)
9. [Configuration Options](#configuration-options)
10. [Data and File Management](#data-and-file-management)
11. [Common Workflows](#common-workflows)
12. [Tips and Best Practices](#tips-and-best-practices)

## Getting Started

### What is Moodle Prototype Manager?

Moodle Prototype Manager is a desktop application that simplifies running and managing Moodle prototype environments using Docker containers. It provides an easy-to-use graphical interface for developers, testers, and educators who need quick access to Moodle instances without complex command-line operations.

### Key Features

- **One-Click Moodle Launch**: Start a complete Moodle environment with a single click
- **Automatic Setup**: Downloads and configures everything needed automatically
- **Credential Management**: Automatically extracts and displays login credentials
- **Health Monitoring**: Real-time status of Docker and internet connectivity
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **No Docker Knowledge Required**: Hides Docker complexity behind a simple interface

### What You'll Need

- A computer running Windows 10/11, macOS 10.15+, or Linux
- Docker Desktop installed and running
- Internet connection for downloading Moodle images
- At least 4GB of available RAM
- 2GB of free disk space

## System Requirements

### Minimum Requirements

**Operating Systems:**
- **Windows**: Windows 10 (64-bit) or Windows 11
- **macOS**: macOS 10.15 (Catalina) or later
- **Linux**: Ubuntu 18.04+, Fedora 30+, or equivalent distributions

**Hardware:**
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 2GB free space (more needed for multiple Moodle instances)
- **Network**: Internet connection for initial setup and image downloads
- **CPU**: Any modern processor (x64/AMD64)

### Docker Requirements

**Docker Desktop:**
- **Windows**: Docker Desktop 4.0 or later
- **macOS**: Docker Desktop 4.0 or later
- **Linux**: Docker Engine 20.10+ or Docker Desktop

**Docker Configuration:**
- At least 2GB RAM allocated to Docker
- Access to Docker daemon (user must be in docker group on Linux)
- Port 8080 available (used by Moodle container)

### Network Requirements

- Internet access for downloading Docker images (approximately 1GB download)
- Port 8080 must be available for Moodle web interface
- Firewall permissions for Docker and the application

## Installation

### Windows Installation

**Method 1: Installer (Recommended)**
1. Download `Moodle Prototype Manager Setup.exe` from the releases page
2. Right-click the installer and select "Run as administrator"
3. Follow the installation wizard
4. Launch from Start Menu or Desktop shortcut

**Method 2: Portable Version**
1. Download `Moodle Prototype Manager Windows x64.zip`
2. Extract to your preferred location (e.g., `C:\Programs\MoodleManager`)
3. Run `Moodle Prototype Manager.exe` from the extracted folder

### macOS Installation

**Method 1: DMG Package (Recommended)**
1. Download `Moodle Prototype Manager.dmg` from the releases page
2. Open the DMG file
3. Drag "Moodle Prototype Manager" to your Applications folder
4. Launch from Applications or Spotlight search

**Method 2: Direct Binary**
1. Download the macOS binary
2. Move to Applications folder
3. Right-click and select "Open" to bypass Gatekeeper (first time only)

**Security Note:** If you see "cannot be opened because it is from an unidentified developer":
1. Right-click the application and select "Open"
2. Click "Open" in the security dialog
3. The application will remember this permission for future launches

### Linux Installation

**Method 1: AppImage (Recommended)**
1. Download `Moodle Prototype Manager.AppImage`
2. Make it executable: `chmod +x Moodle\ Prototype\ Manager.AppImage`
3. Run: `./Moodle\ Prototype\ Manager.AppImage`

**Method 2: Binary**
1. Download the Linux binary
2. Make it executable: `chmod +x moodle-prototype-manager`
3. Run: `./moodle-prototype-manager`

### Docker Desktop Installation

If Docker Desktop is not already installed:

**Windows:**
1. Visit https://docs.docker.com/docker-for-windows/install/
2. Download Docker Desktop for Windows
3. Run the installer and follow the setup wizard
4. Restart your computer when prompted
5. Launch Docker Desktop and complete the initial setup

**macOS:**
1. Visit https://docs.docker.com/docker-for-mac/install/
2. Download Docker Desktop for Mac
3. Drag Docker.app to Applications folder
4. Launch Docker Desktop and complete the initial setup

**Linux:**
1. Follow the installation guide for your distribution at https://docs.docker.com/engine/install/
2. Add your user to the docker group: `sudo usermod -aG docker $USER`
3. Log out and back in for group changes to take effect

## First Launch

### Initial Setup Process

1. **Launch the Application**
   - Windows: Use Start Menu or Desktop shortcut
   - macOS: Open from Applications folder or Spotlight
   - Linux: Run the AppImage or binary

2. **First-Time Interface**
   When you first open the application, you'll see:
   - Moodle logo at the top
   - "Run Moodle" button (initially disabled)
   - Status indicators at the bottom showing system health

3. **Health Check Process**
   The application automatically checks:
   - **Docker Status**: Verifies Docker Desktop is running
   - **Internet Status**: Confirms internet connectivity

   Status indicators will show:
   - 🔴 Red: Problem detected
   - 🟢 Green: System ready
   - 🟡 Yellow: Checking in progress

4. **Resolving Health Check Issues**

   **If Docker Status is Red:**
   - Ensure Docker Desktop is installed and running
   - Check that your user has Docker permissions
   - Try restarting Docker Desktop

   **If Internet Status is Red:**
   - Check your internet connection
   - Verify firewall settings allow the application to access the internet
   - Try accessing a website in your browser to confirm connectivity

5. **Ready to Use**
   When both status indicators are green:
   - The "Run Moodle" button becomes enabled
   - Status text changes to "Ready"
   - You can now start your first Moodle container

## Using the Application

### Starting Moodle for the First Time

1. **Click "Run Moodle"**
   - The button is only available when health checks pass
   - This begins the automatic setup process

2. **Image Download Process**
   If this is your first time or the Moodle image isn't available locally:
   - A download progress modal appears
   - Shows download percentage and status
   - Download size is approximately 1GB
   - Download time varies based on internet speed (5-30 minutes typical)

3. **Container Startup**
   After download completes:
   - A startup modal appears with a loading spinner
   - Message: "Starting Moodle, please wait..."
   - The application monitors container startup progress
   - First-time startup can take 5-20 minutes (Windows may take longer)

4. **Credential Extraction**
   During startup, the application:
   - Monitors container logs for Moodle initialization
   - Automatically extracts the admin password
   - Detects when Moodle is ready to use

5. **Ready to Use**
   When startup completes:
   - The startup modal closes
   - Login credentials appear in the main window
   - Button changes to "Stop Moodle" (red color)
   - Browser launch dialog may appear

### Subsequent Launches

After the first successful launch:

1. **Click "Run Moodle"**
   - No download needed (image already exists)
   - Faster startup (typically 1-5 minutes)

2. **Container Restart**
   - Uses existing container from previous session
   - Preserves Moodle data and configuration
   - Shows same credentials as before

3. **Quick Access**
   - Credentials appear immediately
   - Can access Moodle while container is starting
   - Full functionality available once startup completes

### Stopping Moodle

1. **Click "Stop Moodle"**
   - Button turns red when container is running
   - Initiates graceful container shutdown

2. **Shutdown Process**
   - Application attempts graceful stop first
   - Falls back to forced stop if needed
   - Usually completes within 30 seconds

3. **Post-Shutdown State**
   - Button returns to "Run Moodle" (green)
   - Credentials hidden from view
   - Container is stopped but preserved for next launch

## Understanding the Interface

### Main Window Layout

```
┌─────────────────────────────────────┐
│              HEADER                 │
│            [Moodle Logo]            │
│        Moodle Prototype v1.0        │
├─────────────────────────────────────┤
│                                     │
│              MAIN AREA              │
│                                     │
│          [Run Moodle Button]        │
│                                     │
│       [Credentials Display]         │
│         (when container running)    │
│                                     │
├─────────────────────────────────────┤
│ FOOTER                              │
│ Status: Ready     Docker● Internet● │
└─────────────────────────────────────┘
```

### Interface Elements

**Header Section:**
- **Moodle Logo**: Visual brand identifier
- **Version Text**: Shows application version number

**Main Action Button:**
- **"Run Moodle" (Green)**: Start container (when stopped)
- **"Stop Moodle" (Red)**: Stop container (when running)
- **Disabled (Gray)**: Health checks failed or operation in progress

**Credentials Display (when container running):**
```
Username: admin
Password: [generated-password]  📋
URL:      http://localhost:8080
```

**Footer Status Bar:**
- **Left Side**: Current operation status ("Ready", "Checking...", "Starting...", etc.)
- **Right Side**: Health indicators
  - **Docker●**: Docker daemon connectivity
  - **Internet●**: Internet connection status

### Status Indicators

**Color Meanings:**
- **🟢 Green**: System healthy and ready
- **🔴 Red**: Problem detected, needs attention
- **🟡 Yellow**: Checking or operation in progress

**Status Messages:**
- **"Checking..."**: Initial health checks in progress
- **"Ready"**: All systems operational, ready to start Moodle
- **"Starting..."**: Container startup in progress
- **"Running"**: Moodle container is active
- **"Stopping..."**: Container shutdown in progress
- **Error messages**: Specific issues that need resolution

### Modal Dialogs

**Download Progress Modal:**
- Appears during first-time image download
- Shows progress bar and percentage
- Displays current download status
- Cannot be dismissed (download must complete)

**Startup Wait Modal:**
- Appears during container initialization
- Shows loading spinner
- Message: "Starting Moodle, please wait..."
- Cannot be dismissed (startup must complete)

**Browser Confirmation Dialog:**
- Appears after successful container startup
- Options: "Yes" or "No"
- "Yes" opens Moodle in your default browser
- "No" closes dialog without opening browser

## Working with Containers

### Container Lifecycle

**Container States:**
1. **No Container**: Initial state, no container exists
2. **Container Stopped**: Container exists but not running
3. **Container Starting**: Container is initializing
4. **Container Running**: Container active and serving Moodle
5. **Container Stopping**: Container shutting down

**State Transitions:**
```
No Container ──Run──> Starting ──Ready──> Running
                         │         │
                      [Error]   Stop
                         │         │
                         ▼         ▼
                     Stopped ◄──── Stopping
```

### Understanding Container Behavior

**First Launch (No Container Exists):**
1. Downloads Moodle Docker image (if needed)
2. Creates new container with random admin password
3. Starts container and waits for Moodle initialization
4. Extracts and saves credentials
5. Container is ready for use

**Subsequent Launches (Container Exists):**
1. Checks if existing container is available
2. Starts existing container (faster than creating new)
3. Uses previously saved credentials
4. Container is ready for use

**Container Persistence:**
- Containers are preserved between application sessions
- Moodle data and configuration persist across restarts
- Admin password remains consistent
- Custom changes to Moodle are maintained

### Container Management

**Automatic Management:**
- Application handles all container operations
- No manual Docker commands needed
- Automatic error recovery when possible
- Graceful shutdown on application exit

**Manual Container Recovery:**
If you need to start fresh:
1. Stop the current container using the application
2. Remove container files:
   - Delete `container.id` file (in user data directory)
   - Delete `moodle.txt` file (in user data directory)
3. Restart the application
4. Next launch will create a new container

**User Data Directory Locations:**
- **Windows**: `%APPDATA%\moodle-prototype-manager\`
- **macOS**: `~/.moodle-prototype-manager/`
- **Linux**: `~/.moodle-prototype-manager/`

## Accessing Moodle

### Login Process

1. **Get Credentials**
   When container is running, credentials are displayed in the application:
   ```
   Username: admin
   Password: [random-generated-password]
   URL:      http://localhost:8080
   ```

2. **Open Moodle**
   - Click the URL link in the application (opens browser automatically)
   - Or manually navigate to http://localhost:8080 in any browser
   - Or use the browser confirmation dialog when it appears

3. **Login to Moodle**
   - Enter username: `admin`
   - Enter the password shown in the application
   - Click "Log in"

### Copying Credentials

**Copy Password:**
- Click the 📋 copy button next to the password
- Password is copied to your clipboard
- Button briefly changes color to confirm copy

**Manual Copy:**
- Select the password text with your mouse
- Copy using Ctrl+C (Windows/Linux) or Cmd+C (macOS)

### Browser Access

**Direct URL Access:**
- Moodle is always available at http://localhost:8080 when running
- Bookmark this URL for quick access
- Works in any modern web browser

**Multiple Browser Windows:**
- You can open multiple browser tabs/windows
- All will access the same Moodle instance
- Changes are immediately visible across all windows

### Moodle Features

When you access Moodle, you'll have a fully functional Moodle environment:

**Standard Moodle Features:**
- Course creation and management
- User enrollment and management
- Content creation (lessons, quizzes, assignments)
- Grade book functionality
- Forums and messaging
- File uploads and management

**Admin Capabilities:**
- Full administrative access
- Plugin installation and configuration
- System settings modification
- User and role management
- Backup and restore functionality

**Development Features:**
- Clean Moodle installation
- Debugging options available
- Full database access (if needed)
- Log file access (if needed)

## Configuration Options

### Image Configuration

**Changing Docker Image:**
The application uses a configuration file to determine which Moodle image to use.

1. **Locate Configuration File:**
   - **Development**: `image.docker` in application directory
   - **Production**: `image.docker` in user data directory

2. **Edit Configuration:**
   - Open `image.docker` in a text editor
   - Replace content with desired image name
   - Example: `your-custom/moodle-image:dev`
   - Save the file

3. **Apply Changes:**
   - Restart the application
   - Next launch will use the new image
   - May require download if image not available locally

**Supported Images:**
- Any Docker image that runs Moodle on port 8080
- Must output admin credentials to container logs
- Should follow standard Moodle container patterns

### Application Settings

Currently, the application has minimal user-configurable settings. Configuration is primarily through:

**File-Based Configuration:**
- `image.docker`: Docker image name
- `container.id`: Current container identifier (auto-managed)
- `moodle.txt`: Stored credentials (auto-managed)

**Environment Variables (Advanced):**
- `MOODLE_MANAGER_DEBUG`: Enable debug logging
- `MOODLE_MANAGER_PORT`: Change default port (advanced users)

### Docker Configuration

**Docker Desktop Settings:**
Access through Docker Desktop preferences:

**Resources:**
- **Memory**: Allocate at least 2GB to Docker
- **CPU**: 2+ cores recommended for better performance
- **Disk**: Ensure sufficient space for images and containers

**Advanced:**
- **Port Forwarding**: Ensure port 8080 is available
- **File Sharing**: Enable if planning to mount external directories

## Data and File Management

### Understanding Data Storage

**Application Data:**
- **Configuration**: Stored in user data directories
- **Container State**: Managed automatically by application
- **Credentials**: Encrypted and stored locally

**Moodle Data:**
- **Database**: Stored inside container (persists across restarts)
- **Uploaded Files**: Stored inside container
- **Configuration**: Stored inside container
- **Logs**: Accessible through Docker if needed

### Data Persistence

**What Persists:**
- Moodle database and user data
- Uploaded files and content
- Course configurations
- User accounts and enrollments
- Custom themes and plugins

**What Doesn't Persist:**
- Container logs (cleared on container removal)
- Temporary files
- Debug information

### Backup and Recovery

**Automatic Backup (Built-in to Container):**
- Container data persists across application restarts
- Moodle data survives container stop/start cycles

**Manual Backup (Advanced):**
If you need to backup Moodle data:
1. Use Moodle's built-in backup functionality
2. Export courses and user data through Moodle interface
3. Download files through Moodle file manager

**Recovery Options:**
- Restart existing container to recover from temporary issues
- Create new container for fresh start
- Use Moodle restore functionality for data recovery

### Disk Space Management

**Monitoring Usage:**
- Docker images: ~1GB per Moodle version
- Container data: Varies based on usage (100MB-10GB typical)
- Application files: <50MB

**Cleaning Up:**
1. **Remove Old Containers:**
   - Stop current container in application
   - Delete `container.id` and `moodle.txt` files
   - Restart application for fresh container

2. **Docker Cleanup:**
   - Use Docker Desktop to manage images and containers
   - Remove unused images to free space
   - Clear Docker build cache if needed

## Common Workflows

### Daily Development Workflow

**Starting Your Day:**
1. Launch Moodle Prototype Manager
2. Wait for health checks to complete
3. Click "Run Moodle"
4. Access Moodle in browser when ready
5. Begin development or testing work

**Ending Your Day:**
1. Save any work in Moodle
2. Click "Stop Moodle" in the application
3. Close the application
4. Work is preserved for next session

### Testing Workflow

**Setting Up Test Environment:**
1. Start fresh container (delete `container.id` if needed)
2. Launch Moodle and note new admin credentials
3. Configure test data and scenarios
4. Perform testing activities
5. Document results

**Resetting for Next Test:**
1. Stop container
2. Delete container files for fresh start
3. Restart application
4. New container with clean Moodle installation

### Development Workflow

**Plugin Development:**
1. Start Moodle container
2. Access Moodle admin interface
3. Install development tools and debugging
4. Upload and test plugins
5. Use browser developer tools for debugging

**Theme Development:**
1. Start Moodle container
2. Enable theme debugging in Moodle
3. Upload and activate custom themes
4. Test across different browsers
5. Iterate on design changes

### Demonstration Workflow

**Preparing for Demo:**
1. Start with fresh container for consistent experience
2. Pre-configure demo content and users
3. Test all demo scenarios beforehand
4. Have backup plan with second container if needed

**During Demo:**
1. Start container well before presentation
2. Have browser bookmarked to Moodle URL
3. Keep application visible for troubleshooting
4. Monitor status indicators during demo

## Tips and Best Practices

### Performance Tips

**Optimize Docker Performance:**
- Allocate sufficient RAM to Docker (4-8GB recommended)
- Use SSD storage for better Docker performance
- Close other resource-intensive applications during use
- Regularly clean up unused Docker images

**Application Performance:**
- Keep Docker Desktop updated
- Restart Docker Desktop periodically
- Monitor disk space regularly
- Use wired internet connection for faster downloads

### Security Considerations

**Local Development Security:**
- Moodle runs on localhost only (not accessible externally)
- Use randomly generated admin passwords
- Don't use production data in prototype environment
- Regularly update Docker Desktop for security patches

**Credential Management:**
- Don't share screenshots containing passwords
- Use copy button instead of manual selection when possible
- Consider changing admin password through Moodle if needed

### Troubleshooting Prevention

**Preventive Measures:**
- Always check health indicators before starting
- Keep Docker Desktop running when using application
- Ensure stable internet connection for downloads
- Don't force-quit application during operations

**Best Practices:**
- Let downloads complete fully
- Wait for startup process to finish
- Use "Stop Moodle" button before closing application
- Restart application if status indicators show problems

### Advanced Usage

**Multiple Environments:**
- Use different image configurations for different projects
- Maintain separate application installations if needed
- Document image configurations for team sharing

**Integration with Development Tools:**
- Bookmark Moodle URL for quick access
- Use browser developer tools for debugging
- Integrate with version control for plugin development
- Document container configurations for team use

### Getting Help

**Self-Help Resources:**
- Check application status indicators first
- Review logs in user data directory
- Test Docker Desktop independently
- Restart application for transient issues

**Community Support:**
- Check project documentation and README
- Search existing GitHub issues
- Ask questions in project discussions
- Report bugs with detailed information

**Information for Support Requests:**
- Operating system version
- Docker Desktop version
- Application version
- Error messages or screenshots
- Steps to reproduce issues

This user guide provides comprehensive coverage of all aspects of using the Moodle Prototype Manager application, from installation through advanced usage scenarios.