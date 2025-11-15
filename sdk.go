package uuidify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultBaseURL     = "https://api.uuidify.io"
	defaultUserAgent   = "uuidify-go-sdk/1.0"
	defaultHTTPTimeout = 5 * time.Second
)

// NewDefaultClient creates a client preconfigured with the public API endpoint.
func NewDefaultClient(opts ...ClientOption) (*Client, error) {
	baseOpts := []ClientOption{
		WithHTTPClient(&http.Client{Timeout: defaultHTTPTimeout}),
		WithUserAgent(defaultUserAgent),
	}

	client, err := NewClient(DefaultBaseURL, append(baseOpts, opts...)...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// WithUserAgent ensures every request carries the provided User-Agent header.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) error {
		if ua == "" {
			return nil
		}
		c.RequestEditors = append(c.RequestEditors, func(ctx context.Context, req *http.Request) error {
			req.Header.Set("User-Agent", ua)
			return nil
		})
		return nil
	}
}

// UUIDv1 fetches a UUID v1 value.
func (c *Client) UUIDv1(ctx context.Context) (string, error) {
	return c.singleUUID(ctx, GetParamsVersionV1)
}

// UUIDv4 fetches a UUID v4 value.
func (c *Client) UUIDv4(ctx context.Context) (string, error) {
	return c.singleUUID(ctx, GetParamsVersionV4)
}

// UUIDv7 fetches a UUID v7 value.
func (c *Client) UUIDv7(ctx context.Context) (string, error) {
	return c.singleUUID(ctx, GetParamsVersionV7)
}

// ULID fetches a ULID value.
func (c *Client) ULID(ctx context.Context) (string, error) {
	var resp struct {
		ULID string `json:"ulid"`
	}
	params := &GetParams{Version: ptrVersion(GetParamsVersionUlid)}
	if err := c.invoke(ctx, params, &resp); err != nil {
		return "", err
	}
	return resp.ULID, nil
}

// UUIDBatch fetches multiple UUIDs of the given version.
func (c *Client) UUIDBatch(ctx context.Context, version string, count int) ([]string, error) {
	ver := GetParamsVersion(version)
	if !isSupportedUUIDVersion(ver) {
		return nil, fmt.Errorf("version must be one of v1, v4, v7")
	}
	if count <= 0 || count > 1000 {
		return nil, fmt.Errorf("count must be between 1 and 1000")
	}

	params := &GetParams{
		Version: ptrVersion(ver),
		Count:   ptrCount(count),
	}

	if count == 1 {
		id, err := c.singleUUID(ctx, ver)
		if err != nil {
			return nil, err
		}
		return []string{id}, nil
	}

	var resp struct {
		UUIDs []string `json:"uuids"`
	}
	if err := c.invoke(ctx, params, &resp); err != nil {
		return nil, err
	}
	return resp.UUIDs, nil
}

// ULIDBatch fetches multiple ULIDs.
func (c *Client) ULIDBatch(ctx context.Context, count int) ([]string, error) {
	if count <= 0 || count > 1000 {
		return nil, fmt.Errorf("count must be between 1 and 1000")
	}

	params := &GetParams{
		Version: ptrVersion(GetParamsVersionUlid),
		Count:   ptrCount(count),
	}

	if count == 1 {
		id, err := c.ULID(ctx)
		if err != nil {
			return nil, err
		}
		return []string{id}, nil
	}

	var resp struct {
		ULIDs []string `json:"ulids"`
	}
	if err := c.invoke(ctx, params, &resp); err != nil {
		return nil, err
	}
	return resp.ULIDs, nil
}

func (c *Client) singleUUID(ctx context.Context, version GetParamsVersion) (string, error) {
	var resp struct {
		UUID string `json:"uuid"`
	}
	params := &GetParams{Version: ptrVersion(version)}
	if err := c.invoke(ctx, params, &resp); err != nil {
		return "", err
	}
	return resp.UUID, nil
}

func (c *Client) invoke(ctx context.Context, params *GetParams, v interface{}) error {
	if c == nil {
		return &RequestError{Err: errors.New("client is nil")}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := c.Get(ctx, params)
	if err != nil {
		return &RequestError{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		msg := readBodySnippet(resp.Body)
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return &APIError{StatusCode: resp.StatusCode, Message: msg}
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

func isSupportedUUIDVersion(version GetParamsVersion) bool {
	switch version {
	case GetParamsVersionV1, GetParamsVersionV4, GetParamsVersionV7:
		return true
	default:
		return false
	}
}

func ptrVersion(v GetParamsVersion) *GetParamsVersion {
	vv := v
	return &vv
}

func ptrCount(count int) *int {
	c := count
	return &c
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
