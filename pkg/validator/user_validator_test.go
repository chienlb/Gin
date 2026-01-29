package validator

import (
	"testing"
)

func TestValidateName_Success(t *testing.T) {
	validator := NewUserValidator()

	tests := []string{
		"John Doe",
		"Mary Jane",
		"José María",
		"李明",
		"A B",
	}

	for _, name := range tests {
		err := validator.ValidateName(name)
		if err != nil {
			t.Errorf("Expected valid name %s, got error: %v", name, err)
		}
	}
}

func TestValidateName_Invalid(t *testing.T) {
	validator := NewUserValidator()

	tests := []struct {
		name string
		want string
	}{
		{"", "Name is required"},
		{"X", "Name must be between 2 and 100 characters"},
		{"a", "Name must be between 2 and 100 characters"},
		{"John@Doe", "Name contains invalid characters"},
		{"User123", "Name contains invalid characters"},
	}

	for _, tt := range tests {
		err := validator.ValidateName(tt.name)
		if err == nil {
			t.Errorf("Expected error for name %s", tt.name)
		}
	}
}

func TestValidateEmail_Success(t *testing.T) {
	validator := NewUserValidator()

	tests := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.co.uk",
		"test123@test-domain.com",
	}

	for _, email := range tests {
		err := validator.ValidateEmail(email)
		if err != nil {
			t.Errorf("Expected valid email %s, got error: %v", email, err)
		}
	}
}

func TestValidateEmail_Invalid(t *testing.T) {
	validator := NewUserValidator()

	tests := []string{
		"",
		"invalid",
		"@example.com",
		"user@",
		"user @example.com",
		"user@.com",
	}

	for _, email := range tests {
		err := validator.ValidateEmail(email)
		if err == nil {
			t.Errorf("Expected error for email %s", email)
		}
	}
}

func TestValidatePassword_Success(t *testing.T) {
	validator := NewUserValidator()

	tests := []string{
		"Password123",
		"MyPass1",
		"Secure123Pass",
		"Aa1bcd",
	}

	for _, password := range tests {
		err := validator.ValidatePassword(password)
		if err != nil {
			t.Errorf("Expected valid password, got error: %v", err)
		}
	}
}

func TestValidatePassword_Invalid(t *testing.T) {
	validator := NewUserValidator()

	tests := []struct {
		password string
		want     string
	}{
		{"", "Password is required"},
		{"12345", "Password must be at least 6 characters"},
		{"password", "Password must contain at least one uppercase letter"},
		{"PASSWORD", "Password must contain at least one lowercase letter"},
		{"Password", "Password must contain at least one digit"},
	}

	for _, tt := range tests {
		err := validator.ValidatePassword(tt.password)
		if err == nil {
			t.Errorf("Expected error for password %s", tt.password)
		}
	}
}

func TestValidateCreateRequest(t *testing.T) {
	validator := NewUserValidator()

	err := validator.ValidateCreateRequest("John Doe", "john@example.com", "Password123")
	if err != nil {
		t.Errorf("Expected valid create request, got error: %v", err)
	}

	err = validator.ValidateCreateRequest("X", "invalid-email", "weak")
	if err == nil {
		t.Error("Expected validation error for invalid create request")
	}
}

func TestValidateUpdateRequest(t *testing.T) {
	validator := NewUserValidator()

	err := validator.ValidateUpdateRequest("John Doe", "john@example.com")
	if err != nil {
		t.Errorf("Expected valid update request, got error: %v", err)
	}

	err = validator.ValidateUpdateRequest("X", "invalid-email")
	if err == nil {
		t.Error("Expected validation error for invalid update request")
	}
}
