package repository

import (
	"analytics/internal/models"
	"sync"
	"time"
)

type AnalyticsRepository interface {
	SaveEvent(event models.AnalyticsEvent) error
	GetEventsByShortURL(shortURL string) ([]models.AnalyticsEvent, error)
	GetAllEvents() ([]models.AnalyticsEvent, error)
	GetSummary(shortURL string) (*models.EventSummary, error)
}

type InMemoryAnalyticsRepository struct {
	events []models.AnalyticsEvent
	mu     sync.RWMutex
}

func NewInMemoryAnalyticsRepository() *InMemoryAnalyticsRepository {
	return &InMemoryAnalyticsRepository{
		events: make([]models.AnalyticsEvent, 0),
	}
}

func (r *InMemoryAnalyticsRepository) SaveEvent(event models.AnalyticsEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)
	return nil
}

func (r *InMemoryAnalyticsRepository) GetEventsByShortURL(shortURL string) ([]models.AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.AnalyticsEvent
	for _, event := range r.events {
		if event.ShortURL == shortURL {
			result = append(result, event)
		}
	}
	return result, nil
}

func (r *InMemoryAnalyticsRepository) GetAllEvents() ([]models.AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.events, nil
}

func (r *InMemoryAnalyticsRepository) GetSummary(shortURL string) (*models.EventSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var summary models.EventSummary
	summary.ShortURL = shortURL

	uniqueIPs := make(map[string]bool)
	var lastEventTime time.Time

	for _, event := range r.events {
		if event.ShortURL == shortURL {
			summary.TotalClicks++
			uniqueIPs[event.IPAddress] = true

			if event.Timestamp.After(lastEventTime) {
				lastEventTime = event.Timestamp
			}
		}
	}

	summary.UniqueVisits = len(uniqueIPs)
	if !lastEventTime.IsZero() {
		summary.LastEvent = lastEventTime.Format(time.RFC3339)
	}

	return &summary, nil
}
