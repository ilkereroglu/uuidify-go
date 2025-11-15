package uuidify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func (c *Client) doRequest(ctx context.Context, query map[string]string, v interface{}) error {
	if c == nil {
		return &RequestError{Err: errors.New("client is nil")}
	}

	targetURL, err := c.buildURL(query)
	if err != nil {
		return &RequestError{Err: err}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return &RequestError{Err: err}
	}

	ua := c.UserAgent
	if ua == "" {
		ua = defaultUserAgent
	}
	req.Header.Set("User-Agent", ua)

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &RequestError{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := readBodySnippet(resp.Body)
		if message == "" {
			message = http.StatusText(resp.StatusCode)
		}
		return &APIError{StatusCode: resp.StatusCode, Message: message}
	}

	if v == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}
		return &DecodeError{Err: err}
	}

	return nil
}

func (c *Client) buildURL(query map[string]string) (string, error) {
	base := c.BaseURL
	if base == "" {
		base = defaultBaseURL
	}

	target, err := url.JoinPath(base, "/")
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(target)
	if err != nil {
		return "", err
	}

	if len(query) > 0 {
		values := parsed.Query()
		for key, value := range query {
			if key == "" {
				continue
			}
			values.Set(key, value)
		}
		parsed.RawQuery = values.Encode()
	}

	return parsed.String(), nil
}

func readBodySnippet(r io.Reader) string {
	if r == nil {
		return ""
	}

	limited := io.LimitReader(r, 4096)
	data, err := io.ReadAll(limited)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// APIError captures non-successful HTTP responses from the UUIDify API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return fmt.Sprintf("uuidify API error (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("uuidify API error (%d)", e.StatusCode)
}

// DecodeError wraps errors that occur while decoding API responses.
type DecodeError struct {
	Err error
}

func (e *DecodeError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("uuidify decode error: %v", e.Err)
}

func (e *DecodeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// RequestError wraps lower-level request construction or transport errors.
type RequestError struct {
	Err error
}

func (e *RequestError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("uuidify request error: %v", e.Err)
}

func (e *RequestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
