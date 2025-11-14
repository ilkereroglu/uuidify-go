package uuidify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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
		return "uuidify API error (" + strconv.Itoa(e.StatusCode) + "): " + e.Message
	}
	return "uuidify API error (" + strconv.Itoa(e.StatusCode) + ")"
}

// DecodeError wraps errors that occur while decoding API responses.
type DecodeError struct {
	Err error
}

func (e *DecodeError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return "uuidify decode error: " + e.Err.Error()
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
	return "uuidify request error: " + e.Err.Error()
}

func (e *RequestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
