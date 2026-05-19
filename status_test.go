package analytics

import (
	"testing"

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
		rateLimit bool
	}{
		{408, true, false},
		{410, true, false},
		{429, true, true},
		{460, true, false},
		{500, true, false},
		{502, true, false},
		{503, true, false},
		{504, true, false},
		{508, true, false},
		{501, false, false},
		{505, false, false},
		{511, false, false},
		{400, false, false},
		{401, false, false},
		{403, false, false},
		{413, false, false},
		{200, false, false},
	}
	for _, tc := range cases {
		r, rl := retryableStatus(tc.status)
		assert.Equal(t, tc.retryable, r, "retryable for status %d", tc.status)
		assert.Equal(t, tc.rateLimit, rl, "isRateLimit for status %d", tc.status)
	}
}

func TestParseRetryAfter(t *testing.T) {
	assert.Equal(t, int64(60), parseRetryAfter("60", 300))
	assert.Equal(t, int64(300), parseRetryAfter("9999", 300)) // capped
	assert.Equal(t, int64(0), parseRetryAfter("0", 300))
	assert.Equal(t, int64(0), parseRetryAfter("-1", 300))
	assert.Equal(t, int64(0), parseRetryAfter("", 300))
	assert.Equal(t, int64(0), parseRetryAfter("Wed, 07 May 2026 12:00:00 GMT", 300))
	assert.Equal(t, int64(1), parseRetryAfter("1", 300))
	assert.Equal(t, int64(300), parseRetryAfter("300", 300))
}
