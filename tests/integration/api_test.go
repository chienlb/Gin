package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-demo/internal/config"
	"gin-demo/internal/domain"
	"gin-demo/internal/handler"
	"gin-demo/internal/repository"
	"gin-demo/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load test configuration
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "gin_db_test",
			SSLMode:  "disable",
		},
	}

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Clean database
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", userHandler.GetAllUsers)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.POST("/users", userHandler.CreateUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	return router, db
}

func TestIntegration_CreateUser(t *testing.T) {
	router, _ := setupTestServer(t)

	// Prepare request
	user := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "Password123",
	}
	body, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "success" {
		t.Errorf("Expected status success, got %v", response["status"])
	}

	data := response["data"].(map[string]interface{})
	if data["email"] != "john@example.com" {
		t.Errorf("Expected email john@example.com, got %v", data["email"])
	}
}

func TestIntegration_GetUser(t *testing.T) {
	router, db := setupTestServer(t)

	// Create user first
	userRepo := repository.NewUserRepository(db)
	user := &domain.User{
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Password: "hashed_password",
	}
	userRepo.Create(user)

	// Execute request
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	if data["email"] != "jane@example.com" {
		t.Errorf("Expected email jane@example.com, got %v", data["email"])
	}
}

func TestIntegration_UpdateUser(t *testing.T) {
	router, db := setupTestServer(t)

	// Create user first
	userRepo := repository.NewUserRepository(db)
	user := &domain.User{
		Name:     "Original Name",
		Email:    "original@example.com",
		Password: "hashed_password",
	}
	userRepo.Create(user)

	// Update user
	updateData := map[string]interface{}{
		"name":  "Updated Name",
		"email": "updated@example.com",
	}
	body, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	if data["name"] != "Updated Name" {
		t.Errorf("Expected name Updated Name, got %v", data["name"])
	}
}

func TestIntegration_DeleteUser(t *testing.T) {
	router, db := setupTestServer(t)

	// Create user first
	userRepo := repository.NewUserRepository(db)
	user := &domain.User{
		Name:     "To Delete",
		Email:    "delete@example.com",
		Password: "hashed_password",
	}
	userRepo.Create(user)

	// Delete user
	req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify user is deleted
	req, _ = http.NewRequest("GET", "/api/v1/users/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestIntegration_GetAllUsers(t *testing.T) {
	router, db := setupTestServer(t)

	// Create multiple users
	userRepo := repository.NewUserRepository(db)
	users := []*domain.User{
		{Name: "User 1", Email: "user1@example.com", Password: "pass1"},
		{Name: "User 2", Email: "user2@example.com", Password: "pass2"},
		{Name: "User 3", Email: "user3@example.com", Password: "pass3"},
	}
	for _, u := range users {
		userRepo.Create(u)
	}

	// Execute request
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].([]interface{})
	if len(data) != 3 {
		t.Errorf("Expected 3 users, got %d", len(data))
	}
}

func TestIntegration_ValidationErrors(t *testing.T) {
	router, _ := setupTestServer(t)

	tests := []struct {
		name       string
		user       map[string]interface{}
		statusCode int
	}{
		{
			name:       "Invalid email",
			user:       map[string]interface{}{"name": "John", "email": "invalid", "password": "Pass123"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Weak password",
			user:       map[string]interface{}{"name": "John", "email": "john@example.com", "password": "weak"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Short name",
			user:       map[string]interface{}{"name": "J", "email": "john@example.com", "password": "Pass123"},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.user)
			req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
