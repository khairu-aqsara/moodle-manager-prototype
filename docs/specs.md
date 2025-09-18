# Moodle Prototype Manager - Technical Specifications

## 1. Project Overview

### 1.1 Application Purpose
A cross-platform desktop application built with Golang and Wails that provides a graphical interface for managing the `wenkhairu/moodle-prototype:502-stable` Docker container. The application handles the complete lifecycle of the Moodle prototype including image management, container orchestration, credential extraction, and user interaction.

### 1.2 Target Platforms
- Windows (amd64)
- macOS (Intel and Apple Silicon)

### 1.3 Technology Stack
- **Backend**: Go (Golang) with Wails v2 framework
- **Frontend**: HTML5, CSS3, JavaScript (embedded in Wails)
- **Container Runtime**: Docker Engine
- **Target Image**: `wenkhairu/moodle-prototype:502-stable`

## 2. Functional Requirements

### 2.1 First-Time Application Launch

#### 2.1.1 Initial UI State
- **Header Section**:
  - Moodle logo displayed at the top center (128x128 pixels)
  - Version information displayed below the logo ("Moodle Prototype v[version]")

- **Main Action Area**:
  - "Run Moodle" button (initially disabled)

- **Footer Status Bar**:
  - Left side: "Checking..." status text
  - Right side: Status indicators with colored circles
    - Docker status: Red circle with "Docker" label
    - Internet status: Red circle with "Internet" label

#### 2.1.2 System Health Checks
The application performs background health checks on startup:

**Docker Check Process**:
1. Execute `docker --version` command
2. Verify Docker daemon accessibility
3. Test Docker command execution permissions
4. Update Docker status indicator (Green = Success, Red = Failed)

**Internet Connectivity Check Process**:
1. Perform network connectivity test (ping or HTTP request)
2. Verify external network access
3. Update Internet status indicator (Green = Success, Red = Failed)

**Status Update Logic**:
- Footer text changes from "Checking..." to "Ready" or appropriate error message
- Status circles update colors based on check results
- "Run Moodle" button remains disabled if any check fails
- "Run Moodle" button becomes enabled only when all checks pass

### 2.2 First-Time Container Launch Process

#### 2.2.1 Image Verification and Download
When user clicks "Run Moodle" button:

1. **Image Existence Check**:
   ```bash
   docker images | grep "wenkhairu/moodle-prototype:502-stable"
   ```

2. **Image Download (if not exists)**:
   ```bash
   docker pull wenkhairu/moodle-prototype:502-stable
   ```

3. **Download Progress Modal**:
   - Modal window overlay on main application
   - Progress bar with percentage indicator
   - Real-time download progress parsing from Docker output
   - Cancel option (optional enhancement)

4. **Download Completion**:
   - Re-verify image existence
   - Close progress modal
   - Proceed to container launch

#### 2.2.2 Container Creation and Launch
1. **Container Launch Command**:
   ```bash
   docker run -d -p 8080:8080 wenkhairu/moodle-prototype:502-stable
   ```

2. **Container ID Storage**:
   - Save returned container ID to `container.id` file
   - File location: Application directory root
   - File format: Plain text, single line

3. **Startup Waiting Modal**:
   - Modal window with loading animation
   - Message: "Starting Moodle, please wait..."
   - No user interaction options during this phase

#### 2.2.3 Credential Extraction Process
1. **Log Monitoring**:
   ```bash
   docker logs <container_id>
   ```

2. **Target Log Patterns**:
   - Admin Password: `Generated admin password: <password>`
   - Moodle URL: `Moodle is available at: http://localhost:8080`

3. **Credential Storage**:
   - Save extracted information to `moodle.txt` file
   - File format:
     ```
     password=<extracted_password>
     url=http://localhost:8080
     ```

4. **Polling Strategy**:
   - Check container logs every 2-3 seconds
   - Maximum wait time: 5 minutes
   - Error handling for timeout scenarios

#### 2.2.4 Launch Completion
1. **Modal Dismissal**:
   - Close startup waiting modal
   - Return to main application window

