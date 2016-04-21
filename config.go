package analytics

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ClientConfigFunc func(*Client) error

// Overwrite the default endpoint 'https://api.segment.io'
func WithEndpoint(endpoint string) ClientConfigFunc {
	return func(c *Client) error {
		if endpoint == "" {
			return fmt.Errorf("Endpoint: invalid endpoint '%s'", endpoint)
		}
		c.endpoint = endpoint
		return nil
	}
}

// Overwrite the default flush interval of 5 seconds
func WithFlushInterval(interval time.Duration) ClientConfigFunc {
	return func(c *Client) error {
		if interval <= 0 {
			return fmt.Errorf("FlushInterval: invalid interval '%d'", interval)
		}
		c.interval = interval
		return nil
	}
}

// Overwrite the default buffersize (Note: Segment will reject payloads larger than 500 KB)
func WithBufferSize(bufferSize int) ClientConfigFunc {
	return func(c *Client) error {
		if bufferSize <= 0 {
			return fmt.Errorf("BufferSize: invalid bufferSize '%d'", bufferSize)
		}
		c.bufferSize = bufferSize
		return nil
	}
}

// Overwrite the default http.DefaultClient
func WithHttpClient(client *http.Client) ClientConfigFunc {
	return func(c *Client) error {
		if client == nil {
			return errors.New("HttpClient: <nil> client")
		}
		c.client = client
		return nil
	}
}

// Overwrite the default logger
func WithLogger(logger *log.Logger) ClientConfigFunc {
	return func(c *Client) error {
		if logger == nil {
			return errors.New("Logger: <nil> logger")
		}
		c.logger = logger
		return nil
	}
}

// Configure client for verbose mode
func WithVerbosity(c *Client) error {
	c.verbose = true
	return nil
}
