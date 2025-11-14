package uuidify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
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

	userAgent := c.UserAgent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	req.Header.Set("User-Agent", userAgent)

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
