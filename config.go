package gamehub

import (
	"time"
)

type ClientConfig struct {
	BaseURL    string
	GameID     string
	APIKey     string
	Timeout    time.Duration
	RetryCount int
	UserAgent  string
}

func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		UserAgent:  "GameHub-Client/1.0",
	}
}

func (c *ClientConfig) WithBaseURL(url string) *ClientConfig {
	c.BaseURL = url
	return c
}

func (c *ClientConfig) WithGameID(gameID string) *ClientConfig {
	c.GameID = gameID
	return c
}

func (c *ClientConfig) WithAPIKey(apiKey string) *ClientConfig {
	c.APIKey = apiKey
	return c
}

func (c *ClientConfig) WithTimeout(timeout time.Duration) *ClientConfig {
	c.Timeout = timeout
	return c
}

func (c *ClientConfig) WithRetryCount(count int) *ClientConfig {
	c.RetryCount = count
	return c
}
