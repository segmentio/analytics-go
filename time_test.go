package analytics

import (
	"testing"
	"time"
)

func TestFormatTimeZero(t *testing.T) {
	if s := formatTime(time.Time{}); len(s) != 0 {
		t.Errorf("invalid formatting of zero time value:", s)
	}
}

func TestFormatTimeNonZero(t *testing.T) {
	if s := formatTime(mockTime()); s != "2009-11-10T23:00:00+0000" {
		t.Errorf("invalid formatting of non-zero time value:", s)
	}
}

func TestMakeTimeDefault(t *testing.T) {
	ts := mockTime()

	if tt := makeTime(time.Time{}, ts); tt != ts {
		t.Errorf("invalid default time value:", tt)
	}
}

func TestMakeTimeNonDefault(t *testing.T) {
	ts := mockTime()

	if tt := makeTime(ts, time.Now()); tt != ts {
		t.Errorf("invalid non-default time value:", tt)
	}
}
