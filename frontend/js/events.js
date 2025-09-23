// Event handling and backend communication

// Import Wails bindings
import { HealthCheck, RunMoodle, StopMoodle, GetCredentials, OpenBrowser, GetImageName } from '../wailsjs/go/main/App.js';

// Add IsContainerReady manually until Wails regenerates properly
function IsContainerReady() {
    return window.go?.main?.App?.IsContainerReady?.() || Promise.resolve(false);
}

// Import shared state and functions from app.js
import { AppState, updateStatusText, updateHealthCheckResults } from './app.js';

// Import UI functions
import { 
    showNotification, hideCredentials, showDownloadModal, hideDownloadModal,
    showStartupModal, hideStartupModal, setUIEnabled, setButtonLoading, showBrowserDialog,
    hideBrowserDialog, displayCredentials
} from './ui.js';

// Wails bindings (will be available at runtime)
let wailsBindings = {};

// Check if running in Wails environment
function isWailsEnvironment() {
    return typeof window.go !== 'undefined' && window.go !== null;
}

// Initialize Wails bindings
function initializeWailsBindings() {
    if (isWailsEnvironment()) {
        console.log('Wails environment detected');
        wailsBindings = {
            HealthCheck: HealthCheck,
            RunMoodle: RunMoodle,
            StopMoodle: StopMoodle,
            GetCredentials: GetCredentials,
            OpenBrowser: OpenBrowser,
            IsContainerReady: IsContainerReady,
            GetImageName: GetImageName
        };
    } else {
        console.log('Running in development mode, using mock functions');
        wailsBindings = {
            HealthCheck: mockHealthCheck,
            RunMoodle: mockRunMoodle,
            StopMoodle: mockStopMoodle,
            GetCredentials: mockGetCredentials,
            OpenBrowser: mockOpenBrowser,
            IsContainerReady: mockIsContainerReady,
            GetImageName: mockGetImageName
        };
    }
}

// Mock functions for development
function mockHealthCheck() {
    return Promise.resolve({
        docker: Math.random() > 0.3,
        internet: Math.random() > 0.2
    });
}

function mockRunMoodle() {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve(null); // null means no error
        }, 5000);
    });
}

function mockStopMoodle() {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve(null);
        }, 2000);
    });
}

function mockGetCredentials() {
    return Promise.resolve({
        username: 'admin',
        password: 'admin123!@#',
        url: 'http://localhost:8080'
    });
}

function mockOpenBrowser() {
    return Promise.resolve(null);
}

function mockIsContainerReady() {
    return Promise.resolve(Math.random() > 0.5);
}

function mockGetImageName() {
    return Promise.resolve('wenkhairu/moodle-prototype:502-stable');
}

// Enhanced health check function with independent checking
async function performHealthChecks() {
    // Prevent overlapping health checks
    if (healthCheckInProgress) {
        console.log('Health check already in progress, skipping...');
        return;
    }
    
    healthCheckInProgress = true;
    console.log('Starting health checks...');
    
    updateStatusText('Checking systems...');
    
    // Add visual indication that check is in progress
    const dockerIndicator = document.getElementById('docker-status');
    const internetIndicator = document.getElementById('internet-status');
    
    // Temporarily show checking state
    if (dockerIndicator) dockerIndicator.className = 'status-circle checking';
    if (internetIndicator) internetIndicator.className = 'status-circle checking';
    
    try {
        console.log('Calling backend HealthCheck...');
        const healthStatus = await wailsBindings.HealthCheck();
        console.log('Backend returned:', healthStatus);
        
        // Each status is checked independently
        AppState.dockerStatus = healthStatus.docker || false;
        AppState.internetStatus = healthStatus.internet || false;
        
        console.log('Processed health check results:', {
            docker: AppState.dockerStatus,
            internet: AppState.internetStatus
        });
        
        updateHealthCheckResults();
        
        // Update status text based on results
        if (AppState.dockerStatus && AppState.internetStatus) {
            updateStatusText('All systems ready');
        } else if (!AppState.dockerStatus && !AppState.internetStatus) {
            updateStatusText('Docker and Internet unavailable');
        } else if (!AppState.dockerStatus) {
            updateStatusText('Docker unavailable');
        } else if (!AppState.internetStatus) {
            updateStatusText('Internet unavailable');
        }
        
    } catch (error) {
        console.error('Health check service failed:', error);
        
        // On service failure, maintain last known states or set to unknown
        // Don't automatically set both to false
        if (typeof AppState.dockerStatus === 'undefined') {
            AppState.dockerStatus = false;
        }
        if (typeof AppState.internetStatus === 'undefined') {
            AppState.internetStatus = false;
        }
        
        updateStatusText('Health check failed - ' + error.message);
        updateHealthCheckResults();
    } finally {
        healthCheckInProgress = false;
    }
}

