// TODO: Add the following to a command_execution.go file:
// executeCommand function
// CommandRequest struct

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests the /execute-command endpoint with a valid command and ensures that the output is correct
func TestExecuteCommandValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(executeCommand))
	defer server.Close()

	// Mock a valid command
	cmd := "echo Hello, World!"
	reqBody, _ := json.Marshal(CommandRequest{
		Command:    cmd,
		ListenerID: "listener_123",
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

	// Check output
	var output string
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	if output != "Hello, World!\n" {
		t.Errorf("Expected output 'Hello, World!', got '%s'", output)
	}
}

// This test is to ensure that the given payload is a valid one
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// This test checks whether the command itself is missing or not
func TestExecuteCommandMissingCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(executeCommand))
	defer server.Close()

	reqBody, _ := json.Marshal(CommandRequest{
		ListenerID: "listener_123",
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

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
