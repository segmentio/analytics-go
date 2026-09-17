package analytics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsSuccess(t *testing.T) {
	cases := []struct {
		status int
		want   bool
	}{
		{200, true}, {201, true}, {204, true}, {301, true}, {302, true},
		{400, false}, {429, false}, {500, false}, {0, false}, {199, false},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, isSuccess(tc.status), "status %d", tc.status)
	}
}

func TestRetryableStatus(t *testing.T) {
	cases := []struct {
		status    int
		retryable bool
	}{
		{408, true},
		{410, true},
		{429, true},
		{460, true},
		{500, true},
		{502, true},
		{503, true},
		{504, true},
		{508, true},
		{529, true},
		{501, false},
		{505, false},
		{511, false},
		{400, false},
		{401, false},
		{403, false},
		{413, false},
		{200, false},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.retryable, retryableStatus(tc.status), "retryable for status %d", tc.status)
	}
}

func TestParseRetryAfter(t *testing.T) {
	assert.Equal(t, int64(60), parseRetryAfter("60", 300))
	assert.Equal(t, int64(300), parseRetryAfter("9999", 300)) // capped
	assert.Equal(t, int64(0), parseRetryAfter("0", 300))
	assert.Equal(t, int64(0), parseRetryAfter("-1", 300))
	assert.Equal(t, int64(0), parseRetryAfter("", 300))
	assert.Equal(t, int64(1), parseRetryAfter("1", 300))
	assert.Equal(t, int64(300), parseRetryAfter("300", 300))
}

func TestParseRetryAfterHTTPDate(t *testing.T) {
	// Date ~2 seconds in the future should return ~2
	future := time.Now().Add(2 * time.Second).UTC().Format(time.RFC1123)
	result := parseRetryAfter(future, 300)
	assert.True(t, result >= 1 && result <= 3, "expected ~2 seconds, got %d", result)

	// Date in the past should return 0
	past := time.Now().Add(-10 * time.Second).UTC().Format(time.RFC1123)
	assert.Equal(t, int64(0), parseRetryAfter(past, 300))

	// Garbage string should return 0
	assert.Equal(t, int64(0), parseRetryAfter("not-a-date-or-number", 300))
}
