//go:build live

package uuidify

import (
	"context"
	"testing"
)

func TestLive_UUIDv4(t *testing.T) {
	c := NewClient()

	uuid, err := c.UUIDv4(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if uuid == "" {
		t.Fatal("empty uuid")
	}
}
