package handler

import (
	"encoding/json"
	"net/http"
	"shortener/internal/models"
	"shortener/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.ShortenerService
}

func NewHandler(service *service.ShortenerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var request models.ShortURLRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.OriginalURL == "" {
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.ShortenURL(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, err := h.service.Redirect(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	stats, err := h.service.GetStats(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if stats == nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) GetAllURLs(w http.ResponseWriter, r *http.Request) {
	urls, err := h.service.GetAllURLs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func (h *Handler) SetupRoutes(router *mux.Router) {

	// API endpoints
	router.HandleFunc("/shorten", h.ShortenURL).Methods("POST")
	router.HandleFunc("/{shortURL}", h.Redirect).Methods("GET")
	router.HandleFunc("/stats/{shortURL}", h.GetStats).Methods("GET")
	router.HandleFunc("/urls", h.GetAllURLs).Methods("GET")
	// Health endpoint удален отсюда
}
