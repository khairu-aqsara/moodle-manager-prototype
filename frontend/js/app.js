// Main application JavaScript

// Import UI functions
import { showBrowserDialog, hideBrowserDialog } from './ui.js';

// Import container management functions from events.js
// Note: This will be available after events.js loads
let startMoodleContainer, stopMoodleContainer;

// Application state
export const AppState = {
    dockerStatus: false,
    internetStatus: false,
    containerRunning: false,
    credentials: {
        username: 'admin',
        password: '',
        url: 'http://localhost:8080'
    }
};

// Initialize application
document.addEventListener('DOMContentLoaded', function() {
    console.log('Moodle Prototype Manager starting...');
    
    // Initialize UI
    initializeUI();
    
    // Health checks will be handled by events.js
    
    // Set up event listeners
    setupEventListeners();
});

// Initialize UI elements
function initializeUI() {
    console.log('UI initialized');
}

// Note: Health check functions are now handled by events.js

// Update health check results in UI
export function updateHealthCheckResults() {
    const dockerCircle = document.getElementById('docker-status');
    const internetCircle = document.getElementById('internet-status');
    const runButton = document.getElementById('run-moodle-btn');
    const statusText = document.getElementById('status-text');
    
    // Update Docker status
    if (dockerCircle) {
        dockerCircle.classList.remove('checking', 'red', 'green');
        if (AppState.dockerStatus) {
            dockerCircle.classList.add('green');
        } else {
            dockerCircle.classList.add('red');
        }
    }
    
    // Update Internet status
    if (internetCircle) {
        internetCircle.classList.remove('checking', 'red', 'green');
        if (AppState.internetStatus) {
            internetCircle.classList.add('green');
        } else {
            internetCircle.classList.add('red');
        }
    }
    
    // Update button and status text
    if (AppState.dockerStatus && AppState.internetStatus) {
        runButton.disabled = false;
        statusText.textContent = 'Ready';
    } else {
        runButton.disabled = true;
        let errorMessage = 'Error: ';
        if (!AppState.dockerStatus) errorMessage += 'Docker unavailable ';
        if (!AppState.internetStatus) errorMessage += 'Internet unavailable';
        statusText.textContent = errorMessage.trim();
    }
}

// Set up event listeners
function setupEventListeners() {
    const runButton = document.getElementById('run-moodle-btn');
    const openBrowserBtn = document.getElementById('open-browser-btn');
    const browserYes = document.getElementById('browser-yes');
    const browserNo = document.getElementById('browser-no');
    
    if (runButton) {
        runButton.addEventListener('click', handleRunMoodleClick);
    }
    
    if (openBrowserBtn) {
        openBrowserBtn.addEventListener('click', handleOpenBrowserClick);
    }
    
    if (browserYes) {
        browserYes.addEventListener('click', handleBrowserYes);
    }
    
    if (browserNo) {
        browserNo.addEventListener('click', handleBrowserNo);
    }
    
    console.log('Event listeners set up');
}

// Handle Run Moodle button click
function handleRunMoodleClick() {
    const runButton = document.getElementById('run-moodle-btn');
    
    if (!AppState.containerRunning) {
        // Start container
        if (window.startMoodleContainer) {
            window.startMoodleContainer();
        } else {
            console.error('startMoodleContainer not available yet');
        }
    } else {
        // Stop container
        if (window.stopMoodleContainer) {
            window.stopMoodleContainer();
        } else {
            console.error('stopMoodleContainer not available yet');
        }
    }
}

// Handle Open Browser button click
function handleOpenBrowserClick() {
    showBrowserDialog();
}

// Handle browser dialog Yes
function handleBrowserYes() {
    if (window.handleBrowserYes) {
        window.handleBrowserYes();
    } else {
        console.error('handleBrowserYes not available from events.js');
    }
}

// Handle browser dialog No
function handleBrowserNo() {
    hideBrowserDialog();
}

// Note: Container start/stop functions are handled by events.js

// Update status text
export function updateStatusText(text) {
    const statusText = document.getElementById('status-text');
    if (statusText) {
        statusText.textContent = text;
    }
}