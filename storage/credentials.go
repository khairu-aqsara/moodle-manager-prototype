package storage

import (
	"moodle-prototype-manager/errors"
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
	if creds == nil {
		return errors.NewValidationError("credentials", "credentials object cannot be nil", creds)
	}

	// Validate credentials before saving
	if !creds.IsValid() {
		return errors.NewValidationError("credentials", "credentials are invalid (missing password or URL)", creds)
	}

	err := cm.fileManager.SaveCredentials(creds.Password, creds.URL)
	if err != nil {
		return errors.WrapWithContext(err, "failed to save credentials to file")
	}

	return nil
}

// Load loads credentials from file
func (cm *CredentialManager) Load() (*Credentials, error) {
	if !cm.fileManager.CredentialsExist() {
		// Return default credentials when file doesn't exist (first run)
		return DefaultCredentials(), nil
	}

	data, err := cm.fileManager.LoadCredentials()
	if err != nil {
		return nil, errors.WrapWithContext(err, "failed to load credentials from file")
	}

	creds := DefaultCredentials()

	// Extract password with validation
	if password, exists := data["password"]; exists {
		if password == "" {
			return nil, errors.NewValidationError("password", "password field exists but is empty in credentials file", password)
		}
		creds.Password = password
	}

	// Extract URL with validation
	if url, exists := data["url"]; exists {
		if url == "" {
			return nil, errors.NewValidationError("url", "url field exists but is empty in credentials file", url)
		}
		creds.URL = url
	}

	return creds, nil
}

// Update updates existing credentials
func (cm *CredentialManager) Update(password, url string) error {
	// Validate input parameters
	if err := errors.ValidateNotEmpty("password", password); err != nil {
		return errors.WrapWithContext(err, "invalid password provided to Update")
	}
	if err := errors.ValidateNotEmpty("url", url); err != nil {
		return errors.WrapWithContext(err, "invalid URL provided to Update")
	}

	creds := &Credentials{
		Username: "admin",
		Password: password,
		URL:      url,
	}

	return cm.Save(creds)
}

// Clear removes stored credentials
func (cm *CredentialManager) Clear() error {
	err := cm.fileManager.DeleteCredentials()
	if err != nil {
		return errors.WrapWithContext(err, "failed to clear stored credentials")
	}
	return nil
}

// Exists checks if credentials are stored
func (cm *CredentialManager) Exists() bool {
	return cm.fileManager.CredentialsExist()
}

// IsValid checks if credentials are valid (non-empty password and URL)
func (creds *Credentials) IsValid() bool {
	if creds == nil {
		return false
	}
	return creds.Password != "" && creds.URL != ""
}

// Validate performs comprehensive validation and returns detailed errors
func (creds *Credentials) Validate() error {
	if creds == nil {
		return errors.NewValidationError("credentials", "credentials object is nil", creds)
	}

	multiErr := errors.NewMultiError("credential validation")

	if creds.Username == "" {
		multiErr.Add(errors.NewValidationError("username", "username cannot be empty", creds.Username))
	}

	if creds.Password == "" {
		multiErr.Add(errors.NewValidationError("password", "password cannot be empty", creds.Password))
	}

	if creds.URL == "" {
		multiErr.Add(errors.NewValidationError("url", "URL cannot be empty", creds.URL))
	}

	return multiErr.ToError()
}

// ToMap converts credentials to map for JSON serialization
func (creds *Credentials) ToMap() map[string]string {
	if creds == nil {
		// Return default values if credentials is nil
		return DefaultCredentials().ToMap()
	}

	return map[string]string{
		"username": creds.Username,
		"password": creds.Password,
		"url":      creds.URL,
	}
}

// IsComplete checks if credentials have both password and URL (used by log parser)
func (creds *Credentials) IsComplete() bool {
	if creds == nil {
		return false
	}
	return creds.Password != "" && creds.URL != ""
}