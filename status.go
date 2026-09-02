package analytics

import (
	"strconv"
	"strings"
	"time"
)

// isSuccess returns true for 2xx and 3xx responses (spec item 1).
func isSuccess(status int) bool {
	return status >= 200 && status < 400
}

// retryableStatus returns whether the given HTTP status code is retryable.
func retryableStatus(status int) bool {
	switch status {
	case 408, 410, 429, 460:
		return true
	case 501, 505, 511:
		return false
	default:
		return status >= 500 && status < 600
	}
}

// parseRetryAfter parses the Retry-After header value.
// Supports integer seconds and HTTP-date format (RFC 7231 §7.1.1.1).
// Returns 0 if the value is absent, invalid, zero, or in the past.
// Caps the value at cap.
func parseRetryAfter(header string, cap int64) int64 {
	if header == "" {
		return 0
	}
	header = strings.TrimSpace(header)
	// Try integer seconds first
	n, err := strconv.ParseInt(header, 10, 64)
	if err == nil {
		if n <= 0 {
			return 0
		}
		if n > cap {
			return cap
		}
		return n
	}
	// Try HTTP-date format (RFC 7231 §7.1.1.1)
	t, err := time.Parse(time.RFC1123, header)
	if err != nil {
		// Also try RFC1123Z (with numeric timezone)
		t, err = time.Parse(time.RFC1123Z, header)
		if err != nil {
			return 0
		}
	}
	seconds := int64(time.Until(t).Seconds())
	if seconds <= 0 {
		return 0
	}
	if seconds > cap {
		return cap
	}
	return seconds
}
