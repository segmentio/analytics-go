package analytics

import (
	"testing"
	"time"
)

// New should return a *Client with default value
func TestNew(t *testing.T) {
	cases := []struct {
		key             string
		expected_client *Client
	}{
		{
			"examplekey",
			&Client{
				FlushAt:    20,
				FlushAfter: 5 * time.Second,
				Key:        "examplekey",
				Endpoint:   "https://api.segment.io",
			},
		},
	}

	for _, c := range cases {
		got := New(c.key)
		if got.FlushAt != c.expected_client.FlushAt {
			t.Errorf("got: %v, expected: %v", got.FlushAt, c.expected_client.FlushAt)
		}
		if got.FlushAfter != c.expected_client.FlushAfter {
			t.Errorf("got: %v, expected: %v", got.FlushAfter, c.expected_client.FlushAfter)
		}
		if got.Key != c.expected_client.Key {
			t.Errorf("got: %v, expected: %v", got.Key, c.expected_client.Key)
		}
		if got.Endpoint != c.expected_client.Endpoint {
			t.Errorf("got: %v, expected: %v", got.Endpoint, c.expected_client.Endpoint)
		}
	}
}
