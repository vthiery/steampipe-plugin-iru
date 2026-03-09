package iru

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)


// Client wraps the Iru (Kandji) REST API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// newClient creates a new Client using the given base URL and API token.
func newClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getClient retrieves a configured Client from the plugin connection config.
func getClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	if cfg.URL == nil || *cfg.URL == "" {
		return nil, fmt.Errorf("url must be configured for the iru plugin")
	}
	if cfg.Token == nil || *cfg.Token == "" {
		return nil, fmt.Errorf("token must be configured for the iru plugin")
	}
	return newClient(*cfg.URL, *cfg.Token), nil
}

// get performs an authenticated GET request and unmarshals the response body into result.
func (c *Client) get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("iru API error: status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("unmarshalling response: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Shared sub-types used across multiple tables
// ---------------------------------------------------------------------------

// parseNextCursor extracts the "cursor" query parameter from a paginated "next" URL.
func parseNextCursor(nextURL string) (string, error) {
	u, err := url.Parse(nextURL)
	if err != nil {
		return "", err
	}
	return u.Query().Get("cursor"), nil
}

// AssignedUser represents the user assigned to a device in Iru.
type AssignedUser struct {
	Email    string `json:"email"`
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
	Name     string `json:"name"`
}
