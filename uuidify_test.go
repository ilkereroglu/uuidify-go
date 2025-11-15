package uuidify

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUUIDv4_Single(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("version"); got != "v4" {
			t.Fatalf("expected version v4, got %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"uuid":"1234"}`)
	}))
	defer ts.Close()

c := newTestClient(t, ts)

	uuid, err := c.UUIDv4(context.Background())
	if err != nil {
		t.Fatalf("UUIDv4 returned error: %v", err)
	}
	if uuid != "1234" {
		t.Fatalf("expected uuid 1234, got %s", uuid)
	}
}

func TestUUIDv7_Batch(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("version"); got != "v7" {
			t.Fatalf("expected version v7, got %s", got)
		}
		if got := q.Get("count"); got != "5" {
			t.Fatalf("expected count 5, got %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"uuids":["a","b","c","d","e"]}`)
	}))
	defer ts.Close()

c := newTestClient(t, ts)

	uuids, err := c.UUIDBatch(context.Background(), "v7", 5)
	if err != nil {
		t.Fatalf("UUIDBatch returned error: %v", err)
	}
	if len(uuids) != 5 {
		t.Fatalf("expected 5 uuids, got %d", len(uuids))
	}
}

func TestULID_Single(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("version"); got != "ulid" {
			t.Fatalf("expected version ulid, got %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"ulid":"01J123"}`)
	}))
	defer ts.Close()

c := newTestClient(t, ts)

	id, err := c.ULID(context.Background())
	if err != nil {
		t.Fatalf("ULID returned error: %v", err)
	}
	if id != "01J123" {
		t.Fatalf("expected ulid 01J123, got %s", id)
	}
}

func TestULID_Batch(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("version"); got != "ulid" {
			t.Fatalf("expected version ulid, got %s", got)
		}
		if got := q.Get("count"); got != "3" {
			t.Fatalf("expected count 3, got %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"ulids":["a","b","c"]}`)
	}))
	defer ts.Close()

c := newTestClient(t, ts)

	ids, err := c.ULIDBatch(context.Background(), 3)
	if err != nil {
		t.Fatalf("ULIDBatch returned error: %v", err)
	}
	if len(ids) != 3 {
		t.Fatalf("expected 3 ulids, got %d", len(ids))
	}
}

func TestError_Transport(t *testing.T) {
	t.Parallel()

	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	})}

c, err := NewClient("https://example.com", WithHTTPClient(client))
if err != nil {
	t.Fatalf("failed to create client: %v", err)
}

	if _, err := c.UUIDv4(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	} else {
		var reqErr *RequestError
		if !errors.As(err, &reqErr) {
			t.Fatalf("expected RequestError, got %T", err)
		}
	}
}

func TestError_Decode(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"uuid":`)
	}))
	defer ts.Close()

c := newTestClient(t, ts)

	if _, err := c.UUIDv4(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	} else {
		var decErr *DecodeError
		if !errors.As(err, &decErr) {
			t.Fatalf("expected DecodeError, got %T", err)
		}
	}
}

func TestError_APIStatus(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error":"bad request"}`)
	}))
	defer ts.Close()

c := newTestClient(t, ts)

	if _, err := c.UUIDv4(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	} else {
		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
		}
	}
}

func newTestClient(t *testing.T, ts *httptest.Server) *Client {
	t.Helper()
	client, err := NewClient(
		ts.URL,
		WithHTTPClient(ts.Client()),
		WithUserAgent("uuidify-go-tests"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
