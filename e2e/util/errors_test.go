package util

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsNotFoundError_NilError(t *testing.T) {
	result := IsNotFoundError(nil)
	if result {
		t.Error("Expected IsNotFoundError(nil) to return false")
	}
}

func TestIsNotFoundError_NotFoundLowercase(t *testing.T) {
	err := errors.New("resource not found")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect 'not found' (lowercase)")
	}
}

func TestIsNotFoundError_NotFoundMixedCase(t *testing.T) {
	err := errors.New("Resource Not Found")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect 'Not Found' (mixed case)")
	}
}

func TestIsNotFoundError_NotFoundUppercase(t *testing.T) {
	err := errors.New("RESOURCE NOT FOUND")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect 'NOT FOUND' (uppercase)")
	}
}

func TestIsNotFoundError_404Code(t *testing.T) {
	err := errors.New("HTTP 404: resource does not exist")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect '404' error code")
	}
}

func TestIsNotFoundError_404Only(t *testing.T) {
	err := errors.New("404")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect '404' standalone")
	}
}

func TestIsNotFoundError_DoesNotExistLowercase(t *testing.T) {
	err := errors.New("resource does not exist")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect 'does not exist' (lowercase)")
	}
}

func TestIsNotFoundError_DoesNotExistMixedCase(t *testing.T) {
	err := errors.New("Resource Does Not Exist")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect 'Does Not Exist' (mixed case)")
	}
}

func TestIsNotFoundError_WrappedError(t *testing.T) {
	innerErr := errors.New("not found")
	wrappedErr := fmt.Errorf("failed to get resource: %w", innerErr)
	result := IsNotFoundError(wrappedErr)
	if !result {
		t.Error("Expected IsNotFoundError to detect 'not found' in wrapped error")
	}
}

func TestIsNotFoundError_OtherError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "permission denied",
			err:  errors.New("permission denied"),
			want: false,
		},
		{
			name: "internal server error",
			err:  errors.New("internal server error"),
			want: false,
		},
		{
			name: "connection timeout",
			err:  errors.New("connection timeout"),
			want: false,
		},
		{
			name: "invalid request",
			err:  errors.New("invalid request"),
			want: false,
		},
		{
			name: "empty error message",
			err:  errors.New(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFoundError(tt.err)
			if result != tt.want {
				t.Errorf("IsNotFoundError(%q) = %v, want %v", tt.err.Error(), result, tt.want)
			}
		})
	}
}

func TestIsNotFoundError_ComplexMessages(t *testing.T) {
	tests := []struct {
		name string
		err  string
		want bool
	}{
		{
			name: "API error with not found",
			err:  "API Error: The requested resource was not found",
			want: true,
		},
		{
			name: "404 with details",
			err:  "GET /api/v1/resource returned 404: resource not available",
			want: true,
		},
		{
			name: "does not exist in database",
			err:  "record does not exist in database",
			want: true,
		},
		{
			name: "cluster ID not found",
			err:  "cluster with ID 123 not found",
			want: true,
		},
		{
			name: "partial match should not trigger - founded",
			err:  "company founded in 2020",
			want: false,
		},
		{
			name: "partial match should not trigger - exists",
			err:  "resource exists and is running",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.err)
			result := IsNotFoundError(err)
			if result != tt.want {
				t.Errorf("IsNotFoundError(%q) = %v, want %v", tt.err, result, tt.want)
			}
		})
	}
}

func TestIsNotFoundError_MultipleMatches(t *testing.T) {
	err := errors.New("404 not found: resource does not exist")
	result := IsNotFoundError(err)
	if !result {
		t.Error("Expected IsNotFoundError to detect error with multiple match patterns")
	}
}
