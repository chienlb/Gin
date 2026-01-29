package apperror

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Status  int         `json:"-"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"` // Internal error for logging
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Common error codes
const (
	CodeValidationError     = "VALIDATION_ERROR"
	CodeNotFound            = "NOT_FOUND"
	CodeDuplicateEmail      = "DUPLICATE_EMAIL"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeBadRequest          = "BAD_REQUEST"
	CodeConflict            = "CONFLICT"
)

// New creates a new AppError
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewWithDetails creates an AppError with additional details
func NewWithDetails(code, message string, status int, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}

// NewWithError creates an AppError with an internal error
func NewWithError(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrValidation     = New(CodeValidationError, "Validation failed", http.StatusBadRequest)
	ErrNotFound       = New(CodeNotFound, "Resource not found", http.StatusNotFound)
	ErrUnauthorized   = New(CodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrConflict       = New(CodeConflict, "Resource already exists", http.StatusConflict)
	ErrInternalServer = New(CodeInternalServerError, "Internal server error", http.StatusInternalServerError)
)

// Validation errors
func ValidationError(field, reason string) *AppError {
	return NewWithDetails(
		CodeValidationError,
		"Validation failed",
		http.StatusBadRequest,
		map[string]string{
			"field":  field,
			"reason": reason,
		},
	)
}

// User-related errors
func DuplicateEmailError(email string) *AppError {
	return NewWithDetails(
		CodeDuplicateEmail,
		"User with this email already exists",
		http.StatusConflict,
		map[string]string{"email": email},
	)
}

func UserNotFoundError(id int) *AppError {
	return NewWithDetails(
		CodeNotFound,
		"User not found",
		http.StatusNotFound,
		map[string]interface{}{"user_id": id},
	)
}
