package analytics

import (
	"errors"
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

// Identify should return an expected error if "userId" or "anonymousId" doesn't provided
func TestIdentify(t *testing.T) {
	valid_client := New("this is a secret, you can't read it. don't you?")
	cases := []struct {
		client       *Client
		msg          Message
		expected_err error
	}{
		{
			valid_client,
			Message{"userId": "1", "traits": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"anonymousId": "1", "traits": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"traits": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
	}

	for _, c := range cases {
		err := valid_client.Identify(c.msg)
		if (err == nil && c.expected_err != nil) || (err != nil && c.expected_err == nil) {
			t.Errorf("got: %v, expected: %v", err, c.expected_err)

		}
		if err != nil && c.expected_err != nil {
			if err.Error() != c.expected_err.Error() {
				t.Errorf("got: %v, expected: %v", err, c.expected_err)

			}
		}

	}
}