// Enhanced container start function
async function startMoodleContainer() {
    console.log('Starting Moodle container...');
    
    try {
        setUIEnabled(false);
        
        // Show download modal
        showDownloadModal();
        updateStatusText('Starting Moodle container...');

        // Listen for real Docker pull progress events
        let progressListener = null;
        let hasReceivedProgress = false;

        // Set up progress listener before starting the pull
        progressListener = window.runtime.EventsOn('docker:pull:progress', (data) => {
            hasReceivedProgress = true;
            const percentage = data.percentage || 0;
            const status = data.status || 'Downloading...';

            // Update progress bar
            const progressFill = document.getElementById('download-progress');
            if (progressFill) {
                // Only update if we have valid percentage, otherwise keep current value
                if (percentage >= 0) {
                    progressFill.style.width = Math.min(percentage, 100) + '%';
                }
            }

            // Update status text in modal if available
            const modalStatus = document.querySelector('.modal-status');
            if (modalStatus) {
                modalStatus.textContent = status;
            }

            console.log(`Docker pull progress: ${percentage.toFixed(1)}% - ${status}`);
        });

        // Start the container (backend will handle image pull if needed)
        try {
            await wailsBindings.RunMoodle();

            // Complete download progress
            const progressFill = document.getElementById('download-progress');
            if (progressFill) {
                progressFill.style.width = '100%';
            }
        } catch (runError) {
            // Clean up event listener on error
            if (progressListener) {
                window.runtime.EventsOff('docker:pull:progress');
            }
            
            // Extract error message properly
            let errorMessage = 'Failed to start container';
            if (runError && typeof runError === 'string') {
                errorMessage = runError;
            } else if (runError && runError.message) {
                errorMessage = runError.message;
            } else if (runError && runError.toString) {
                errorMessage = runError.toString();
            }
            
            throw new Error(errorMessage);
        }

        // Clean up event listener
        if (progressListener) {
            window.runtime.EventsOff('docker:pull:progress');
        }

        // Update UI to show we're waiting for container
        hideDownloadModal();
        showStartupModal();
        updateStatusText('Waiting for container to start...');
        
        // Poll for container readiness by checking if moodle.txt file exists
        let attempts = 0;
        // No timeout limit - Windows installations can take 20-30+ minutes
        // The modal will stay open until the container is actually ready
        
        // Keep the modal open and update status
        updateStatusText('Container is starting, please wait...');
        
        let isReady = false;
        while (!isReady) {
            try {
                isReady = await wailsBindings.IsContainerReady();
                
                if (isReady) {
                    console.log('Container is ready!');
                    break;
                }
                
                // Update status with progress indication
                const minutes = Math.floor(attempts * 2 / 60);
                const seconds = (attempts * 2) % 60;
                updateStatusText(`Container starting... ${minutes}:${seconds.toString().padStart(2, '0')}`);
                
                // Wait 2 seconds before next check
                await new Promise(resolve => setTimeout(resolve, 2000));
                attempts++;
            } catch (error) {
                console.error('Error checking container readiness:', error);
                // Don't break on error, just continue checking
                await new Promise(resolve => setTimeout(resolve, 2000));
                attempts++;
            }
        }
        
        // Only hide modal after we're done checking
        hideStartupModal();
        
        // Container is ready (since we only exit the loop when isReady is true)
        // Get the credentials now that container is ready
        try {
            const credentials = await wailsBindings.GetCredentials();
            
            // Update state and UI
            AppState.containerRunning = true;
            AppState.credentials = credentials;
            updateStatusText('Container running - Moodle is ready!');
            displayCredentials(credentials);
            
            // Update button to stop mode
            const runButton = document.getElementById('run-moodle-btn');
            runButton.textContent = 'Stop Moodle';
            runButton.classList.add('stop');
            runButton.disabled = false;
            
            showNotification('Moodle container started successfully!');
            
            // Show browser dialog
            setTimeout(() => {
                showBrowserDialog();
            }, 1000);
        } catch (credError) {
            console.error('Failed to get credentials:', credError);
            throw new Error('Container started but failed to get credentials: ' + credError.message);
        }
        
    } catch (error) {
        console.error('Failed to start container:', error);
        
        // Extract error message safely
        let errorMessage = 'Unknown error';
        if (error && error.message) {
            errorMessage = error.message;
        } else if (error && typeof error === 'string') {
            errorMessage = error;
        } else if (error && error.toString) {
            errorMessage = error.toString();
        }
        
        showNotification('Failed to start Moodle: ' + errorMessage, 'error');
        
        // Clean up any running intervals
        if (typeof downloadInterval !== 'undefined') {
            clearInterval(downloadInterval);
        }
        
        // Hide modals and reset UI
        hideDownloadModal();
        hideStartupModal();
        updateStatusText('Container start failed');
        setUIEnabled(true);
    }
}

