package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"shortener/internal/models"
	"shortener/internal/repository"
	"strings"
	"time"
)

type ShortenerService struct {
	repo      repository.URLRepository
	analytics AnalyticsClient
	baseURL   string
}

type AnalyticsClient interface {
	SendEvent(shortURL, eventType string) error
}

func NewShortenerService(repo repository.URLRepository, analytics AnalyticsClient, baseURL string) *ShortenerService {
	return &ShortenerService{
		repo:      repo,
		analytics: analytics,
		baseURL:   baseURL,
	}
}

func (s *ShortenerService) ShortenURL(request models.ShortURLRequest) (*models.ShortURLResponse, error) {
	shortURL := request.CustomAlias
	if shortURL == "" {
		shortURL = generateShortURL(request.OriginalURL)
	}

	urlStats := models.URLStats{
		ShortURL:    shortURL,
		OriginalURL: request.OriginalURL,
		CreatedAt:   time.Now(),
		AccessCount: 0,
	}

	_, err := s.repo.Create(urlStats)
	if err != nil {
		return nil, err
	}

	go s.analytics.SendEvent(shortURL, "created")

	response := &models.ShortURLResponse{
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, shortURL),
		OriginalURL: request.OriginalURL,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	return response, nil
}

func (s *ShortenerService) Redirect(shortURL string) (string, error) {
	urlStats, err := s.repo.GetByShortURL(shortURL)
	if err != nil {
		return "", err
	}

	if urlStats == nil {
		return "", fmt.Errorf("URL not found")
	}

	err = s.repo.IncrementAccess(shortURL)
	if err != nil {
		return "", err
	}

	go s.analytics.SendEvent(shortURL, "accessed")

	return urlStats.OriginalURL, nil
}

func (s *ShortenerService) GetStats(shortURL string) (*models.URLStats, error) {
	return s.repo.GetByShortURL(shortURL)
}

func (s *ShortenerService) GetAllURLs() ([]models.URLStats, error) {
	return s.repo.GetAll()
}

func generateShortURL(originalURL string) string {
	hash := md5.Sum([]byte(originalURL + time.Now().String()))
	shortURL := hex.EncodeToString(hash[:])[:8]
	return shortURL
}

// HTTPAnalyticsClient реализация для отправки событий в analytics сервис
type HTTPAnalyticsClient struct {
	analyticsURL string
	httpClient   *http.Client
}

func NewHTTPAnalyticsClient(analyticsURL string) *HTTPAnalyticsClient {
	return &HTTPAnalyticsClient{
		analyticsURL: analyticsURL,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPAnalyticsClient) SendEvent(shortURL, eventType string) error {
	url := fmt.Sprintf("%s/events", c.analyticsURL)
	body := fmt.Sprintf(`{"short_url":"%s","event_type":"%s"}`, shortURL, eventType)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
