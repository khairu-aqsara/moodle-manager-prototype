package docker

import (
	"testing"
)

func TestLogParser(t *testing.T) {
	parser := NewLogParser()
	
	// Test log content that contains credentials
	testLog := `
Starting Moodle container...
Initializing database...
Generated admin password: MySecurePassword123!
Setting up web server...
Moodle is available at: http://localhost:8080
Container ready for connections.
	`
	
	creds := parser.ExtractCredentials(testLog)
	
	if creds.Password != "MySecurePassword123!" {
		t.Errorf("Expected password 'MySecurePassword123!', got '%s'", creds.Password)
	}
	
	if creds.URL != "http://localhost:8080" {
		t.Errorf("Expected URL 'http://localhost:8080', got '%s'", creds.URL)
	}
	
	if !creds.IsComplete() {
		t.Error("Credentials should be complete")
	}
	
	if !creds.HasPassword() {
		t.Error("Should have password")
	}
	
	if !creds.HasURL() {
		t.Error("Should have URL")
	}
	
	// Test incomplete log
	incompleteLog := "Starting container..."
	incompleteCreds := parser.ExtractCredentials(incompleteLog)
	
	if incompleteCreds.IsComplete() {
		t.Error("Incomplete credentials should not be marked as complete")
	}
	
	t.Logf("Log parser tests completed successfully")
}