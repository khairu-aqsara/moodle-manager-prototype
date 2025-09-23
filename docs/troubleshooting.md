# Moodle Prototype Manager - Troubleshooting Guide

## Table of Contents

1. [Quick Diagnostic Steps](#quick-diagnostic-steps)
2. [Application Won't Start](#application-wont-start)
3. [Health Check Issues](#health-check-issues)
4. [Docker-Related Problems](#docker-related-problems)
5. [Container Issues](#container-issues)
6. [Network and Connectivity Problems](#network-and-connectivity-problems)
7. [Image Download Issues](#image-download-issues)
8. [Startup and Timing Problems](#startup-and-timing-problems)
9. [Credential and Access Issues](#credential-and-access-issues)
10. [Performance Issues](#performance-issues)
11. [Platform-Specific Problems](#platform-specific-problems)
12. [File and Permission Issues](#file-and-permission-issues)
13. [Advanced Troubleshooting](#advanced-troubleshooting)
14. [Getting Additional Help](#getting-additional-help)

## Quick Diagnostic Steps

Before diving into specific issues, try these general diagnostic steps:

### 1. Check System Status
```
Open Moodle Prototype Manager
├── Look at footer status indicators
│   ├── Docker● (should be green)
│   ├── Internet● (should be green)
│   └── Status text (should show "Ready")
├── If any indicator is red, see relevant section below
└── If all green but still having issues, continue troubleshooting
```

### 2. Restart Sequence
Often resolves transient issues:
1. Close Moodle Prototype Manager completely
2. Stop any running Moodle containers: Click "Stop Moodle" if available
3. Restart Docker Desktop
4. Wait 30 seconds for Docker to fully start
5. Relaunch Moodle Prototype Manager

### 3. Verify Prerequisites
- **Docker Desktop**: Ensure it's running (check system tray/menu bar)
- **Internet Connection**: Test by browsing to any website
- **Port Availability**: Ensure port 8080 is not used by other applications
- **Disk Space**: Verify at least 2GB free space available

### 4. Check Application Version
- Look at header in application for version number
- Compare with latest release on GitHub
- Update if running an old version

## Application Won't Start

### Symptoms
- Application doesn't launch when clicked
- Error messages on startup
- Application crashes immediately
- White screen or empty window

### Windows-Specific Issues

**"Application failed to start" or DLL errors:**
```powershell
# Check for missing Visual C++ Redistributables
# Download from Microsoft and install:
# https://aka.ms/vs/17/release/vc_redist.x64.exe

# Check Windows version
winver

# Verify 64-bit system
systeminfo | findstr "System Type"
```

**"Windows protected your PC" message:**
1. Click "More info"
2. Click "Run anyway"
3. For permanent solution, see code signing section

**Antivirus blocking:**
1. Check antivirus quarantine/logs
2. Add application to antivirus whitelist
3. Temporarily disable real-time protection to test

### macOS-Specific Issues

**"App can't be opened because it is from an unidentified developer":**
1. Right-click the application
2. Select "Open"
3. Click "Open" in the security dialog

**"App is damaged and can't be opened":**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine "/Applications/Moodle Prototype Manager.app"

# If still doesn't work, re-download from official source
```

**Gatekeeper issues:**
```bash
# Check Gatekeeper status
spctl --status

# Assess application
spctl -a -v "/Applications/Moodle Prototype Manager.app"
```

### Linux-Specific Issues

**Permission denied:**
```bash
# Make binary executable
chmod +x moodle-prototype-manager
# or for AppImage
chmod +x Moodle\ Prototype\ Manager.AppImage
```

**Missing dependencies:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-37

# Fedora/RHEL
sudo dnf install gtk3 webkit2gtk3

# Check missing libraries
ldd moodle-prototype-manager
```

**AppImage issues:**
```bash
# Try extracting and running directly
./Moodle\ Prototype\ Manager.AppImage --appimage-extract
./squashfs-root/AppRun
```

### General Solutions

**Clean reinstall:**
1. Uninstall current version completely
2. Remove user data directory:
   - Windows: `%APPDATA%\moodle-prototype-manager`
   - macOS: `~/.moodle-prototype-manager`
   - Linux: `~/.moodle-prototype-manager`
3. Restart system
4. Download and install latest version

**Run from terminal/command prompt:**
```bash
# Windows
"C:\Program Files\Moodle Prototype Manager\Moodle Prototype Manager.exe"

# macOS
"/Applications/Moodle Prototype Manager.app/Contents/MacOS/Moodle Prototype Manager"

# Linux
./moodle-prototype-manager
```
This may show error messages not visible in GUI mode.

## Health Check Issues

### Docker Status Red

**Symptoms:** Red circle next to "Docker" in status bar

**Common Causes and Solutions:**

**Docker Desktop not running:**
```bash
# Check if Docker Desktop is running
docker version

# If command fails:
# - Windows: Start Docker Desktop from Start Menu
# - macOS: Start Docker Desktop from Applications
# - Linux: Start Docker service
systemctl start docker  # Linux with systemd
```

**Docker daemon not accessible:**
```bash
# Check Docker daemon status
docker info

# If permission denied on Linux:
sudo usermod -aG docker $USER
# Then log out and back in

# If Docker Desktop is starting:
# Wait 1-2 minutes for complete startup
```

**Docker Desktop installation issues:**
1. Uninstall Docker Desktop completely
2. Restart computer
3. Download latest version from docker.com
4. Install with administrator privileges
5. Restart computer again

**WSL2 issues (Windows):**
```powershell
# Check WSL2 installation
wsl --list --verbose

# Update WSL2 if needed
wsl --update

# Restart WSL2
wsl --shutdown
```

### Internet Status Red

**Symptoms:** Red circle next to "Internet" in status bar

**Network connectivity test:**
```bash
# Test basic connectivity
ping google.com

# Test HTTP connectivity
curl -I http://www.google.com

# Test HTTPS connectivity
curl -I https://www.google.com
```

**Common Causes:**
- No internet connection
- Firewall blocking application
- Proxy configuration issues
- DNS resolution problems

**Solutions:**

**Firewall configuration:**
- Add Moodle Prototype Manager to firewall exceptions
- Temporarily disable firewall to test
- Check corporate firewall settings

**Proxy configuration:**
```bash
# Check proxy settings
echo $http_proxy
echo $https_proxy

# For Docker Desktop with proxy:
# Configure proxy in Docker Desktop settings
```

**DNS issues:**
```bash
# Try different DNS servers
# Windows: Network settings → Change adapter options
# macOS: System Preferences → Network → Advanced → DNS
# Linux: Edit /etc/resolv.conf

# Test with public DNS
nslookup google.com 8.8.8.8
```

## Docker-Related Problems

### Docker Desktop Won't Start

**Windows:**
```powershell
# Check Hyper-V is enabled
Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V

# Check virtualization in BIOS
# Restart → Enter BIOS → Enable Intel VT-x/AMD-V

# Reset Docker Desktop
# Right-click Docker Desktop → Troubleshoot → Reset to factory defaults
```

**macOS:**
```bash
# Check virtualization framework
sysctl kern.hv_support

# Reset Docker Desktop
# Docker Desktop → Troubleshoot → Reset to factory defaults
```

**Linux:**
```bash
# Check Docker service
systemctl status docker

# Start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# Check user permissions
groups $USER | grep docker
```

### Docker Commands Fail

**Permission denied errors:**
```bash
# Linux: Add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Test after logout/login
docker run hello-world
```

**Docker daemon not responding:**
```bash
# Restart Docker Desktop or daemon
# Windows/macOS: Restart Docker Desktop
# Linux:
sudo systemctl restart docker
```

### Docker Images Issues

**Image pull failures:**
```bash
# Test Docker Hub connectivity
docker pull hello-world

# Check available disk space
df -h

# Clean up Docker space
docker system prune -a
```

**Corrupted images:**
```bash
# Remove specific image
docker rmi wenkhairu/moodle-prototype:502-stable

# Clean all images and start fresh
docker system prune -a --volumes
```

## Container Issues

### Container Won't Start

**Port already in use:**
```bash
# Check what's using port 8080
# Windows:
netstat -ano | findstr :8080
# macOS/Linux:
lsof -i :8080

# Stop conflicting service or use different port
```

**Container exists but won't start:**
```bash
# Check container status
docker ps -a | grep moodle

# Remove problematic container
docker rm [container-id]

# Delete container.id file to force new container creation
```

**Insufficient resources:**
```bash
# Check Docker resource allocation
docker system df

# Increase Docker Desktop memory allocation:
# Docker Desktop → Settings → Resources → Memory
```

### Container Starts But Moodle Inaccessible

**Container running but port not accessible:**
```bash
# Verify container port mapping
docker ps

# Test port connectivity
telnet localhost 8080
# or
curl http://localhost:8080
```

**Moodle still initializing:**
- First startup can take 5-20 minutes
- Check container logs: `docker logs [container-id]`
- Look for "Moodle is available at:" message

### Container Startup Timeout

**Symptoms:** Startup modal appears but never completes

**Solutions:**
1. **Check container logs:**
   ```bash
   # Find container ID
   docker ps -a | grep moodle

   # Check logs
   docker logs [container-id]
   ```

2. **Look for error messages:**
   - Database initialization errors
   - Permission issues
   - Resource constraints

3. **Try manual container start:**
   ```bash
   # Stop application
   # Start container manually
   docker start [container-id]

   # Monitor logs
   docker logs -f [container-id]
   ```

4. **Create fresh container:**
   - Stop application
   - Delete `container.id` file
   - Restart application

## Network and Connectivity Problems

### Can't Access Moodle at localhost:8080

**Verify Moodle is responding:**
```bash
# Test HTTP response
curl -I http://localhost:8080

# Test with different user agent
curl -H "User-Agent: Mozilla/5.0" http://localhost:8080
```

**Browser-specific issues:**
1. Try different browser
2. Clear browser cache
3. Disable browser extensions
4. Try incognito/private mode

**Localhost resolution issues:**
```bash
# Check hosts file
# Windows: C:\Windows\System32\drivers\etc\hosts
# macOS/Linux: /etc/hosts

# Should contain:
127.0.0.1 localhost
```

### Firewall Blocking Access

**Windows Firewall:**
1. Windows Security → Firewall & network protection
2. Allow an app through firewall
3. Add Moodle Prototype Manager
4. Check both Private and Public networks

**macOS Firewall:**
1. System Preferences → Security & Privacy → Firewall
2. Firewall Options
3. Add Moodle Prototype Manager

**Third-party firewalls:**
- Add application to whitelist
- Allow ports 8080 inbound and outbound
- Check for application-specific blocking rules

### Corporate Network Issues

**Proxy configuration:**
```bash
# Configure Docker Desktop proxy
# Docker Desktop → Settings → Resources → Proxies

# Set environment variables
export http_proxy=http://proxy.company.com:8080
export https_proxy=https://proxy.company.com:8080
```

**VPN conflicts:**
- Disconnect VPN temporarily to test
- Configure VPN to allow local traffic
- Use VPN split tunneling if available

## Image Download Issues

### Download Starts But Never Completes

**Check download progress:**
1. Download modal should show percentage
2. If stuck at same percentage for >10 minutes, cancel and retry
3. Check internet speed and stability

**Network timeout issues:**
```bash
# Test Docker Hub connectivity
docker pull hello-world

# Check Docker daemon timeout settings
# Increase timeout in Docker Desktop settings
```

### Download Fails With Error

**Network connectivity:**
```bash
# Test Docker registry connectivity
curl -I https://registry-1.docker.io

# Test DNS resolution
nslookup registry-1.docker.io
```

**Authentication issues:**
```bash
# If using private registry, login first
docker login

# Clear Docker credentials
docker logout
```

**Disk space issues:**
```bash
# Check available space
df -h

# Clean Docker system
docker system prune -a
```

### Very Slow Download

**Optimize download speed:**
1. Use wired internet connection
2. Close bandwidth-heavy applications
3. Try downloading during off-peak hours

**Change Docker registry mirror:**
```bash
# Configure registry mirror in Docker Desktop
# Docker Desktop → Settings → Docker Engine
# Add mirror configuration
```

**Restart download:**
1. Stop application
2. Remove partially downloaded image: `docker rmi wenkhairu/moodle-prototype:502-stable`
3. Restart application and try again

## Startup and Timing Problems

### Moodle Takes Too Long to Start

**Expected startup times:**
- First time: 5-20 minutes (Windows may take longer)
- Subsequent starts: 1-5 minutes
- Varies by system performance and Docker configuration

**Optimize startup performance:**
1. **Increase Docker resources:**
   - Docker Desktop → Settings → Resources
   - Memory: 4-8GB
   - CPUs: 2-4 cores

2. **Use SSD storage:**
   - Store Docker images on SSD if possible
   - Move Docker Desktop storage location to SSD

3. **Close other applications:**
   - Free up system resources
   - Stop unnecessary background processes

### Startup Hangs at Specific Point

**Check container logs:**
```bash
# Find container
docker ps -a | grep moodle

# Check logs
docker logs [container-id]

# Monitor in real-time
docker logs -f [container-id]
```

**Common hang points:**
1. **Database initialization:** Look for MySQL/MariaDB messages
2. **Web server startup:** Look for Apache/Nginx messages
3. **Moodle installation:** Look for Moodle setup messages

**Recovery options:**
1. Wait longer (up to 30 minutes for first startup)
2. Restart container manually
3. Create new container if persistent

### Application Freezes During Startup

**Application not responding:**
1. Don't force-quit immediately
2. Wait 5-10 minutes for timeout
3. If still frozen, force-quit and restart

**Check system resources:**
```bash
# Monitor resource usage
# Windows: Task Manager
# macOS: Activity Monitor
# Linux: htop or system monitor
```

**Recovery steps:**
1. Force-quit application
2. Stop Docker containers: `docker stop $(docker ps -q)`
3. Restart Docker Desktop
4. Restart application

## Credential and Access Issues

### Can't See Admin Password

**Password not displayed:**
1. Ensure container is fully started (green status)
2. Check that credentials section is visible
3. Try stopping and restarting container

**Password field empty:**
1. Stop application
2. Delete `moodle.txt` file from user data directory
3. Restart application and container

### Can't Login to Moodle

**Verify credentials:**
1. Copy password using copy button (📋)
2. Ensure username is exactly "admin" (lowercase)
3. Check for extra spaces in password

**Moodle access issues:**
```bash
# Verify Moodle is responding
curl -I http://localhost:8080

# Check container logs for errors
docker logs [container-id] | grep -i error
```

**Clear browser data:**
1. Clear cookies for localhost
2. Clear cached data
3. Try different browser

### Copy Button Not Working

**Clipboard API issues:**
- Try manually selecting and copying password
- Ensure browser supports clipboard API
- Grant clipboard permissions if prompted

**Browser security:**
- Some browsers block clipboard access on localhost
- Try accessing via 127.0.0.1:8080 instead
- Use HTTPS if configured

## Performance Issues

### Application Runs Slowly

**System resource constraints:**
```bash
# Check memory usage
# Ensure at least 4GB RAM available for Docker

# Check CPU usage
# Close unnecessary applications

# Check disk space
# Ensure at least 2GB free space
```

**Docker Desktop optimization:**
```
Docker Desktop → Settings → Resources:
- Memory: 4-8GB
- CPUs: 2-4 cores
- Swap: 1-2GB
```

### Moodle Responds Slowly

**Container resource limits:**
- Increase Docker memory allocation
- Check container resource usage: `docker stats`

**Network latency:**
```bash
# Test localhost latency
ping localhost

# Check for network issues
netstat -an | grep 8080
```

**Database performance:**
- First-time database setup is always slow
- Performance improves after initial setup
- Consider persistent database optimizations

### High CPU or Memory Usage

**Docker Desktop consuming resources:**
1. Adjust resource limits in Docker Desktop settings
2. Close unused containers: `docker stop $(docker ps -q)`
3. Clean up Docker system: `docker system prune`

**Application memory leaks:**
- Restart application periodically
- Monitor memory usage in task manager
- Report persistent memory leaks as bugs

## Platform-Specific Problems

### Windows-Specific Issues

**WSL2 integration problems:**
```powershell
# Check WSL2 status
wsl --list --verbose

# Restart WSL2
wsl --shutdown

# Update WSL2
wsl --update
```

**Hyper-V conflicts:**
```powershell
# Disable conflicting features
dism.exe /online /disable-feature /featurename:VirtualMachinePlatform

# Re-enable for Docker
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform
```

**Path length limitations:**
- Use shorter installation paths
- Enable long path support in Windows
- Move application closer to root directory

### macOS-Specific Issues

**Apple Silicon (M1/M2) compatibility:**
```bash
# Check architecture
uname -m

# Use Apple Silicon native Docker Desktop
# Ensure using ARM64 version of application
```

**File system permissions:**
```bash
# Fix permissions on user data directory
chmod -R 755 ~/.moodle-prototype-manager

# Reset file attributes
xattr -r -d com.apple.quarantine ~/.moodle-prototype-manager
```

**System Extension blocking:**
1. System Preferences → Security & Privacy
2. Allow Docker Desktop system extension
3. Restart Docker Desktop

### Linux-Specific Issues

**Docker daemon not starting:**
```bash
# Check systemd status
systemctl status docker

# Check for conflicts
journalctl -u docker.service

# Reinstall Docker if needed
sudo apt remove docker-ce docker-ce-cli containerd.io
sudo apt install docker-ce docker-ce-cli containerd.io
```

**X11/Wayland display issues:**
```bash
# Check display server
echo $XDG_SESSION_TYPE

# For Wayland compatibility issues:
GDK_BACKEND=x11 ./moodle-prototype-manager
```

**AppImage FUSE issues:**
```bash
# Install FUSE if needed
sudo apt install fuse

# Or extract and run without FUSE
./Moodle\ Prototype\ Manager.AppImage --appimage-extract-and-run
```

## File and Permission Issues

### Can't Create Configuration Files

**Permission denied errors:**
```bash
# Check user data directory permissions
ls -la ~/.moodle-prototype-manager

# Fix permissions
chmod 755 ~/.moodle-prototype-manager
chmod 644 ~/.moodle-prototype-manager/*
```

**Directory doesn't exist:**
```bash
# Create user data directory manually
mkdir -p ~/.moodle-prototype-manager

# Set proper permissions
chmod 755 ~/.moodle-prototype-manager
```

### Configuration Files Corrupted

**Reset configuration:**
```bash
# Backup existing files
cp ~/.moodle-prototype-manager/moodle.txt ~/moodle.txt.backup

# Remove corrupted files
rm ~/.moodle-prototype-manager/container.id
rm ~/.moodle-prototype-manager/moodle.txt

# Restart application
```

### File System Errors

**Disk full errors:**
```bash
# Check disk space
df -h

# Clean up Docker
docker system prune -a --volumes

# Clean system temporary files
```

**File system corruption:**
```bash
# Check file system
# Windows: chkdsk C: /f
# macOS: diskutil verifyVolume /
# Linux: fsck /dev/sdX
```

## Advanced Troubleshooting

### Collecting Diagnostic Information

**Application logs:**
Look for log files in:
- Windows: `%APPDATA%\moodle-prototype-manager\logs\`
- macOS: `~/.moodle-prototype-manager/logs/`
- Linux: `~/.moodle-prototype-manager/logs/`

**Docker logs:**
```bash
# Application container logs
docker logs [container-id] > container.log

# Docker system info
docker system info > docker-info.txt

# Docker daemon logs
# Windows: %LOCALAPPDATA%\Docker\log\
# macOS: ~/Library/Containers/com.docker.docker/Data/log/
# Linux: journalctl -u docker.service
```

**System information:**
```bash
# Windows
systeminfo > system-info.txt
dxdiag /t dxdiag-report.txt

# macOS
system_profiler SPHardwareDataType > hardware-info.txt
sw_vers > version-info.txt

# Linux
uname -a > system-info.txt
lscpu > cpu-info.txt
free -h > memory-info.txt
```

### Network Debugging

**Port connectivity testing:**
```bash
# Test port 8080
telnet localhost 8080

# Network interface information
# Windows: ipconfig /all
# macOS/Linux: ifconfig -a

# Routing table
# Windows: route print
# macOS/Linux: netstat -rn
```

**Docker network inspection:**
```bash
# List Docker networks
docker network ls

# Inspect default bridge network
docker network inspect bridge

# Check container network settings
docker inspect [container-id] | grep -A 10 NetworkSettings
```

### Docker Engine Debugging

**Enable Docker debug logging:**
```json
// Docker daemon.json configuration
{
  "debug": true,
  "log-level": "debug"
}
```

**Docker system diagnostics:**
```bash
# Docker system events
docker system events

# Docker system disk usage
docker system df -v

# Detailed container inspection
docker inspect [container-id]
```

### Manual Container Management

**If application container management fails:**
```bash
# List all containers
docker ps -a

# Start container manually
docker start [container-id]

# Execute commands in container
docker exec -it [container-id] /bin/bash

# Check container logs
docker logs -f [container-id]

# Remove problematic container
docker rm -f [container-id]
```

## Getting Additional Help

### Before Seeking Help

**Gather this information:**
1. **System details:**
   - Operating system and version
   - Application version
   - Docker Desktop version

2. **Error information:**
   - Exact error messages
   - Screenshots of issues
   - Steps to reproduce problem

3. **Log files:**
   - Application logs
   - Docker container logs
   - System logs (if relevant)

4. **Configuration:**
   - Contents of `image.docker` file
   - Docker Desktop settings
   - Network configuration

### Self-Help Resources

**Official documentation:**
- README.md in project repository
- GitHub repository wiki
- API documentation
- Build and deployment guides

**Community resources:**
- GitHub Issues search
- GitHub Discussions
- Docker community forums
- Moodle developer forums

### Reporting Bugs

**Create detailed bug reports:**
```markdown
## Bug Report

**Environment:**
- OS: Windows 11 Pro 22H2
- Application Version: 1.0.0
- Docker Desktop Version: 4.15.0

**Description:**
Brief description of the issue

**Steps to Reproduce:**
1. Step one
2. Step two
3. Step three

**Expected Behavior:**
What should happen

**Actual Behavior:**
What actually happens

**Logs:**
```
Paste relevant log entries here
```

**Additional Context:**
Any other relevant information
```

### Getting Community Support

**GitHub Discussions:**
- Ask questions about usage
- Share tips and tricks
- Discuss feature requests
- Help other users

**GitHub Issues:**
- Report confirmed bugs
- Request new features
- Suggest improvements
- Track development progress

### Emergency Recovery

**If nothing else works:**
1. **Complete reset:**
   ```bash
   # Stop all containers
   docker stop $(docker ps -q)

   # Remove all containers
   docker rm $(docker ps -aq)

   # Remove all images
   docker rmi $(docker images -q)

   # Remove application data
   rm -rf ~/.moodle-prototype-manager

   # Restart Docker Desktop
   # Reinstall application
   ```

2. **System restart:**
   - Save all work
   - Restart computer
   - Start Docker Desktop
   - Launch application

3. **Contact support:**
   - Provide detailed information
   - Include log files
   - Be available for follow-up questions

This troubleshooting guide covers the most common issues and their solutions. For issues not covered here, please check the project documentation or seek community support through the official channels.