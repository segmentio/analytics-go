package analytics

import (
	"errors"
	"testing"
	"time"
)

var (
	client = New("SECRET")
)

// TestCases represent case on test
type TestCases struct {
	msg          Message
	expected_err error
}

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
	cases := []TestCases{
		{Message{"userId": "1", "traits": map[string]interface{}{}}, nil},
		{Message{"anonymousId": "1", "traits": map[string]interface{}{}}, nil},
		{Message{"traits": map[string]interface{}{}}, errors.New("You must pass either an 'anonymousId' or 'userId'.")},
	}

	LoopAssertion(t, cases, client.Identify)
}

// Track should return an expected error if the following are true:
// 1. "event" doesn't exists
// 2. "userId" or "anonymousId" doesn't exists
func TestTrack(t *testing.T) {
	cases := []TestCases{
		{
			Message{"userId": "1", "event": "registered",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			Message{"anonymousId": "1", "event": "registered",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			Message{"event": "registered", "properties": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
		{
			Message{"userId": "1", "properties": map[string]interface{}{}},
			errors.New("You must pass 'event'."),
		},
		{
			Message{"anonymousId": "1", "properties": map[string]interface{}{}},
			errors.New("You must pass 'event'."),
		},
	}

	LoopAssertion(t, cases, client.Track)
}

// Page should return an expected error if userId" or "anonymousId" doesn't exists
func TestPage(t *testing.T) {
	cases := []TestCases{
		{
			Message{"userId": "1", "name": "About", "category": "Help",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			Message{"anonymousId": "1", "name": "About", "category": "Help",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			Message{"name": "Home", "properties": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
	}
	LoopAssertion(t, cases, client.Page)
}

// Screen should return an expected error if userId" or "anonymousId" doesn't exists
func TestScreen(t *testing.T) {
	cases := []TestCases{
		{
			Message{"userId": "1", "name": "Cover", "category": "Guide",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			Message{"anonymousId": "1", "name": "About", "category": "About",
				"properties": map[string]interface{}{}},
			nil,
		},
		{
			Message{"name": "Post New", "properties": map[string]interface{}{}},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
	}
	LoopAssertion(t, cases, client.Screen)
}

// Group should return an expected error if the following are true:
// 1. "groupId" doesn't exists
// 2. "userId" or "anonymousId" doesn't exists
func TestGroup(t *testing.T) {
	cases := []TestCases{
		{
			Message{
				"userId":  "1",
				"groupId": "1",
				"traits":  map[string]interface{}{},
			},
			nil,
		},
		{
			Message{
				"anonymousId": "1",
				"groupId":     "1",
				"traits":      map[string]interface{}{},
			},
			nil,
		},
		{
			Message{
				"groupId": "1",
				"traits":  map[string]interface{}{},
			},
			errors.New("You must pass either an 'anonymousId' or 'userId'."),
		},
		{
			Message{
				"userId": "1",
				"traits": map[string]interface{}{},
			},
			errors.New("You must pass a 'groupId'."),
		},
		{
			Message{
				"anonymousId": "1",
				"traits":      map[string]interface{}{},
			},
			errors.New("You must pass a 'groupId'."),
		},
	}
	LoopAssertion(t, cases, client.Group)
}

// Alias should return an expected error if "userId" or "previousId" doesn't exists
func TestAlias(t *testing.T) {
	cases := []TestCases{
		{Message{"userId": "1", "previousId": "1"}, nil},
		{Message{"userId": "1"}, errors.New("You must pass a 'previousId'.")},
		{Message{"previousId": "1"}, errors.New("You must pass a 'userId'.")},
		{Message{"someId": "1", "otherId": "1"}, errors.New("You must pass a 'userId'.")},
	}

	LoopAssertion(t, cases, client.Alias)
}

func LoopAssertion(t *testing.T, cases []TestCases, exported func(msg Message) error) {
	for _, c := range cases {
		err := exported(c.msg)
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
