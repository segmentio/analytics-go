package analytics

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestDefaultOptions(t *testing.T) {
	options := newOptions()

	if options.Size != DefaultSize {
		t.Fatalf("Default options 'size' missmatch, have %v want %v", options.Size, DefaultSize)
	}

	if options.Interval != DefaultInterval {
		t.Fatalf("Default options 'interval' missmatch, have %v want %v", options.Interval, DefaultInterval)
	}

	if options.Verbose != DefaultVerbosity {
		t.Fatalf("Default options 'verbose' missmatch, have %v want %v", options.Verbose, DefaultVerbosity)
	}

	if options.Logger == nil {
		t.Fatalf("Default options 'logger' should not be nil")
	}
}

func TestOptionsVars(t *testing.T) {
	size := 42
	interval := 123 * time.Second
	client := &http.Client{}
	logger := log.New(os.Stdout, "logger ", log.Llongfile)

	options := newOptions(
		Size(size),
		Interval(interval),
		Verbose(true),
		HTTPClient(*client),
		Log(logger),
	)

	if options.Size != size {
		t.Fatalf("Options 'size' missmatch, have %v want %v", options.Size, size)
	}

	if options.Interval != interval {
		t.Fatalf("Options 'interval' missmatch, have %v want %v", options.Interval, interval)
	}

	if options.Verbose != true {
		t.Fatalf("Options 'verbose' missmatch, have %v want %v", options.Verbose, true)
	}

	if options.Logger == nil {
		t.Fatalf("Options 'logger' should not be nil")
	}
}
