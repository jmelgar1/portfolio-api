package s3

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/url"
	"time"
)

type Handler struct {
	service S3Service
}

func NewHandler() *Handler {
	service, err := NewService()
	if err != nil {
		log.Printf("Failed to initialize S3 service: %v", err)
		return &Handler{service: nil}
	}

	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /resume", h.handleGetSignedURL)
}

type SignedURLResponse struct {
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) handleGetSignedURL(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		h.writeError(w, "S3 service not available", http.StatusServiceUnavailable)
		return
	}

	key := "resume/Resume.pdf"

	expiration := 15 * time.Minute
	if expStr := r.URL.Query().Get("expires_in"); expStr != "" {
		if parsedExp, err := time.ParseDuration(expStr); err == nil && parsedExp > 0 && parsedExp <= 24*time.Hour {
			expiration = parsedExp
		}
	}

	signedURL, err := h.service.GenerateSignedURL(key, expiration)
	if err != nil {
		log.Printf("Failed to generate signed URL: %v", err)
		h.writeError(w, "failed to generate signed URL", http.StatusInternalServerError)
		return
	}

	response := SignedURLResponse{
		URL:       signedURL,
		ExpiresAt: time.Now().Add(expiration).Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
