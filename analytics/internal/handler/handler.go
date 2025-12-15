package handler

import (
	"analytics/internal/models"
	"analytics/internal/repository"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type AnalyticsHandler struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsHandler(repo repository.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

func (h *AnalyticsHandler) ReceiveEvent(w http.ResponseWriter, r *http.Request) {
	var eventRequest struct {
		ShortURL  string `json:"short_url"`
		EventType string `json:"event_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := models.AnalyticsEvent{
		ShortURL:  eventRequest.ShortURL,
		EventType: eventRequest.EventType,
		Timestamp: time.Now(),
		IPAddress: extractIP(r),
		UserAgent: r.UserAgent(),
	}

	err := h.repo.SaveEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "event recorded",
	})
}

func (h *AnalyticsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	events, err := h.repo.GetEventsByShortURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	summary, err := h.repo.GetSummary(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (h *AnalyticsHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.repo.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *AnalyticsHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "analytics",
	})
}

// метод для настройки маршрутов
func (h *AnalyticsHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/events", h.ReceiveEvent).Methods("POST")
	router.HandleFunc("/events/{shortURL}", h.GetEvents).Methods("GET")
	router.HandleFunc("/summary/{shortURL}", h.GetSummary).Methods("GET")
	router.HandleFunc("/events", h.GetAllEvents).Methods("GET")
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

func extractIP(r *http.Request) string {
	// Получаем IP из заголовков, если за прокси
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}
