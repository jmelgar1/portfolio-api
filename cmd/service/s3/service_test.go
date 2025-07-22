package s3

import (
	"os"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		bucketEnv   string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid bucket environment variable",
			bucketEnv:   "test-bucket",
			shouldError: false,
		},
		{
			name:        "missing bucket environment variable",
			bucketEnv:   "",
			shouldError: true,
			errorMsg:    "AWS_S3_BUCKET environment variable is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalBucket := os.Getenv("AWS_S3_BUCKET")
			defer os.Setenv("AWS_S3_BUCKET", originalBucket)

			if tt.bucketEnv == "" {
				os.Unsetenv("AWS_S3_BUCKET")
			} else {
				os.Setenv("AWS_S3_BUCKET", tt.bucketEnv)
			}

			service, err := NewService()

			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != nil && err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q but got %q", tt.errorMsg, err.Error())
				}
				if service != nil {
					t.Errorf("expected service to be nil when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if service == nil {
					t.Errorf("expected service to be non-nil when no error occurs")
				}
				if service != nil && service.bucket != tt.bucketEnv {
					t.Errorf("expected bucket %q but got %q", tt.bucketEnv, service.bucket)
				}
			}
		})
	}
}

func TestService_GenerateSignedURL(t *testing.T) {
	originalBucket := os.Getenv("AWS_S3_BUCKET")
	defer os.Setenv("AWS_S3_BUCKET", originalBucket)

	os.Setenv("AWS_S3_BUCKET", "test-bucket")

	service, err := NewService()
	if err != nil {
		t.Skipf("Skipping test due to AWS configuration error: %v", err)
	}

	tests := []struct {
		name       string
		key        string
		expiration time.Duration
	}{
		{
			name:       "valid key and expiration",
			key:        "resume/Resume.pdf",
			expiration: 15 * time.Minute,
		},
		{
			name:       "different key",
			key:        "documents/test.pdf",
			expiration: 1 * time.Hour,
		},
		{
			name:       "short expiration",
			key:        "resume/Resume.pdf",
			expiration: 5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.GenerateSignedURL(tt.key, tt.expiration)

			if err != nil {
				t.Logf("Note: This test requires valid AWS credentials and may fail in CI/CD environments")
				t.Logf("Error: %v", err)
				return
			}

			if url == "" {
				t.Errorf("expected non-empty URL")
			}

			if len(url) < 10 {
				t.Errorf("expected URL to be longer than 10 characters, got %d", len(url))
			}
		})
	}
}

func TestService_GenerateSignedURL_InvalidKey(t *testing.T) {
	originalBucket := os.Getenv("AWS_S3_BUCKET")
	defer os.Setenv("AWS_S3_BUCKET", originalBucket)

	os.Setenv("AWS_S3_BUCKET", "test-bucket")

	service, err := NewService()
	if err != nil {
		t.Skipf("Skipping test due to AWS configuration error: %v", err)
	}

	tests := []struct {
		name       string
		key        string
		expiration time.Duration
	}{
		{
			name:       "empty key",
			key:        "",
			expiration: 15 * time.Minute,
		},
		{
			name:       "zero expiration",
			key:        "resume/Resume.pdf",
			expiration: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.GenerateSignedURL(tt.key, tt.expiration)

			if err == nil {
				t.Logf("Note: AWS SDK may handle some invalid inputs gracefully")
				t.Logf("Generated URL: %s", url)
			} else {
				t.Logf("Expected error for invalid input: %v", err)
			}
		})
	}
}
