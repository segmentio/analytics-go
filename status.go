package analytics

import (
	"strconv"
	"strings"
)

// isSuccess returns true for 2xx and 3xx responses (spec item 1).
func isSuccess(status int) bool {
	return status >= 200 && status < 400
}

// retryableStatus returns the retry strategy for a given HTTP status code.
// Returns (retryable bool, isRateLimit bool).
func retryableStatus(status int) (retryable bool, isRateLimit bool) {
	switch status {
	case 429:
		return true, true
	case 408, 410, 460:
		return true, false
	case 501, 505, 511:
		return false, false
	default:
		if status >= 500 && status < 600 {
			return true, false
		}
		return false, false
	}
}

// parseRetryAfter parses the Retry-After header value (integer seconds only).
// Returns 0 if the value is absent, invalid, zero, or an HTTP-date.
// Caps the value at cap.
func parseRetryAfter(header string, cap int64) int64 {
	if header == "" {
		return 0
	}
	n, err := strconv.ParseInt(strings.TrimSpace(header), 10, 64)
	if err != nil || n <= 0 {
		return 0
	}
	if n > cap {
		return cap
	}
	return n
}