2. **Credential Display**:
   - Show extracted credentials in main window
   - Display format:
     ```
     Username: admin
     Password: <extracted_password>
     URL: http://localhost:8080
     ```

3. **Browser Launch Offer**:
   - Dialog: "Would you like to open Moodle in your browser?"
   - Options: "Yes" / "No" / "Open Browser" button
   - If accepted: Launch default browser with Moodle URL

4. **Button State Update**:
   - Change "Run Moodle" button to "Stop Moodle"
   - Button color change to indicate active state

### 2.3 Subsequent Application Launches

#### 2.3.1 Startup Process
1. Perform same health checks as first-time launch
2. Display same initial UI state
3. Enable "Run Moodle" button if health checks pass

#### 2.3.2 Container Restart Process
When user clicks "Run Moodle" button:

1. **Container ID Verification**:
   - Check for existence of `container.id` file
   - Validate container ID format
   - Verify container exists in Docker

2. **Container Restart Logic**:
   ```bash
   docker container <container_id> start
   ```

3. **Fallback to First-Time Process**:
   - If `container.id` file missing → Execute first-time launch process
   - If container ID invalid → Execute first-time launch process
   - If container doesn't exist → Execute first-time launch process

4. **Credential Display**:
   - Read credentials from `moodle.txt` file
   - Display credentials in main window
   - Offer browser launch option

### 2.4 Container Management

#### 2.4.1 Stop Container Process
When "Stop Moodle" button is clicked:

1. **Graceful Shutdown Command**:
   ```bash
   docker container <container_id> stop
   ```

2. **Shutdown Progress Indicator**:
   - Loading spinner or progress indicator
   - Status text: "Stopping Moodle..."
   - Disable button during operation

3. **Completion Actions**:
   - Button text changes back to "Run Moodle"
   - Hide credential display
   - Update footer status

#### 2.4.2 Application Exit Behavior
When user closes the application:

1. **Container Status Check**:
   - Check if container is running using `container.id`
   - Command: `docker container <container_id> inspect --format='{{.State.Running}}'`

2. **Cleanup Process**:
   - If container running → Execute stop container process
   - Wait for graceful shutdown completion
   - Exit application after cleanup

3. **Direct Exit**:
   - If no container running → Exit immediately

## 3. User Interface Specifications

### 3.1 Main Window Layout

#### 3.1.1 Window Properties
- **Dimensions**: 800x600 pixels (minimum)
- **Resizable**: Optional (recommended fixed size for v1)
- **Window Title**: "Moodle Prototype Manager v[version]"

#### 3.1.2 Layout Structure
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
│         [Open Browser Button]       │
│         (when container running)    │
│                                     │
├─────────────────────────────────────┤
│ FOOTER                              │
│ Status: Ready     Docker● Internet● │
└─────────────────────────────────────┘
```

#### 3.1.3 Component Details

**Header Section**:
- Moodle logo: 128x128 pixels, centered
- Version text: Below logo, centered, smaller font

**Main Action Button**:
- Size: 200x50 pixels
- States: "Run Moodle" (green) / "Stop Moodle" (red)
- Disabled state: Grayed out with disabled cursor

**Credentials Display**:
- Monospace font for password display
- Selectable text for copy operations
- Format:
  ```
  Username: admin
  Password: [generated_password]
  URL: http://localhost:8080
  ```

**Footer Status Bar**:
- Height: 30 pixels
- Left side: Status text
- Right side: Status indicators with 10px colored circles

### 3.2 Modal Windows

#### 3.2.1 Download Progress Modal
- **Dimensions**: 400x200 pixels
- **Modal Type**: Blocking overlay
- **Components**:
  - Title: "Downloading Moodle Image..."
  - Progress bar: Full width with percentage
  - Cancel button (optional)

#### 3.2.2 Startup Wait Modal
- **Dimensions**: 300x150 pixels
- **Modal Type**: Blocking overlay
- **Components**:
  - Loading spinner/animation
  - Text: "Starting Moodle, please wait..."
  - No user controls

## 4. Technical Architecture

### 4.1 Application Structure

#### 4.1.1 Go Backend Components
```
main.go                 // Application entry point
app.go                  // Wails application context
docker/
  ├── manager.go        // Docker operations
  ├── health.go         // Health check functions
  └── logs.go           // Log parsing utilities
