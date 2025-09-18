package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger

// InitLogger initializes the logger to write to moodle.log
func InitLogger() {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return
	}

	// Create or open the log file
	logFile := filepath.Join(logDir, "moodle.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	// Create logger with timestamp
	logger = log.New(file, "", log.LstdFlags)
	
	// Log initialization
	LogInfo("=== Moodle Prototype Manager Started ===")
}

// LogInfo logs an info message
func LogInfo(message string) {
	logMessage("INFO", message)
}

// LogError logs an error message
func LogError(message string, err error) {
	if err != nil {
		logMessage("ERROR", fmt.Sprintf("%s: %v", message, err))
	} else {
		logMessage("ERROR", message)
	}
}

// LogDebug logs a debug message
func LogDebug(message string) {
	logMessage("DEBUG", message)
}

// LogWarning logs a warning message
func LogWarning(message string) {
	logMessage("WARNING", message)
}

// logMessage writes a formatted log message
func logMessage(level, message string) {
	if logger != nil {
		logger.Printf("[%s] %s", level, message)
	}
	// Also print to console for immediate feedback
	fmt.Printf("[%s] %s: %s\n", time.Now().Format("15:04:05"), level, message)
}