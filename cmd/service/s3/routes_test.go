package s3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type mockService struct {
	shouldError bool
	errorMsg    string
	returnURL   string
}

func (m *mockService) GenerateSignedURL(key string, expiration time.Duration) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("%s", m.errorMsg)
	}
	return m.returnURL, nil
}

func TestHandler_handleGetSignedURL(t *testing.T) {
	tests := []struct {
		name               string
		handler            *Handler
		query              string
		expectedStatus     int
		expectedURL        string
		shouldContainError bool
	}{
		{
			name: "successful request",
			handler: &Handler{
				service: &mockService{
					shouldError: false,
					returnURL:   "https://test-bucket.s3.amazonaws.com/resume/Resume.pdf?signed-url",
				},
			},
			query:          "",
			expectedStatus: http.StatusOK,
			expectedURL:    "https://test-bucket.s3.amazonaws.com/resume/Resume.pdf?signed-url",
		},
		{
			name: "successful request with custom expiration",
			handler: &Handler{
				service: &mockService{
					shouldError: false,
					returnURL:   "https://test-bucket.s3.amazonaws.com/resume/Resume.pdf?signed-url",
				},
			},
			query:          "?expires_in=30m",
			expectedStatus: http.StatusOK,
			expectedURL:    "https://test-bucket.s3.amazonaws.com/resume/Resume.pdf?signed-url",
		},
		{
			name: "service unavailable",
			handler: &Handler{
				service: nil,
			},
			query:              "",
			expectedStatus:     http.StatusServiceUnavailable,
			shouldContainError: true,
		},
		{
			name: "service error",
			handler: &Handler{
				service: &mockService{
					shouldError: true,
					errorMsg:    "AWS error",
				},
			},
			query:              "",
			expectedStatus:     http.StatusInternalServerError,
			shouldContainError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/resume"+tt.query, nil)
			w := httptest.NewRecorder()

			tt.handler.handleGetSignedURL(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type to be application/json")
			}

			if tt.shouldContainError {
				var errorResp ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Errorf("failed to decode error response: %v", err)
				}
				if errorResp.Error == "" {
					t.Errorf("expected error message in response")
				}
			} else {
				var response SignedURLResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.URL != tt.expectedURL {
					t.Errorf("expected URL %q but got %q", tt.expectedURL, response.URL)
				}
				if response.ExpiresAt == "" {
					t.Errorf("expected ExpiresAt to be set")
				}
			}
		})
	}
}

func TestHandler_handleGetSignedURL_ExpirationParsing(t *testing.T) {
	handler := &Handler{
		service: &mockService{
			shouldError: false,
			returnURL:   "https://test-bucket.s3.amazonaws.com/resume/Resume.pdf?signed-url",
		},
	}

	tests := []struct {
		name           string
		expiresIn      string
		expectedStatus int
	}{
		{
			name:           "valid duration - minutes",
			expiresIn:      "30m",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid duration - hours",
			expiresIn:      "2h",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid duration - seconds",
			expiresIn:      "300s",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid duration - too long",
			expiresIn:      "25h",
			expectedStatus: http.StatusOK, // Should fall back to default
		},
		{
			name:           "invalid duration - negative",
			expiresIn:      "-1h",
			expectedStatus: http.StatusOK, // Should fall back to default
		},
		{
			name:           "invalid duration - malformed",
			expiresIn:      "invalid",
			expectedStatus: http.StatusOK, // Should fall back to default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/resume?expires_in="+tt.expiresIn, nil)
			w := httptest.NewRecorder()

			handler.handleGetSignedURL(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response SignedURLResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.URL == "" {
					t.Errorf("expected URL to be set")
				}
			}
		})
	}
}

func TestHandler_writeError(t *testing.T) {
	handler := &Handler{}

	tests := []struct {
		name       string
		message    string
		statusCode int
	}{
		{
			name:       "bad request error",
			message:    "invalid request",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "internal server error",
			message:    "something went wrong",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "service unavailable",
			message:    "service down",
			statusCode: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handler.writeError(w, tt.message, tt.statusCode)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d but got %d", tt.statusCode, w.Code)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type to be application/json")
			}

			var errorResp ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
				t.Errorf("failed to decode error response: %v", err)
			}

			if errorResp.Error != tt.message {
				t.Errorf("expected error message %q but got %q", tt.message, errorResp.Error)
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	originalBucket := os.Getenv("AWS_S3_BUCKET")
	defer os.Setenv("AWS_S3_BUCKET", originalBucket)

	tests := []struct {
		name             string
		bucketEnv        string
		expectNilService bool
	}{
		{
			name:             "valid environment",
			bucketEnv:        "test-bucket",
			expectNilService: false,
		},
		{
			name:             "missing bucket environment",
			bucketEnv:        "",
			expectNilService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.bucketEnv == "" {
				os.Unsetenv("AWS_S3_BUCKET")
			} else {
				os.Setenv("AWS_S3_BUCKET", tt.bucketEnv)
			}

			handler := NewHandler()

			if handler == nil {
				t.Errorf("expected handler to be non-nil")
			}

			if tt.expectNilService {
				if handler.service != nil {
					t.Errorf("expected service to be nil when bucket env is missing")
				}
			} else {
				if handler.service == nil {
					t.Errorf("expected service to be non-nil when bucket env is set")
				}
			}
		})
	}
}

func TestHandler_RegisterRoutes(t *testing.T) {
	handler := NewHandler()
	router := http.NewServeMux()

	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/resume", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Errorf("route not registered properly")
	}
}
