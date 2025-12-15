package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AnalyticsClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAnalyticsClient(baseURL string) *AnalyticsClient {
	return &AnalyticsClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *AnalyticsClient) SendEvent(shortURL, eventType string) error {
	url := fmt.Sprintf("%s/events", c.baseURL)
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

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("analytics service returned status: %d", resp.StatusCode)
	}

	return nil
}