storage/
  ├── files.go          // File I/O operations
  └── credentials.go    // Credential management
ui/
  ├── state.go          // UI state management
  └── events.go         // UI event handlers
```

#### 4.1.2 Frontend Structure
```
frontend/
├── index.html          // Main application UI
├── css/
│   └── styles.css      // Application styling
├── js/
│   ├── app.js          // Main application logic
│   ├── ui.js           // UI manipulation
│   └── events.js       // Event handling
└── assets/
    └── images/         // UI images and icons
```

### 4.2 Data Flow Architecture

#### 4.2.1 Application State Management
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   UI Layer  │◄──►│ Go Backend  │◄──►│   Docker    │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Events    │    │ File System │    │  Container  │
└─────────────┘    └─────────────┘    └─────────────┘
```

#### 4.2.2 File System Integration
- **container.id**: Container identifier storage
- **moodle.txt**: Credential storage
- **Application logs**: Debug and error information
- **Configuration**: Application settings (future enhancement)

### 4.3 Error Handling Strategy

#### 4.3.1 Docker Operation Errors
- **Docker not installed**: Clear error message with installation guidance
- **Docker not running**: Instructions to start Docker Desktop
- **Permission errors**: Guidance for Docker permissions
- **Image pull failures**: Network troubleshooting suggestions
- **Container launch failures**: Port conflict detection and resolution

#### 4.3.2 Network and Connectivity Errors
- **Internet connectivity**: Offline mode graceful degradation
- **Docker registry access**: Alternative download methods
- **Container port conflicts**: Port availability checking and alternatives

#### 4.3.3 File System Errors
- **Permission errors**: User guidance for file access permissions
- **Disk space**: Available space checking and warnings
- **File corruption**: Validation and recovery procedures

## 5. Performance Requirements

### 5.1 Response Time Targets
- **Application startup**: < 3 seconds
- **Health checks**: < 5 seconds total
- **Container operations**: Based on Docker performance
- **UI responsiveness**: < 100ms for user interactions

### 5.2 Resource Requirements
- **Memory usage**: < 50MB idle, < 100MB during operations
- **CPU usage**: Minimal idle, burst during Docker operations
- **Disk space**: < 10MB application, variable for Docker images
- **Network bandwidth**: Based on Docker image size

## 6. Security Considerations

### 6.1 Local Security
- **File permissions**: Restrict access to credential files
- **Docker socket access**: Validate Docker command execution
- **Process isolation**: Separate Docker operations from UI

### 6.2 Credential Management
- **Password storage**: Plain text in local file (acceptable for prototype)
- **File access**: Local user access only
- **Network exposure**: No network credential transmission

## 7. Testing Strategy

### 7.1 Unit Testing
- Docker operation functions
- File I/O operations
- Credential parsing logic
- UI state management

### 7.2 Integration Testing
- Complete application workflows
- Docker integration scenarios
- Cross-platform compatibility
- Error condition handling

### 7.3 User Acceptance Testing
- First-time user experience
- Subsequent launch scenarios
- Error recovery procedures
- Performance validation

## 8. Deployment and Distribution

### 8.1 Build Process
- Cross-platform compilation
- Asset bundling
- Executable packaging
- Code signing (for macOS)

### 8.2 Distribution Strategy
- Direct executable download
- Platform-specific packages
- Update mechanism (future enhancement)
- Installation instructions

## 9. Future Enhancements

### 9.1 Planned Features
- Configuration management
- Multiple container support
- Automatic updates
- Logging and diagnostics
- Container health monitoring

### 9.2 Platform Extensions
- Linux support
- Container registry alternatives
- Custom Docker image support
- Advanced Docker options

---

*This specification document serves as the definitive technical reference for the Moodle Prototype Manager application development.*