// Enhanced container stop function
async function stopMoodleContainer() {
    console.log('Stopping Moodle container...');
    
    try {
        setButtonLoading('run-moodle-btn', true, 'Stopping...');
        
        try {
            await wailsBindings.StopMoodle();
        } catch (stopError) {
            // Extract error message properly
            let errorMessage = 'Failed to stop container';
            if (stopError && typeof stopError === 'string') {
                errorMessage = stopError;
            } else if (stopError && stopError.message) {
                errorMessage = stopError.message;
            } else if (stopError && stopError.toString) {
                errorMessage = stopError.toString();
            }
            
            throw new Error(errorMessage);
        }
        
        // Container stopped successfully
        AppState.containerRunning = false;
        
        const runButton = document.getElementById('run-moodle-btn');
        runButton.textContent = 'Run Moodle';
        runButton.classList.remove('stop');
        runButton.disabled = false;
        
        hideCredentials();
        showNotification('Moodle container stopped successfully', 'success');
        
    } catch (error) {
        console.error('Failed to stop container:', error);
        showNotification('Failed to stop Moodle: ' + error.message, 'error');
        
        // Reset button state
        setButtonLoading('run-moodle-btn', false);
    }
}

// Load credentials from backend and check if container is running
async function loadCredentials() {
    try {
        // First check if we have stored credentials and if container might be running
        const credentials = await wailsBindings.GetCredentials();

        // If we have a valid password, it means credentials were saved previously
        if (credentials && credentials.password && credentials.password !== '') {
            console.log('Found stored credentials, checking if container is running...');

            // Check if container is actually ready/running
            try {
                const isReady = await wailsBindings.IsContainerReady();
                if (isReady) {
                    console.log('Container is running, displaying credentials');

                    // Update app state
                    AppState.containerRunning = true;
                    AppState.credentials = credentials;

                    // Update button to stop mode
                    const runButton = document.getElementById('run-moodle-btn');
                    if (runButton) {
                        runButton.textContent = 'Stop Moodle';
                        runButton.classList.add('stop');
                        runButton.disabled = false;
                    }

                    // Display credentials using the proper function
                    displayCredentials(credentials);
                    updateStatusText('Container running - Moodle is ready!');
                } else {
                    console.log('Stored credentials found but container is not running');
                    // Container is not running, clear the stored credentials or keep for next run
                    updateStatusText('Ready to run Moodle');
                }
            } catch (readyError) {
                console.error('Failed to check container readiness:', readyError);
                updateStatusText('Ready to run Moodle');
            }
        } else {
            console.log('No stored credentials found');
            updateStatusText('Ready to run Moodle');
        }

    } catch (error) {
        console.error('Failed to load credentials:', error);
        updateStatusText('Ready to run Moodle');
    }
}

