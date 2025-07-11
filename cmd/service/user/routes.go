package user

import (
	"log"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /images", h.handleImages)
}

func (h *Handler) handleImages(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle Images - GET request")
}
