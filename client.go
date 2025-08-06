package gamehub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SudoLab/gamehub-client/internal"
)

type Client struct {
	config     *ClientConfig
	httpClient *http.Client
}

func NewClient(config *ClientConfig) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func NewClientWithCredentials(baseURL, gameID, apiKey string) *Client {
	config := DefaultConfig().
		WithBaseURL(baseURL).
		WithGameID(gameID).
		WithAPIKey(apiKey)

	return NewClient(config)
}

// Auth methods

// GetUser retrieves user information using session token
func (c *Client) GetUser(ctx context.Context, sessionToken string) (*User, error) {
	req, err := c.newRequest(ctx, "GET", "/api/v1/auth/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+sessionToken)

	var user User
	if err := c.doRequest(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Coins methods

// GetUserCoins retrieves user's coin balance (internal API)
func (c *Client) GetUserCoins(ctx context.Context, userID int64) (int, error) {
	endpoint := fmt.Sprintf("/api/v1/internal/users/%d/coins", userID)
	req, err := c.newInternalRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return 0, err
	}

	var result struct {
		Balance int `json:"balance"`
	}
	if err := c.doRequest(req, &result); err != nil {
		return 0, err
	}

	return result.Balance, nil
}

// DeductCoins deducts coins from user's balance (internal API)
func (c *Client) DeductCoins(ctx context.Context, userID int64, amount int, reason, referenceID string) error {
	payload := map[string]any{
		"user_id":      userID,
		"amount":       amount,
		"game_id":      c.config.GameID,
		"reason":       reason,
		"reference_id": referenceID,
	}

	req, err := c.newInternalRequest(ctx, "POST", "/api/v1/internal/coins/deduct", payload)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// AddCoins adds coins to user's balance (internal API)
func (c *Client) AddCoins(ctx context.Context, userID int64, amount int, reason, referenceID string) error {
	payload := map[string]any{
		"user_id":      userID,
		"amount":       amount,
		"game_id":      c.config.GameID,
		"reason":       reason,
		"reference_id": referenceID,
	}

	req, err := c.newInternalRequest(ctx, "POST", "/api/v1/internal/coins/add", payload)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// Games methods

// ReportScore reports a game score (internal API)
func (c *Client) ReportScore(ctx context.Context, userID int64, score int64) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"game_id": c.config.GameID,
		"score":   score,
	}

	req, err := c.newInternalRequest(ctx, "POST", "/api/v1/internal/games/report-score", payload)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// GetAvailableGames retrieves all available games
func (c *Client) GetAvailableGames(ctx context.Context) ([]Game, error) {
	req, err := c.newRequest(ctx, "GET", "/api/v1/games", nil)
	if err != nil {
		return nil, err
	}

	var games []Game
	if err := c.doRequest(req, &games); err != nil {
		return nil, err
	}

	return games, nil
}

// Rankings methods

// GetGlobalRankings retrieves global rankings
func (c *Client) GetGlobalRankings(ctx context.Context, limit, offset int) ([]RankingEntry, error) {
	endpoint := fmt.Sprintf("/api/v1/rankings/global?limit=%d&offset=%d", limit, offset)
	req, err := c.newRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var rankings []RankingEntry
	if err := c.doRequest(req, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

// Helper methods

func (c *Client) newRequest(ctx context.Context, method, endpoint string, body any) (*http.Request, error) {
	url := c.config.BaseURL + endpoint

	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.config.UserAgent)

	return req, nil
}

func (c *Client) newInternalRequest(ctx context.Context, method, endpoint string, body any) (*http.Request, error) {
	req, err := c.newRequest(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.config.APIKey)
	req.Header.Set("X-Game-ID", c.config.GameID)

	return req, nil
}

func (c *Client) doRequest(req *http.Request, result any) error {
	resp, err := internal.DoWithRetry(c.httpClient, req, c.config.RetryCount)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return c.handleAPIError(resp.StatusCode, apiResp.Error)
	}

	if result != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}

	return nil
}

func (c *Client) handleAPIError(statusCode int, message string) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrUserNotFound
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusBadRequest:
		if message == "insufficient coins" {
			return ErrInsufficientCoins
		}
		return NewError(statusCode, message, "bad_request")
	case http.StatusInternalServerError:
		return ErrServerError
	default:
		return NewError(statusCode, message, "unknown_error")
	}
}
