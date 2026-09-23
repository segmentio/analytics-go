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
