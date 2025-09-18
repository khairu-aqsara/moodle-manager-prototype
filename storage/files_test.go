package storage

import (
	"testing"
)

func TestFileManager(t *testing.T) {
	fm := NewFileManager()
	
	// Clean up any existing test files
	defer func() {
		fm.DeleteContainerID()
		fm.DeleteCredentials()
	}()
	
	// Test container ID operations
	testContainerID := "test-container-id-123"
	
	// Save container ID
	err := fm.SaveContainerID(testContainerID)
	if err != nil {
		t.Fatalf("Failed to save container ID: %v", err)
	}
	
	// Check if container ID file exists
	if !fm.ContainerIDExists() {
		t.Error("Container ID file should exist after saving")
	}
	
	// Load container ID
	loadedID, err := fm.LoadContainerID()
	if err != nil {
		t.Fatalf("Failed to load container ID: %v", err)
	}
	
	if loadedID != testContainerID {
		t.Errorf("Expected container ID %s, got %s", testContainerID, loadedID)
	}
	
	// Test credentials operations
	testPassword := "test-password-123"
	testURL := "http://test.localhost:8080"
	
	// Save credentials
	err = fm.SaveCredentials(testPassword, testURL)
	if err != nil {
		t.Fatalf("Failed to save credentials: %v", err)
	}
	
	// Check if credentials file exists
	if !fm.CredentialsExist() {
		t.Error("Credentials file should exist after saving")
	}
	
	// Load credentials
	creds, err := fm.LoadCredentials()
	if err != nil {
		t.Fatalf("Failed to load credentials: %v", err)
	}
	
	if creds["password"] != testPassword {
		t.Errorf("Expected password %s, got %s", testPassword, creds["password"])
	}
	
	if creds["url"] != testURL {
		t.Errorf("Expected URL %s, got %s", testURL, creds["url"])
	}
	
	t.Logf("File manager tests completed successfully")
}