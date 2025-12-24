package manapool

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name      string
		err       *APIError
		wantError string
	}{
		{
			name: "with request ID",
			err: &APIError{
				StatusCode: 404,
				Message:    "not found",
				RequestID:  "req-12345",
			},
			wantError: "manapool API error (status 404, request req-12345): not found",
		},
		{
			name: "without request ID",
			err: &APIError{
				StatusCode: 500,
				Message:    "internal server error",
			},
			wantError: "manapool API error (status 500): internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantError {
				t.Errorf("APIError.Error() = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestAPIError_IsNotFound(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			want:       true,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			want:       false,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			got := err.IsNotFound()
			if got != tt.want {
				t.Errorf("APIError.IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsUnauthorized(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			want:       true,
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			want:       false,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			got := err.IsUnauthorized()
			if got != tt.want {
				t.Errorf("APIError.IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsForbidden(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			want:       true,
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			want:       false,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			got := err.IsForbidden()
			if got != tt.want {
				t.Errorf("APIError.IsForbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsRateLimited(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "429 Too Many Requests",
			statusCode: http.StatusTooManyRequests,
			want:       true,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			want:       false,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			got := err.IsRateLimited()
			if got != tt.want {
				t.Errorf("APIError.IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			want:       true,
		},
		{
			name:       "502 Bad Gateway",
			statusCode: http.StatusBadGateway,
			want:       true,
		},
		{
			name:       "503 Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
			want:       true,
		},
		{
			name:       "504 Gateway Timeout",
			statusCode: http.StatusGatewayTimeout,
			want:       true,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			want:       false,
		},
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			got := err.IsServerError()
			if got != tt.want {
				t.Errorf("APIError.IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "limit",
		Message: "must be non-negative",
	}

	want := "validation error for field 'limit': must be non-negative"
	got := err.Error()

	if got != want {
		t.Errorf("ValidationError.Error() = %q, want %q", got, want)
	}
}

func TestNetworkError_Error(t *testing.T) {
	tests := []struct {
		name      string
		err       *NetworkError
		wantError string
	}{
		{
			name: "with underlying error",
			err: &NetworkError{
				Message: "connection failed",
				Err:     errors.New("dial timeout"),
			},
			wantError: "network error: connection failed: dial timeout",
		},
		{
			name: "without underlying error",
			err: &NetworkError{
				Message: "rate limiter error",
			},
			wantError: "network error: rate limiter error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantError {
				t.Errorf("NetworkError.Error() = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestNetworkError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	netErr := &NetworkError{
		Message: "network error",
		Err:     underlyingErr,
	}

	got := netErr.Unwrap()
	if got != underlyingErr {
		t.Errorf("NetworkError.Unwrap() = %v, want %v", got, underlyingErr)
	}

	// Test nil case
	netErr2 := &NetworkError{
		Message: "network error",
	}

	got2 := netErr2.Unwrap()
	if got2 != nil {
		t.Errorf("NetworkError.Unwrap() = %v, want nil", got2)
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(404, "not found")

	if err.StatusCode != 404 {
		t.Errorf("NewAPIError().StatusCode = %d, want 404", err.StatusCode)
	}
	if err.Message != "not found" {
		t.Errorf("NewAPIError().Message = %q, want %q", err.Message, "not found")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("field", "invalid")

	if err.Field != "field" {
		t.Errorf("NewValidationError().Field = %q, want %q", err.Field, "field")
	}
	if err.Message != "invalid" {
		t.Errorf("NewValidationError().Message = %q, want %q", err.Message, "invalid")
	}
}

func TestNewNetworkError(t *testing.T) {
	underlyingErr := errors.New("underlying")
	err := NewNetworkError("network failed", underlyingErr)

	if err.Message != "network failed" {
		t.Errorf("NewNetworkError().Message = %q, want %q", err.Message, "network failed")
	}
	if err.Err != underlyingErr {
		t.Errorf("NewNetworkError().Err = %v, want %v", err.Err, underlyingErr)
	}
}

func TestAPIError_Response(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
	}

	err := &APIError{
		StatusCode: 404,
		Message:    "not found",
		Response:   resp,
	}

	if err.Response != resp {
		t.Errorf("APIError.Response = %v, want %v", err.Response, resp)
	}
}
