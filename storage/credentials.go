package storage

import (
	"fmt"
)

// Credentials represents Moodle login credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

// DefaultCredentials returns default credential structure
func DefaultCredentials() *Credentials {
	return &Credentials{
		Username: "admin",
		Password: "",
		URL:      "http://localhost:8080",
	}
}

// CredentialManager handles credential operations
type CredentialManager struct {
	fileManager *FileManager
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager() *CredentialManager {
	return &CredentialManager{
		fileManager: NewFileManager(),
	}
}

// Save saves credentials to file
func (cm *CredentialManager) Save(creds *Credentials) error {
	return cm.fileManager.SaveCredentials(creds.Password, creds.URL)
}

// Load loads credentials from file
func (cm *CredentialManager) Load() (*Credentials, error) {
	if !cm.fileManager.CredentialsExist() {
		return DefaultCredentials(), nil
	}
	
	data, err := cm.fileManager.LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}
	
	creds := DefaultCredentials()
	
	if password, exists := data["password"]; exists {
		creds.Password = password
	}
	
	if url, exists := data["url"]; exists {
		creds.URL = url
	}
	
	return creds, nil
}

// Update updates existing credentials
func (cm *CredentialManager) Update(password, url string) error {
	creds := &Credentials{
		Username: "admin",
		Password: password,
		URL:      url,
	}
	
	return cm.Save(creds)
}

// Clear removes stored credentials
func (cm *CredentialManager) Clear() error {
	return cm.fileManager.DeleteCredentials()
}

// Exists checks if credentials are stored
func (cm *CredentialManager) Exists() bool {
	return cm.fileManager.CredentialsExist()
}

// IsValid checks if credentials are valid (non-empty password and URL)
func (creds *Credentials) IsValid() bool {
	return creds.Password != "" && creds.URL != ""
}

// ToMap converts credentials to map for JSON serialization
func (creds *Credentials) ToMap() map[string]string {
	return map[string]string{
		"username": creds.Username,
		"password": creds.Password,
		"url":      creds.URL,
	}
}