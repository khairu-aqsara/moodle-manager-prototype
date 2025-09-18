package docker

import (
	"regexp"
	"strings"
)

// CredentialInfo holds extracted credential information
type CredentialInfo struct {
	Password string `json:"password"`
	URL      string `json:"url"`
}

// LogParser handles parsing container logs for credentials
type LogParser struct {
	passwordRegex *regexp.Regexp
	urlRegex      *regexp.Regexp
}

// NewLogParser creates a new log parser
func NewLogParser() *LogParser {
	return &LogParser{
		// Match both patterns: "Password: " and "Generated admin password: "
		passwordRegex: regexp.MustCompile(`(?:Generated admin password:|Password:)\s*(.+)`),
		urlRegex:      regexp.MustCompile(`Moodle is available at:\s*(.+)`),
	}
}

// ExtractCredentials parses container logs to extract admin credentials
func (lp *LogParser) ExtractCredentials(logs string) *CredentialInfo {
	creds := &CredentialInfo{}
	
	// Extract password
	if matches := lp.passwordRegex.FindStringSubmatch(logs); len(matches) > 1 {
		creds.Password = strings.TrimSpace(matches[1])
	}
	
	// Extract URL
	if matches := lp.urlRegex.FindStringSubmatch(logs); len(matches) > 1 {
		creds.URL = strings.TrimSpace(matches[1])
	}
	
	return creds
}

// IsCredentialComplete checks if we have all required credentials
func (ci *CredentialInfo) IsComplete() bool {
	return ci.Password != "" && ci.URL != ""
}

// HasPassword checks if password is extracted
func (ci *CredentialInfo) HasPassword() bool {
	return ci.Password != ""
}

// HasURL checks if URL is extracted
func (ci *CredentialInfo) HasURL() bool {
	return ci.URL != ""
}