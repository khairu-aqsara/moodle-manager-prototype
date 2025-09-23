# Moodle Prototype Manager - Frontend Documentation

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [HTML Structure](#html-structure)
4. [CSS Design System](#css-design-system)
5. [JavaScript Components](#javascript-components)
6. [User Interface Components](#user-interface-components)
7. [State Management](#state-management)
8. [Event Handling](#event-handling)
9. [Modal System](#modal-system)
10. [Responsive Design](#responsive-design)
11. [Integration with Backend](#integration-with-backend)

## Overview

The frontend of the Moodle Prototype Manager is built using modern web technologies integrated with the Wails framework. It provides a clean, intuitive desktop application interface for managing Docker containers running Moodle prototypes.

### Technology Stack

- **HTML5**: Semantic markup with accessibility considerations
- **CSS3**: Modern styling with gradients, animations, and responsive design
- **JavaScript (ES6+)**: Modular architecture with import/export
- **Wails Integration**: Go backend communication via JavaScript bindings

### Design Principles

1. **Simplicity**: Clean, uncluttered interface focused on core functionality
2. **Consistency**: Unified visual language across all components
3. **Responsiveness**: Adapts to different window sizes gracefully
4. **Accessibility**: Proper semantic HTML and keyboard navigation
5. **Performance**: Efficient DOM manipulation and minimal resource usage

## Architecture

### File Structure

```
frontend/
├── index.html          # Main application layout
├── css/
│   └── styles.css     # Complete stylesheet with component styles
├── js/
│   ├── app.js         # Application state and initialization
│   ├── ui.js          # UI manipulation functions
│   └── events.js      # Event handling and backend integration
├── assets/
│   └── images/        # Static assets (logos, icons)
└── dist/              # Built frontend files (auto-generated)
```

### Module System

The frontend uses ES6 modules for clean separation of concerns:

```javascript
// Cross-module dependencies
app.js     ←→ ui.js     ←→ events.js
    ↓           ↓           ↓
AppState   UI Functions  Event Handlers
```

## HTML Structure

### Main Layout (`index.html`)

The application follows a standard desktop application layout:

```html
<!DOCTYPE html>
<html lang="en">
<body>
    <header class="header">
        <!-- Logo and version information -->
    </header>

    <main class="main-content">
        <!-- Primary action button and credentials display -->
    </main>

    <footer class="footer">
        <!-- Status indicators and health check results -->
    </footer>

    <!-- Modal overlays for operations -->
</body>
</html>
```

### Semantic Structure

#### Header Section
```html
<header class="header">
    <div class="logo-container">
        <img src="assets/images/moodle-logo.png" alt="Moodle Logo" class="logo">
        <h1 class="version-text" id="version-text">Moodle Prototype Manager v1.0</h1>
    </div>
</header>
```

**Features:**
- Fallback logo display using CSS gradient if image fails to load
- Dynamic version text updated via JavaScript
- Centered layout with consistent spacing

#### Main Content Area
```html
<main class="main-content">
    <div class="action-container">
        <button id="run-moodle-btn" class="action-button" disabled>
            Run Moodle
        </button>
    </div>

    <div id="credentials-display" class="credentials-container" style="display: none;">
        <!-- Credentials table -->
    </div>
</main>
```

**Features:**
- Primary action button with state management
- Collapsible credentials display
- Flexible layout adapting to content

#### Footer Status Bar
```html
<footer class="footer">
    <div class="status-left">
        <span id="status-text">Checking...</span>
    </div>
    <div class="status-right">
        <div class="status-indicator">
            <span class="status-circle red" id="docker-status"></span>
            <span class="status-label">Docker</span>
        </div>
        <div class="status-indicator">
            <span class="status-circle red" id="internet-status"></span>
            <span class="status-label">Internet</span>
        </div>
    </div>
</footer>
```

**Features:**
- Dynamic status text updates
- Color-coded health indicators
- Responsive layout with proper spacing

### Modal Components

#### Download Progress Modal
```html
<div id="download-modal" class="modal" style="display: none;">
    <div class="modal-content download-modal">
        <h3>Downloading Moodle Image...</h3>
        <div class="progress-container">
            <div class="progress-bar">
                <div class="progress-fill" id="download-progress"></div>
            </div>
            <div class="progress-labels">
                <span>0%</span>
                <span>100%</span>
            </div>
        </div>
    </div>
</div>
```

#### Startup Wait Modal
```html
<div id="startup-modal" class="modal" style="display: none;">
    <div class="modal-content startup-modal">
        <div class="loading-spinner"></div>
        <p>Starting Moodle, please wait...</p>
    </div>
</div>
```

#### Browser Confirmation Dialog
```html
<div id="browser-dialog" class="modal" style="display: none;">
    <div class="modal-content dialog-modal">
        <h3>Open Browser</h3>
        <p>Would you like to open Moodle in your browser?</p>
        <div class="dialog-buttons">
            <button id="browser-yes" class="dialog-button primary">Yes</button>
            <button id="browser-no" class="dialog-button secondary">No</button>
        </div>
    </div>
</div>
```

## CSS Design System

### Color Palette

```css
/* Primary Colors */
--primary-gradient: linear-gradient(135deg, #f98012 0%, #ff6b35 100%);
--primary-hover: linear-gradient(135deg, #e6720f 0%, #e55a2b 100%);

/* Status Colors */
--success-color: #28a745;
--error-color: #dc3545;
--warning-color: #ffc107;
--info-color: #17a2b8;

/* Neutral Colors */
--background-color: #f5f5f5;
--text-primary: #333;
--text-secondary: #666;
--text-muted: #999;
--border-color: #e0e0e0;
```

### Typography

```css
/* Font Stack */
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;

/* Size Scale */
--font-size-large: 16px;
--font-size-normal: 14px;
--font-size-small: 12px;
--font-size-tiny: 11px;

/* Weights */
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
```

### Spacing System

```css
/* Spacing Scale (based on 4px grid) */
--spacing-xs: 4px;
--spacing-sm: 8px;
--spacing-md: 15px;
--spacing-lg: 20px;
--spacing-xl: 30px;
```

### Component Styles

#### Buttons

**Primary Action Button:**
```css
.action-button {
    width: 140px;
    height: 36px;
    background: linear-gradient(135deg, #f98012 0%, #ff6b35 100%);
    border-radius: 20px;
    color: white;
    font-weight: 600;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 4px 15px rgba(249, 128, 18, 0.3);
}
```

**Features:**
- Gradient background with hover effects
- Smooth cubic-bezier transitions
- Shimmer effect on hover using pseudo-elements
- Disabled state handling
- Stop state styling (red gradient)

**Dialog Buttons:**
```css
.dialog-button.primary {
    background: linear-gradient(135deg, #f98012 0%, #ff6b35 100%);
    color: white;
    padding: 8px 16px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(249, 128, 18, 0.3);
}

.dialog-button.secondary {
    background-color: #6c757d;
    color: white;
    border-radius: 12px;
}
```

#### Status Indicators

**Status Circles:**
```css
.status-circle {
    width: 8px;
    height: 8px;
    border-radius: 50%;
}

.status-circle.red { background-color: #dc3545; }
.status-circle.green { background-color: #28a745; }
.status-circle.checking {
    background-color: #ffc107;
    animation: pulse 1.5s infinite;
}
```

**Pulse Animation:**
```css
@keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
}
```

#### Credentials Table

```css
.credentials-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
    border: 1px solid #ddd;
}

.credentials-table td.label {
    font-weight: 500;
    background: #f9f9f9;
    width: 25%;
}

.credentials-table td.value {
    font-family: monospace;
    background: white;
}
```

**Copy Button:**
```css
.copy-button {
    width: 24px;
    height: 20px;
    background: #f8f9fa;
    border: 1px solid #dee2e6;
    border-radius: 4px;
    transition: all 0.2s ease;
}

.copy-button.copied {
    background: #d4edda;
    border-color: #c3e6cb;
    color: #155724;
}
```

### Modal System Styling

**Base Modal:**
```css
.modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.modal-content {
    background: white;
    border-radius: 12px;
    padding: 30px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}
```

**Progress Bar:**
```css
.progress-bar {
    width: 100%;
    height: 8px;
    background-color: #e9ecef;
    border-radius: 4px;
    overflow: hidden;
}

.progress-fill {
    height: 100%;
    background-color: #007bff;
    transition: width 0.3s ease;
    width: 0%;
}
```

**Loading Spinner:**
```css
.loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid #e3e3e3;
    border-top: 4px solid #007bff;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}
```

## JavaScript Components

### Application State Management (`app.js`)

#### AppState Object

Central state management for the entire application:

```javascript
export const AppState = {
    dockerStatus: false,      // Docker daemon accessibility
    internetStatus: false,    // Internet connectivity status
    containerRunning: false,  // Current container state
    credentials: {
        username: 'admin',
        password: '',         // Extracted from logs
        url: 'http://localhost:8080'
    }
};
```

#### Core Functions

**`updateHealthCheckResults()`**
- Updates UI elements based on AppState
- Controls button enable/disable logic
- Updates status indicator colors
- Sets appropriate status messages

**`updateStatusText(text)`**
- Updates footer status text
- Provides user feedback for operations

**Event Listener Setup:**
```javascript
document.addEventListener('DOMContentLoaded', function() {
    initializeUI();
    setupEventListeners();
});
```

### UI Management (`ui.js`)

#### Modal Control Functions

**Download Progress Modal:**
```javascript
export function showDownloadModal() {
    const modal = document.getElementById('download-modal');
    if (modal) modal.style.display = 'flex';
}

export function updateDownloadProgress(percentage, status) {
    const progressFill = document.getElementById('download-progress');
    const statusText = document.querySelector('.modal-status');

    if (progressFill) {
        progressFill.style.width = `${percentage}%`;
    }
    if (statusText) {
        statusText.textContent = status;
    }
}
```

**Startup Wait Modal:**
```javascript
export function showStartupModal() {
    const modal = document.getElementById('startup-modal');
    if (modal) modal.style.display = 'flex';
}
```

**Browser Confirmation Dialog:**
```javascript
export function showBrowserDialog() {
    const modal = document.getElementById('browser-dialog');
    if (modal) modal.style.display = 'flex';
}
```

#### Credential Display Management

**`displayCredentials(credentials)`**
- Shows the credentials table
- Updates password and URL fields
- Handles missing or invalid data gracefully

```javascript
export function displayCredentials(credentials) {
    const credentialsDisplay = document.getElementById('credentials-display');
    const passwordElement = document.getElementById('password');
    const urlElement = document.getElementById('url');

    if (credentialsDisplay) {
        credentialsDisplay.style.display = 'block';
    }

    if (passwordElement) {
        passwordElement.textContent = credentials.password || '-';
    }

    if (urlElement) {
        const url = credentials.url || 'http://localhost:8080';
        urlElement.textContent = url;
        urlElement.href = url;
    }
}
```

#### Notification System

**`showNotification(message, type, duration)`**
- Creates temporary notification overlays
- Supports multiple types: success, error, warning, info
- Auto-dismisses after specified duration
- Positioned at top-right of window

```javascript
export function showNotification(message, type = 'info', duration = 5000) {
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;

    // Dynamic styling based on type
    switch (type) {
        case 'success':
            notification.style.backgroundColor = '#28a745';
            break;
        case 'error':
            notification.style.backgroundColor = '#dc3545';
            break;
        // ... other types
    }

    document.body.appendChild(notification);
    setTimeout(() => {
        if (notification.parentNode) {
            notification.parentNode.removeChild(notification);
        }
    }, duration);
}
```

#### Button State Management

**`setButtonLoading(buttonId, isLoading, loadingText)`**
- Manages loading states for buttons
- Preserves original button text
- Disables button during loading
- Restores state when complete

### Event Handling (`events.js`)

#### Wails Backend Integration

The events module handles all communication with the Go backend through Wails bindings:

**Health Check Integration:**
```javascript
async function performHealthChecks() {
    try {
        const healthStatus = await window.go.main.App.HealthCheck();
        AppState.dockerStatus = healthStatus.docker;
        AppState.internetStatus = healthStatus.internet;
        updateHealthCheckResults();
    } catch (error) {
        console.error('Health check failed:', error);
        // Handle error state
    }
}
```

**Container Operations:**
```javascript
async function startMoodleContainer() {
    try {
        setButtonLoading('run-moodle-btn', true, 'Starting...');
        await window.go.main.App.RunMoodle();
        // Handle success
    } catch (error) {
        console.error('Failed to start container:', error);
        showNotification('Failed to start Moodle: ' + error, 'error');
    } finally {
        setButtonLoading('run-moodle-btn', false);
    }
}

async function stopMoodleContainer() {
    try {
        setButtonLoading('run-moodle-btn', true, 'Stopping...');
        await window.go.main.App.StopMoodle();
        // Handle success
    } catch (error) {
        console.error('Failed to stop container:', error);
        showNotification('Failed to stop Moodle: ' + error, 'error');
    } finally {
        setButtonLoading('run-moodle-btn', false);
    }
}
```

#### Progress Event Handling

**Docker Pull Progress:**
```javascript
// Listen for backend progress events
window.wails.Events.On('docker:pull:progress', (data) => {
    const { percentage, status } = data;
    updateDownloadProgress(percentage, status);
});
```

#### User Interaction Handlers

**Main Button Click Handler:**
```javascript
function handleRunMoodleClick() {
    if (!AppState.containerRunning) {
        startMoodleContainer();
    } else {
        stopMoodleContainer();
    }
}
```

**Browser Dialog Handlers:**
```javascript
function handleBrowserYes() {
    window.go.main.App.OpenBrowser()
        .then(() => {
            console.log('Browser opened successfully');
            hideBrowserDialog();
        })
        .catch(error => {
            console.error('Failed to open browser:', error);
            showNotification('Failed to open browser', 'error');
        });
}
```

**Copy Password Functionality:**
```javascript
document.getElementById('copy-password-btn')?.addEventListener('click', async () => {
    const passwordElement = document.getElementById('password');
    if (passwordElement && navigator.clipboard) {
        try {
            await navigator.clipboard.writeText(passwordElement.textContent);
            // Visual feedback for successful copy
            const button = document.getElementById('copy-password-btn');
            button.classList.add('copied');
            setTimeout(() => {
                button.classList.remove('copied');
            }, 2000);
        } catch (error) {
            console.error('Failed to copy password:', error);
        }
    }
});
```

## User Interface Components

### Primary Action Button

**States:**
1. **Disabled**: Gray, unclickable when health checks fail
2. **Run State**: Orange gradient, "Run Moodle" text
3. **Stop State**: Red gradient, "Stop Moodle" text
4. **Loading**: Disabled with loading text

**Visual Effects:**
- Hover shimmer animation
- Smooth state transitions
- Box shadow depth changes
- Transform animations on interaction

### Status Indicators

**Health Check Circles:**
- **Red**: Failed or unavailable
- **Green**: Healthy and available
- **Yellow**: Checking (with pulse animation)

**Implementation:**
```javascript
function updateHealthIndicator(elementId, status) {
    const element = document.getElementById(elementId);
    if (element) {
        element.classList.remove('red', 'green', 'checking');
        element.classList.add(status ? 'green' : 'red');
    }
}
```

### Credentials Display

**Features:**
- Tabular layout with label/value pairs
- Monospace font for password readability
- Copy button with visual feedback
- Clickable URL with backend integration

**Security Considerations:**
- Password is displayed in plain text (appropriate for local development tool)
- Copy functionality uses modern Clipboard API
- URL clicks are handled by backend for proper browser launching

### Progress Visualization

**Download Progress Bar:**
- Smooth width transitions
- Percentage labels at corners
- Status text updates
- Responsive to real-time backend events

**Loading Spinner:**
- CSS-only animation
- Consistent with design system colors
- Appropriate size for modal context

## State Management

### AppState Pattern

The application uses a centralized state object that's imported across modules:

```javascript
// Shared state accessible from any module
import { AppState } from './app.js';

// State updates trigger UI updates
function updateContainerState(running) {
    AppState.containerRunning = running;
    updateButtonState();
    if (running) {
        showCredentials();
    } else {
        hideCredentials();
    }
}
```

### State Synchronization

**Backend → Frontend:**
- Health check results update AppState
- Container status changes trigger UI updates
- Credentials are fetched and cached in AppState

**Frontend → Backend:**
- User actions call backend methods
- UI state reflects backend operation status
- Error states are handled gracefully

### Persistent State

While most state is ephemeral (reloaded on app restart), some state is persisted by the backend:
- Container ID (in `container.id` file)
- Moodle credentials (in `moodle.txt` file)
- Docker image configuration (in `image.docker` file)

## Event Handling

### Event Flow Architecture

```
User Action → UI Event → Backend Call → Backend Response → UI Update
     ↓                      ↑                      ↑           ↓
 Button Click → handleRunMoodleClick() → App.RunMoodle() → Success/Error → Update Button State
```

### Error Handling Strategy

**Three-Layer Error Handling:**

1. **Backend Layer**: Go functions return errors with context
2. **Integration Layer**: JavaScript catches and processes backend errors
3. **UI Layer**: Error states are communicated to users via notifications and status updates

**Example Error Flow:**
```javascript
try {
    await window.go.main.App.RunMoodle();
    // Success path
    showNotification('Moodle started successfully', 'success');
    updateContainerState(true);
} catch (error) {
    // Error path
    console.error('Container start failed:', error);
    showNotification(`Failed to start Moodle: ${error.message}`, 'error');
    updateContainerState(false);
} finally {
    // Cleanup path
    setButtonLoading('run-moodle-btn', false);
}
```

### Async Operation Handling

**Pattern for Long-Running Operations:**
1. Show loading state immediately
2. Display appropriate modal (progress/waiting)
3. Handle backend events for progress updates
4. Update UI based on final result
5. Clean up loading states and modals

## Modal System

### Modal Architecture

**Base Modal Structure:**
- Overlay background with opacity transition
- Centered content container
- Specific sizing for different modal types
- Z-index stacking for proper layering

**Modal Types:**

1. **Progress Modal** (400x200px)
   - Download progress bar
   - Real-time status updates
   - No user interaction required

2. **Wait Modal** (300x150px)
   - Loading spinner
   - Simple message
   - Non-dismissible during operation

3. **Dialog Modal** (280px width)
   - User decision required
   - Primary/secondary button options
   - Dismissible by user action

### Modal Management

**Show/Hide Pattern:**
```javascript
function showModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = 'flex';
        // Optional: focus management for accessibility
    }
}

function hideModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = 'none';
    }
}
```

**Modal State Management:**
- Only one modal visible at a time
- Modals are hidden by default (CSS: `display: none`)
- No modal stacking in current implementation

## Responsive Design

### Breakpoint Strategy

**Single Breakpoint Approach:**
```css
@media (max-width: 600px) {
    /* Mobile/narrow window adjustments */
    .credentials-container {
        min-width: 300px;
        margin: 0 20px;
    }

    .modal-content {
        margin: 0 20px;
        max-width: calc(100% - 40px);
    }
}
```

### Flexible Layout

**CSS Grid and Flexbox:**
- Header: Flexbox for center alignment
- Main: Flexbox column with center alignment
- Footer: Flexbox with space-between
- Modals: Flexbox for center positioning

**Adaptive Components:**
- Buttons maintain size but stack on narrow screens
- Tables scroll horizontally if needed
- Modals scale with viewport constraints
- Text scales appropriately with system font settings

### Window Size Considerations

**Desktop Application Context:**
- Fixed minimum window size (800x600)
- Optimized for desktop interaction patterns
- No touch-specific interactions
- Keyboard navigation support

## Integration with Backend

### Wails Framework Integration

**Binding Pattern:**
```javascript
// Backend methods are available at:
window.go.main.App.MethodName()

// Examples:
await window.go.main.App.HealthCheck()
await window.go.main.App.RunMoodle()
await window.go.main.App.GetCredentials()
```

**Event System:**
```javascript
// Listen for backend events
window.wails.Events.On('event-name', (data) => {
    // Handle event data
});

// Backend can emit events like:
// docker:pull:progress
// container:status:changed
```

### Error Boundary Implementation

**Global Error Handling:**
```javascript
window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason);
    showNotification('An unexpected error occurred', 'error');
});

window.addEventListener('error', (event) => {
    console.error('JavaScript error:', event.error);
    showNotification('Application error occurred', 'error');
});
```

### Development vs Production

**Environment Detection:**
```javascript
// Check if running in Wails context
const isProduction = window.go && window.wails;

if (isProduction) {
    // Use Wails bindings
    await window.go.main.App.HealthCheck();
} else {
    // Development fallbacks or mock data
    console.log('Running in development mode');
}
```

## Performance Considerations

### DOM Manipulation Optimization

- Minimal DOM queries (cache element references)
- Batch DOM updates when possible
- Use CSS transitions instead of JavaScript animations
- Efficient event listener management

### Memory Management

- Remove event listeners when not needed
- Clean up intervals and timeouts
- Avoid memory leaks in notification system
- Proper modal cleanup

### Asset Optimization

- Efficient CSS (minimal specificity)
- Single CSS file to reduce requests
- Optimized images with fallbacks
- Minimal JavaScript bundle size

## Accessibility Features

### Semantic HTML

- Proper heading hierarchy
- Descriptive alt text for images
- Semantic button and input elements
- ARIA labels where appropriate

### Keyboard Navigation

- Tab order follows logical flow
- Enter key activates primary buttons
- Escape key closes modals (where appropriate)
- Focus management in modal dialogs

### Visual Accessibility

- High contrast color ratios
- Clear visual hierarchy
- Consistent interactive element styling
- Loading states with appropriate feedback

This frontend documentation provides comprehensive coverage of all visual and interactive components of the Moodle Prototype Manager application, serving as a complete reference for frontend development and maintenance.