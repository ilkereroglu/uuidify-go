package uuidify

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBaseURL     = "https://api.uuidify.io"
	defaultUserAgent   = "uuidify-go-sdk/1.0"
	defaultHTTPTimeout = 5 * time.Second
)

// Client is a minimal UUIDify API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// ClientOption configure a Client.
type ClientOption func(*Client)

// NewClient constructs a new Client, applying any provided options.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{Timeout: defaultHTTPTimeout},
		UserAgent:  defaultUserAgent,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(c)
	}

	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: defaultHTTPTimeout}
	}
	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}

	return c
}

// WithBaseURL overrides the default base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = url
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = client
	}
}

// WithUserAgent overrides the default User-Agent header value.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.UserAgent = ua
	}
}

// UUIDv1 fetches a UUID v1 value.
func (c *Client) UUIDv1(ctx context.Context) (string, error) {
	var resp UUIDResponse
	if err := c.doRequest(ctx, map[string]string{"version": "v1"}, &resp); err != nil {
		return "", err
	}
	return resp.UUID, nil
}

// UUIDv4 fetches a UUID v4 value.
func (c *Client) UUIDv4(ctx context.Context) (string, error) {
	var resp UUIDResponse
	if err := c.doRequest(ctx, map[string]string{"version": "v4"}, &resp); err != nil {
		return "", err
	}
	return resp.UUID, nil
}

// UUIDv7 fetches a UUID v7 value.
func (c *Client) UUIDv7(ctx context.Context) (string, error) {
	var resp UUIDResponse
	if err := c.doRequest(ctx, map[string]string{"version": "v7"}, &resp); err != nil {
		return "", err
	}
	return resp.UUID, nil
}

// ULID fetches a ULID value.
func (c *Client) ULID(ctx context.Context) (string, error) {
	var resp ULIDResponse
	if err := c.doRequest(ctx, map[string]string{"version": "ulid"}, &resp); err != nil {
		return "", err
	}
	return resp.ULID, nil
}

// UUIDBatch fetches multiple UUIDs of the given version.
func (c *Client) UUIDBatch(ctx context.Context, version string, count int) ([]string, error) {
	if !isSupportedUUIDVersion(version) {
		return nil, fmt.Errorf("version must be one of v1, v4, v7")
	}
	if count <= 0 || count > 1000 {
		return nil, fmt.Errorf("count must be between 1 and 1000")
	}

	query := map[string]string{
		"version": version,
		"count":   strconv.Itoa(count),
	}

	if count == 1 {
		var resp UUIDResponse
		if err := c.doRequest(ctx, query, &resp); err != nil {
			return nil, err
		}
		return []string{resp.UUID}, nil
	}

	var resp UUIDListResponse
	if err := c.doRequest(ctx, query, &resp); err != nil {
		return nil, err
	}
	return resp.UUIDs, nil
}

// ULIDBatch fetches multiple ULIDs.
func (c *Client) ULIDBatch(ctx context.Context, count int) ([]string, error) {
	if count <= 0 || count > 1000 {
		return nil, fmt.Errorf("count must be between 1 and 1000")
	}

	query := map[string]string{
		"version": "ulid",
		"count":   strconv.Itoa(count),
	}

	if count == 1 {
		var resp ULIDResponse
		if err := c.doRequest(ctx, query, &resp); err != nil {
			return nil, err
		}
		return []string{resp.ULID}, nil
	}

	var resp ULIDListResponse
	if err := c.doRequest(ctx, query, &resp); err != nil {
		return nil, err
	}
	return resp.ULIDs, nil
}

func isSupportedUUIDVersion(version string) bool {
	switch version {
	case "v1", "v4", "v7":
		return true
	default:
		return false
	}
}
