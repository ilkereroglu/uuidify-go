package uuidify

import "fmt"

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

// Unwrap allows DecodeError to participate in errors.Is/errors.As checks.
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

// Unwrap allows RequestError to expose the original error.
func (e *RequestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
