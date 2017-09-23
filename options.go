package analytics

import (
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// DefaultInterval default interval
	DefaultInterval = 5 * time.Second

	// DefaultSize default size
	DefaultSize = 250

	// DefaultLogger default goes to Stderr with 'segment' prefix
	DefaultLogger = log.New(os.Stderr, "segment ", log.LstdFlags)

	// DefaultVerbosity is false
	DefaultVerbosity = false

	// DefaultHTTPClient is http.DefaultClient
	DefaultHTTPClient = *http.DefaultClient
)

// Options contains library options
type Options struct {
	Interval time.Duration
	Size     int
	Logger   Logger
	Verbose  bool
	Client   http.Client
}

// Option function for setup some opts
type Option func(opts *Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Interval: DefaultInterval,
		Size:     DefaultSize,
		Logger:   DefaultLogger,
		Verbose:  DefaultVerbosity,
		Client:   DefaultHTTPClient,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Interval setup interval
func Interval(interval time.Duration) Option {
	return func(o *Options) {
		o.Interval = interval
	}
}

// Size setup size
func Size(size int) Option {
	return func(o *Options) {
		o.Size = size
	}
}

// Log setup custom logger
func Log(logger Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// Verbose setup verbosity
func Verbose(v bool) Option {
	return func(o *Options) {
		o.Verbose = v
	}
}

// HTTPClient setup http.Client to use
func HTTPClient(client http.Client) Option {
	return func(o *Options) {
		o.Client = client
	}
}
