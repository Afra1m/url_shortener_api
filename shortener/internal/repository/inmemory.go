package repository

import (
	"shortener/internal/models"
	"sync"
	"time"
)

type URLRepository interface {
	Create(url models.URLStats) (string, error)
	GetByShortURL(shortURL string) (*models.URLStats, error)
	GetAll() ([]models.URLStats, error)
	IncrementAccess(shortURL string) error
}

type InMemoryURLRepository struct {
	urls map[string]models.URLStats
	mu   sync.RWMutex
}

func NewInMemoryURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		urls: make(map[string]models.URLStats),
	}
}

func (r *InMemoryURLRepository) Create(url models.URLStats) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[url.ShortURL] = url
	return url.ShortURL, nil
}

func (r *InMemoryURLRepository) GetByShortURL(shortURL string) (*models.URLStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if url, exists := r.urls[shortURL]; exists {
		return &url, nil
	}
	return nil, nil
}

func (r *InMemoryURLRepository) GetAll() ([]models.URLStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	urls := make([]models.URLStats, 0, len(r.urls))
	for _, url := range r.urls {
		urls = append(urls, url)
	}
	return urls, nil
}

func (r *InMemoryURLRepository) IncrementAccess(shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if url, exists := r.urls[shortURL]; exists {
		url.AccessCount++
		url.LastAccessed = time.Now()
		r.urls[shortURL] = url
	}
	return nil
}
