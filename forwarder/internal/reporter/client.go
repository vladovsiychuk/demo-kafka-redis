package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	Fetch(ctx context.Context, partner string, date string) (*Data, error)
}

type httpClient struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *httpClient) Fetch(ctx context.Context, partner, date string) (*Data, error) {
	url := fmt.Sprintf("%s/%s?date=%s", c.baseURL, partner, date)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
