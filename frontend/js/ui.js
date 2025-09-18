// UI manipulation functions

// Import AppState from app.js
import { AppState } from './app.js';

// Show download modal
export function showDownloadModal() {
    const modal = document.getElementById('download-modal');
    if (modal) {
        modal.style.display = 'flex';
    }
}

// Hide download modal
export function hideDownloadModal() {
    const modal = document.getElementById('download-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}


// Show startup modal
export function showStartupModal() {
    const modal = document.getElementById('startup-modal');
    if (modal) {
        modal.style.display = 'flex';
    }
}

// Hide startup modal
export function hideStartupModal() {
    const modal = document.getElementById('startup-modal');
    if (modal) {
        modal.style.display = 'none';
    }
}


// Update UI when container is running
function updateUIForRunningContainer() {
    const runButton = document.getElementById('run-moodle-btn');
    
    // Update button
    runButton.textContent = 'Stop Moodle';
    runButton.classList.add('stop');
    runButton.disabled = false;
    
    // Show credentials
    showCredentials();
    
    // Update app state
    AppState.containerRunning = true;
}

// Show credentials display
export function showCredentials() {
    const credentialsDisplay = document.getElementById('credentials-display');
    const passwordElement = document.getElementById('password');
    const urlElement = document.getElementById('url');
    
    if (credentialsDisplay) {
        credentialsDisplay.style.display = 'block';
    }
    
    if (passwordElement) {
        passwordElement.textContent = AppState.credentials.password;
    }
    
    if (urlElement) {
        urlElement.textContent = AppState.credentials.url;
    }
}

// Hide credentials display
export function hideCredentials() {
    const credentialsDisplay = document.getElementById('credentials-display');
    if (credentialsDisplay) {
        credentialsDisplay.style.display = 'none';
    }
}

// Show browser dialog
export function showBrowserDialog() {
    const modal = document.getElementById('browser-dialog');
    if (modal) {
        modal.style.display = 'flex';
    }
}

// Hide browser dialog
export function hideBrowserDialog() {
    const modal = document.getElementById('browser-dialog');
    if (modal) {
        modal.style.display = 'none';
    }
}

// Update button loading state
export function setButtonLoading(buttonId, isLoading, loadingText = 'Loading...') {
    const button = document.getElementById(buttonId);
    if (button) {
        if (isLoading) {
            button.disabled = true;
            button.setAttribute('data-original-text', button.textContent);
            button.textContent = loadingText;
        } else {
            button.disabled = false;
            const originalText = button.getAttribute('data-original-text');
            if (originalText) {
                button.textContent = originalText;
                button.removeAttribute('data-original-text');
            }
        }
    }
}

// Show notification (can be used for errors or success messages)
export function showNotification(message, type = 'info', duration = 5000) {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    
    // Style the notification
    notification.style.position = 'fixed';
    notification.style.top = '20px';
    notification.style.right = '20px';
    notification.style.padding = '12px 20px';
    notification.style.borderRadius = '6px';
    notification.style.color = 'white';
    notification.style.fontSize = '14px';
    notification.style.zIndex = '10000';
    notification.style.maxWidth = '300px';
    notification.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
    
    // Set background color based on type
    switch (type) {
        case 'success':
            notification.style.backgroundColor = '#28a745';
            break;
        case 'error':
            notification.style.backgroundColor = '#dc3545';
            break;
        case 'warning':
            notification.style.backgroundColor = '#ffc107';
            notification.style.color = '#212529';
            break;
        default:
            notification.style.backgroundColor = '#17a2b8';
    }
    
    // Add to document
    document.body.appendChild(notification);
    
    // Remove after duration
    setTimeout(() => {
        if (notification.parentNode) {
            notification.parentNode.removeChild(notification);
        }
    }, duration);
}


// Display credentials in the UI
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

// Handle URL click - use backend OpenBrowser instead of direct navigation
window.handleUrlClick = function(event) {
    event.preventDefault();
    
    // Call the backend OpenBrowser function
    if (window.go?.main?.App?.OpenBrowser) {
        window.go.main.App.OpenBrowser()
            .then(() => {
                console.log('Browser opened successfully');
            })
            .catch(error => {
                console.error('Failed to open browser:', error);
                // Fallback to direct navigation
                window.open(event.target.href, '_blank');
            });
    } else {
        // Fallback for development
        window.open(event.target.href, '_blank');
    }
}

// Utility function to enable/disable UI elements
export function setUIEnabled(enabled) {
    const runButton = document.getElementById('run-moodle-btn');
    const openBrowserBtn = document.getElementById('open-browser-btn');
    
    if (runButton) {
        runButton.disabled = !enabled;
    }
    
    if (openBrowserBtn) {
        openBrowserBtn.disabled = !enabled;
    }
}