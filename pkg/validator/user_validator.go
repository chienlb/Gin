package validator

import (
	"regexp"
	"strings"
	"unicode"

	"gin-demo/pkg/apperror"
)

// UserValidator validates user input
type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateCreateRequest validates user creation request
func (v *UserValidator) ValidateCreateRequest(name, email, password string) *apperror.AppError {
	if err := v.ValidateName(name); err != nil {
		return err
	}
	if err := v.ValidateEmail(email); err != nil {
		return err
	}
	if err := v.ValidatePassword(password); err != nil {
		return err
	}
	return nil
}

// ValidateName validates user name
func (v *UserValidator) ValidateName(name string) *apperror.AppError {
	name = strings.TrimSpace(name)

	if name == "" {
		return apperror.ValidationError("name", "name is required")
	}
	if len(name) < 2 {
		return apperror.ValidationError("name", "name must be at least 2 characters")
	}
	if len(name) > 100 {
		return apperror.ValidationError("name", "name must not exceed 100 characters")
	}

	// Check if name contains only valid characters
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' && r != '\'' {
			return apperror.ValidationError("name", "name contains invalid characters")
		}
	}

	return nil
}

// ValidateEmail validates email format
func (v *UserValidator) ValidateEmail(email string) *apperror.AppError {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if email == "" {
		return apperror.ValidationError("email", "email is required")
	}

	if len(email) > 100 {
		return apperror.ValidationError("email", "email is too long")
	}

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return apperror.ValidationError("email", "invalid email format")
	}

	return nil
}

// ValidatePassword validates password strength
func (v *UserValidator) ValidatePassword(password string) *apperror.AppError {
	if password == "" {
		return apperror.ValidationError("password", "password is required")
	}

	if len(password) < 6 {
		return apperror.ValidationError("password", "password must be at least 6 characters")
	}

	if len(password) > 128 {
		return apperror.ValidationError("password", "password must not exceed 128 characters")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	// Check for at least one lowercase letter
	hasLower := false
	// Check for at least one digit
	hasDigit := false

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return apperror.ValidationError("password",
			"password must contain uppercase, lowercase, and digit")
	}

	return nil
}

// ValidateUpdateRequest validates user update request
func (v *UserValidator) ValidateUpdateRequest(name, email string) *apperror.AppError {
	if name != "" {
		if err := v.ValidateName(name); err != nil {
			return err
		}
	}

	if email != "" {
		if err := v.ValidateEmail(email); err != nil {
			return err
		}
	}

	return nil
}