// Load and display Docker image name
async function loadImageName() {
    try {
        const imageName = await wailsBindings.GetImageName();
        // Try both ID and class selector as fallback
        const versionElement = document.getElementById('version-text') || document.querySelector('.version-text');
        
        if (versionElement) {
            versionElement.textContent = imageName;
            console.log('Updated version text to:', imageName);
        } else {
            console.error('Version text element not found');
        }
        
        console.log('Loaded image name:', imageName);
        
    } catch (error) {
        console.error('Failed to load image name:', error);
        // Fallback to default text if loading fails
        const versionElement = document.getElementById('version-text') || document.querySelector('.version-text');
        if (versionElement) {
            versionElement.textContent = 'Moodle Prototype';
        }
    }
}

// Handle application startup
document.addEventListener('DOMContentLoaded', function() {
    // Initialize Wails bindings
    initializeWailsBindings();

    // Load and display image name
    setTimeout(async function() {
        console.log('Attempting to load image name...');
        await loadImageName();
    }, 500);

    // Perform initial health checks
    setTimeout(performHealthChecks, 1000);

    // Load credentials if container is already running
    setTimeout(loadCredentials, 1500);

    // Add event listener for copy password button
    const copyPasswordBtn = document.getElementById('copy-password-btn');
    if (copyPasswordBtn) {
        copyPasswordBtn.addEventListener('click', handleCopyPassword);
    }
});

// Handle keyboard shortcuts
document.addEventListener('keydown', function(event) {
    // Ctrl+R or Cmd+R to refresh health checks
    if ((event.ctrlKey || event.metaKey) && event.key === 'r') {
        event.preventDefault();
        performHealthChecks();
        showNotification('Refreshing health checks...', 'info');
    }
    
    // Escape key to close modals
    if (event.key === 'Escape') {
        hideBrowserDialog();
        // Don't close other modals as they represent ongoing operations
    }
});

// Handle window beforeunload (app closing)
window.addEventListener('beforeunload', function(event) {
    // If container is running, we might want to stop it
    // This will be handled by the Go backend's OnShutdown
    if (AppState.containerRunning) {
        console.log('Application closing with container running - backend will handle cleanup');
    }
});

// Track if health check is in progress to prevent overlapping
let healthCheckInProgress = false;

// Periodic health check (optional - every 30 seconds)
setInterval(function() {
    // Only perform periodic checks if not already checking
    if (!healthCheckInProgress) {
        performHealthChecks();
    }
}, 30000);

// Handle browser opening
async function handleBrowserYes() {
    hideBrowserDialog();

    try {
        await wailsBindings.OpenBrowser();
        showNotification('Opening browser...', 'success');
    } catch (error) {
        console.error('Failed to open browser:', error);
        showNotification('Failed to open browser: ' + error.message, 'error');
    }
}

// Handle copy password functionality
async function handleCopyPassword() {
    const passwordElement = document.getElementById('password');
    const copyButton = document.getElementById('copy-password-btn');

    if (!passwordElement || !copyButton) {
        console.error('Password element or copy button not found');
        return;
    }

    const password = passwordElement.textContent;

    if (!password || password === '-') {
        showNotification('No password available to copy', 'error');
        return;
    }

    try {
        // Use the Clipboard API if available
        if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(password);
        } else {
            // Fallback for older browsers or non-secure contexts
            const textArea = document.createElement('textarea');
            textArea.value = password;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            textArea.style.top = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
        }

        // Visual feedback
        const originalText = copyButton.innerHTML;
        const originalClass = copyButton.className;

        copyButton.innerHTML = '✓';
        copyButton.classList.add('copied');

        // Reset after 2 seconds
        setTimeout(() => {
            copyButton.innerHTML = originalText;
            copyButton.className = originalClass;
        }, 2000);

        showNotification('Password copied to clipboard!', 'success');

    } catch (error) {
        console.error('Failed to copy password:', error);
        showNotification('Failed to copy password', 'error');
    }
}

// Export functions for global access
window.performHealthChecks = performHealthChecks;
window.startMoodleContainer = startMoodleContainer;
window.stopMoodleContainer = stopMoodleContainer;
window.handleBrowserYes = handleBrowserYes;
window.handleCopyPassword = handleCopyPassword;

// Also export as modules
export { performHealthChecks, startMoodleContainer, stopMoodleContainer, handleBrowserYes, handleCopyPassword };