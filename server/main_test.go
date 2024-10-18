package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Initialize the API key for testing
	apiKey = generateAPIKey()

	// Run the tests
	exitCode := m.Run()

	// Clean up and exit
	os.Exit(exitCode)
}

// TestCreateUser tests the user creation endpoint.
func TestCreateUser(t *testing.T) {
	// Set up a new HTTP server
	server := httptest.NewServer(http.HandlerFunc(createUser))
	defer server.Close()

	// Test user creation
	user := User{Username: "testuser", Password: "testpass"}
	userJSON, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Check if user is added
	usersMtx.Lock()
	defer usersMtx.Unlock()
	if len(users) != 1 || users[0].Username != "testuser" {
		t.Error("User not added correctly")
	}
}

// TestLogin tests the login functionality.
func TestLogin(t *testing.T) {
	// Set up a new HTTP server
	server := httptest.NewServer(http.HandlerFunc(login))
	defer server.Close()

	// First, register a user
	user := User{Username: "testuser", Password: "testpass"}
	usersMtx.Lock()
	users = append(users, user)
	usersMtx.Unlock()

	// Test login
	loginUser := User{Username: "testuser", Password: "testpass"}
	loginJSON, _ := json.Marshal(loginUser)

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestUnauthorizedLogin tests login with invalid credentials.
func TestUnauthorizedLogin(t *testing.T) {
	// Set up a new HTTP server
	server := httptest.NewServer(http.HandlerFunc(login))
	defer server.Close()

	// Test login with non-existing user
	loginUser := User{Username: "invaliduser", Password: "wrongpass"}
	loginJSON, _ := json.Marshal(loginUser)

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

// This test checks for a valid response with an empty payload
func TestExecuteCommandInvalidPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(executeCommand))
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// This test checks for a valid response when the command is missing
func TestExecuteCommandMissingCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(executeCommand))
	defer server.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"ListenerID": "listener_123",
	})

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
