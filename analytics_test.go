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

// Identify should return an expected error if "userId" or "anonymousId" doesn't exists
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

// Track should return an expected error if the following are true:
// 1. "event" doesn't exists
// 2. "userId" or "anonymousId" doesn't exists
func TestTrack(t *testing.T) {
	valid_client := New("this is an another secret. please don't read it.")
	cases := []struct {
		client       *Client
		msg          Message
		expected_err error
	}{
		{
			valid_client,
			Message{"userId": "1", "event": "registered",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"anonymousId": "1", "event": "registered",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"event": "registered", "properties": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
		{
			valid_client,
			Message{"userId": "1", "properties": map[string]interface{}{}},
			errors.New("You must pass 'event'."),
		},
		{
			valid_client,
			Message{"anonymousId": "1", "properties": map[string]interface{}{}},
			errors.New("You must pass 'event'."),
		},
	}

	for _, c := range cases {
		err := valid_client.Track(c.msg)
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

// Page should return an expected error if userId" or "anonymousId" doesn't exists
func TestPage(t *testing.T) {
	valid_client := New("i warn you to not read this. but ... you ... argh ...")
	cases := []struct {
		client       *Client
		msg          Message
		expected_err error
	}{
		{
			valid_client,
			Message{"userId": "1", "name": "About", "category": "Help",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"anonymousId": "1", "name": "About", "category": "Help",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"name": "Home", "properties": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
	}

	for _, c := range cases {
		err := valid_client.Page(c.msg)
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

// Screen should return an expected error if userId" or "anonymousId" doesn't exists
func TestScreen(t *testing.T) {
	valid_client := New("i'm out")
	cases := []struct {
		client       *Client
		msg          Message
		expected_err error
	}{
		{
			valid_client,
			Message{"userId": "1", "name": "Cover", "category": "Guide",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"anonymousId": "1", "name": "About", "category": "About",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			valid_client,
			Message{"name": "Post New", "properties": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
	}

	for _, c := range cases {
		err := valid_client.Screen(c.msg)
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

// Group should return an expected error if the following are true:
// 1. "groupId" doesn't exists
// 2. "userId" or "anonymousId" doesn't exists
func TestGroup(t *testing.T) {
	valid_client := New("OH: this is a secret boy!")
	cases := []struct {
		client       *Client
		msg          Message
		expected_err error
	}{
		{
			valid_client,
			Message{
				"userId":  "1",
				"groupId": "1",
				"traits":  map[string]interface{}{},
			},
			nil,
		},
		{
			valid_client,
			Message{
				"anonymousId": "1",
				"groupId":     "1",
				"traits":      map[string]interface{}{},
			},
			nil,
		},
		{
			valid_client,
			Message{
				"groupId": "1",
				"traits":  map[string]interface{}{},
			},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
		{
			valid_client,
			Message{
				"userId": "1",
				"traits": map[string]interface{}{},
			},
			errors.New("You must pass a 'groupId'."),
		},
		{
			valid_client,
			Message{
				"anonymousId": "1",
				"traits":      map[string]interface{}{},
			},
			errors.New("You must pass a 'groupId'."),
		},
	}

	for _, c := range cases {
		err := valid_client.Group(c.msg)
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

// Alias should return an expected error if "userId" or "previousId" doesn't exists
func TestAlias(t *testing.T) {
	valid_client := New("Sugar ... Yes Please .. (Maroon 5 Sugar)")
	cases := []struct {
		client       *Client
		msg          Message
		expected_err error
	}{
		{
			valid_client,
			Message{"userId": "1", "previousId": "1"},
			nil,
		},
		{
			valid_client,
			Message{"userId": "1"},
			errors.New("You must pass a 'previousId'."),
		},
		{
			valid_client,
			Message{"previousId": "1"},
			errors.New("You must pass a 'userId'."),
		},
		{
			valid_client,
			Message{"someId": "1", "otherId": "1"},
			errors.New("You must pass a 'userId'."),
		},
	}

	for _, c := range cases {
		err := valid_client.Alias(c.msg)
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
