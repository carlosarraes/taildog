package config

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/carlosarraes/taildog/internal/client"
)

type Config struct {
	APIKey string
	AppKey string
	Site   string

	Timeout    time.Duration
	MaxRetries int

	Query string
}

func NewConfig(apiKey, appKey, site, query string) (*Config, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, errors.New("apiKey is empty or contains only whitespace")
	}

	appKey = strings.TrimSpace(appKey)
	if appKey == "" {
		return nil, errors.New("appKey is empty or contains only whitespace")
	}

	site = strings.TrimSpace(site)
	if site == "" {
		site = "datadoghq.com"
	}

	return &Config{
		APIKey:     apiKey,
		AppKey:     appKey,
		Site:       site,
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Query:      query,
	}, nil
}

func (c *Config) ToClientConfig() *client.Config {
	return &client.Config{
		APIKey:     c.APIKey,
		AppKey:     c.AppKey,
		Site:       c.Site,
		Timeout:    c.Timeout,
		MaxRetries: c.MaxRetries,
	}
}

func TestAuthentication(ctx context.Context, ddClient *client.Client) error {
	_, err := ddClient.FetchLogs(ctx, "*")
	if err != nil {
		return fmt.Errorf("authentication test failed: %w", err)
	}
	return nil
}
