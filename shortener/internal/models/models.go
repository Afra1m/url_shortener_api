package models

import "time"

type ShortURLRequest struct {
	OriginalURL string `json:"original_url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

type ShortURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

type URLStats struct {
	ShortURL     string    `json:"short_url"`
	OriginalURL  string    `json:"original_url"`
	CreatedAt    time.Time `json:"created_at"`
	AccessCount  int       `json:"access_count"`
	LastAccessed time.Time `json:"last_accessed,omitempty"`
}
