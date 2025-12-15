package models

import "time"

type AnalyticsEvent struct {
	ShortURL  string    `json:"short_url"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

type EventSummary struct {
	ShortURL     string `json:"short_url"`
	TotalClicks  int    `json:"total_clicks"`
	UniqueVisits int    `json:"unique_visits"`
	LastEvent    string `json:"last_event,omitempty"`
}
