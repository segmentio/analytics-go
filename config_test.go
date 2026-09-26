package analytics

import (
	"testing"
	"time"
)

func TestConfigZeroValue(t *testing.T) {
	c := Config{}

	if err := c.validate(); err != nil {
		t.Error("validating the zero-value configuration failed:", err)
	}
}

func TestConfigInvalidInterval(t *testing.T) {
	c := Config{
		Interval: -1 * time.Second,
	}

	if err := c.validate(); err == nil {
		t.Error("no error returned when validating a malformed config")

	} else if e, ok := err.(ConfigError); !ok {
		t.Error("invalid error returned when checking a malformed config:", err)

	} else if e.Field != "Interval" || e.Value.(time.Duration) != (-1*time.Second) {
		t.Error("invalid field error reported:", e)
	}
}

func TestConfigInvalidBatchSize(t *testing.T) {
	c := Config{
		BatchSize: -1,
	}

	if err := c.validate(); err == nil {
		t.Error("no error returned when validating a malformed config")

	} else if e, ok := err.(ConfigError); !ok {
		t.Error("invalid error returned when checking a malformed config:", err)

	} else if e.Field != "BatchSize" || e.Value.(int) != -1 {
		t.Error("invalid field error reported:", e)
	}
}

func TestDefaultRetryAfterNeverExceedsTheCeiling(t *testing.T) {
	for attempt := 0; attempt < 30; attempt++ {
		if d := defaultRetryAfter(attempt); d > 60*time.Second {
			t.Fatalf("attempt %d returned %s, above the 60s ceiling", attempt, d)
		}
	}
}

func TestDefaultRetryAfterJittersAtTheCeiling(t *testing.T) {
	// backo-go jittered before clamping, so every attempt past the ceiling
	// returned exactly the cap and a fleet stayed in lockstep. Guard against
	// regressing to that.
	seen := make(map[time.Duration]struct{})
	var min time.Duration = 60 * time.Second

	for i := 0; i < 50; i++ {
		d := defaultRetryAfter(20) // well past the ceiling
		seen[d] = struct{}{}
		if d < min {
			min = d
		}
	}

	if len(seen) < 2 {
		t.Errorf("expected jittered values at the ceiling, got the same value %d times", len(seen))
	}
	if min < 30*time.Second {
		t.Errorf("jitter should subtract at most 50%%, but saw %s", min)
	}
}

func TestConfigRejectsNegativeRetryFields(t *testing.T) {
	// These used to survive validate() and reach the retry loop, where a negative
	// MaxRetries dropped every batch after its first failure instead of erroring here.
	for _, test := range []struct {
		field  string
		config Config
	}{
		{"MaxRetries", Config{MaxRetries: -1}},
		{"MaxTotalBackoffDuration", Config{MaxTotalBackoffDuration: -1 * time.Second}},
		{"MaxRateLimitDuration", Config{MaxRateLimitDuration: -1 * time.Second}},
		{"ShutdownTimeout", Config{ShutdownTimeout: -1 * time.Second}},
	} {
		if err := test.config.validate(); err == nil {
			t.Errorf("negative %s should be rejected", test.field)
		}
	}
}

func TestConfigZeroRetryFieldsTakeDefaults(t *testing.T) {
	// Zero means "use the default" for these, per Config's zero-value convention.
	// There is deliberately no way to ask for no retries at all.
	c := Config{}
	if err := c.validate(); err != nil {
		t.Fatalf("zero values should be valid: %s", err)
	}

	c = makeConfig(c)
	if c.MaxRetries != DefaultMaxRetries {
		t.Errorf("MaxRetries = %d, want the default %d", c.MaxRetries, DefaultMaxRetries)
	}
	if c.ShutdownTimeout != DefaultShutdownTimeout {
		t.Errorf("ShutdownTimeout = %s, want the default %s", c.ShutdownTimeout, DefaultShutdownTimeout)
	}
	if c.MaxRateLimitDuration != DefaultMaxRateLimitDuration {
		t.Errorf("MaxRateLimitDuration = %s, want the default %s", c.MaxRateLimitDuration, DefaultMaxRateLimitDuration)
	}
	if c.MaxTotalBackoffDuration != DefaultMaxTotalBackoffDuration {
		t.Errorf("MaxTotalBackoffDuration = %s, want the default %s", c.MaxTotalBackoffDuration, DefaultMaxTotalBackoffDuration)
	}
}

func TestDefaultRateLimitBudgetExceedsTheRetryAfterCeiling(t *testing.T) {
	// A single maximal Retry-After must not be able to consume the whole budget.
	// When the two are equal, the elapsed check runs before the wait, so one wait
	// spends the budget and the batch is dropped having been attempted once.
	ceiling := time.Duration(maxRetryAfterSeconds) * time.Second
	if DefaultMaxRateLimitDuration <= ceiling {
		t.Fatalf("DefaultMaxRateLimitDuration (%s) must exceed the Retry-After ceiling (%s); "+
			"at parity a single capped Retry-After leaves no room for a retry",
			DefaultMaxRateLimitDuration, ceiling)
	}
	if attempts := DefaultMaxRateLimitDuration / ceiling; attempts < 2 {
		t.Errorf("budget allows only %d capped wait(s); want room for at least 2 attempts", attempts)
	}
}
