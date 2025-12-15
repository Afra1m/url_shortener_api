package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"shortener/internal/handler"
	"shortener/internal/repository"
	"shortener/internal/service"
	"shortener/pkg/utils"

	"github.com/gorilla/mux"
)

// rootHandler обрабатывает корневой запрос
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "url-shortener",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string]string{
			"shorten":  "POST /shorten",
			"redirect": "GET /{shortURL}",
			"stats":    "GET /stats/{shortURL}",
			"all_urls": "GET /urls",
			"health":   "GET /health",
		},
	})
}

// healthHandler обрабатывает health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "url-shortener",
	})
}

func main() {
	// Конфигурация
	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", "8080")
	analyticsURL := getEnv("ANALYTICS_URL", "http://analytics:8081")

	// Инициализация компонентов
	repo := repository.NewInMemoryURLRepository()
	analyticsClient := service.NewHTTPAnalyticsClient(analyticsURL)
	baseURL := utils.GenerateBaseURL(host, port)
	shortenerService := service.NewShortenerService(repo, analyticsClient, baseURL)
	h := handler.NewHandler(shortenerService)

	// Настройка маршрутизатора
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// Добавляем основные маршруты
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// Настраиваем остальные маршруты через handler
	h.SetupRoutes(router)

	// Запуск сервера
	addr := fmt.Sprintf(":%s", port)
	log.Printf("URL Shortener service starting on %s", addr)
	log.Printf("Base URL: %s", baseURL)
	log.Printf("Analytics URL: %s", analyticsURL)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
