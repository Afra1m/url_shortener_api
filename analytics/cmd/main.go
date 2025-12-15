package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"analytics/internal/handler"
	"analytics/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	// Конфигурация
	port := getEnv("PORT", "8081")

	// Инициализация компонентов
	repo := repository.NewInMemoryAnalyticsRepository()
	analyticsHandler := handler.NewAnalyticsHandler(repo)

	// Настройка маршрутизатора
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	analyticsHandler.SetupRoutes(router) // Используем метод SetupRoutes

	// Запуск сервера
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Analytics service starting on %s", addr)

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
