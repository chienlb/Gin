package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var baseURL string

func TestMain(m *testing.M) {
	// Get API URL from environment or use default
	baseURL = os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Wait for API to be ready
	waitForAPI()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func waitForAPI() {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("API is ready")
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Warning: API not responding, tests may fail")
}

func TestE2E_HealthCheck(t *testing.T) {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to reach health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestE2E_UserCRUD(t *testing.T) {
	// Create user
	createData := map[string]interface{}{
		"name":     "E2E Test User",
		"email":    fmt.Sprintf("e2e-%d@example.com", time.Now().Unix()),
		"password": "TestPass123",
	}
	body, _ := json.Marshal(createData)

	resp, err := http.Post(baseURL+"/api/v1/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createResponse)
	data := createResponse["data"].(map[string]interface{})
	userID := int(data["id"].(float64))

	// Get user
	resp, err = http.Get(fmt.Sprintf("%s/api/v1/users/%d", baseURL, userID))
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Update user
	updateData := map[string]interface{}{
		"name":  "E2E Updated User",
		"email": createData["email"],
	}
	body, _ = json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/users/%d", baseURL, userID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Delete user
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/users/%d", baseURL, userID), nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify deletion
	resp, err = http.Get(fmt.Sprintf("%s/api/v1/users/%d", baseURL, userID))
	if err != nil {
		t.Fatalf("Failed to verify deletion: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestE2E_GetAllUsers(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/v1/users")
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if response["status"] != "success" {
		t.Errorf("Expected success status")
	}
}

func TestE2E_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]interface{}
		statusCode int
	}{
		{
			name:       "Invalid email",
			data:       map[string]interface{}{"name": "Test", "email": "invalid", "password": "Pass123"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Missing fields",
			data:       map[string]interface{}{"name": "Test"},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.data)
			resp, err := http.Post(baseURL+"/api/v1/users", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}
