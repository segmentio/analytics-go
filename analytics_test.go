package analytics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func fixture(name string) string {
	f, err := os.Open(filepath.Join("fixtures", name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mockId() string { return "I'm unique" }

func mockTime() time.Time {
	// time.Unix(0, 0) fails on Circle
	return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
}

func mockServer() (chan []byte, *httptest.Server) {
	done := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, r.Body)

		var v interface{}
		err := json.Unmarshal(buf.Bytes(), &v)
		if err != nil {
			panic(err)
		}

		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			panic(err)
		}

		done <- b
	}))

	return done, server
}

func ExampleTrack() {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		BatchSize: 1,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: Properties{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
	})

	fmt.Printf("%s\n", <-body)
	// Output:
	// {
	//   "batch": [
	//     {
	//       "event": "Download",
	//       "messageId": "I'm unique",
	//       "properties": {
	//         "application": "Segment Desktop",
	//         "platform": "osx",
	//         "version": "1.1.0"
	//       },
	//       "timestamp": "2009-11-10T23:00:00Z",
	//       "type": "track",
	//       "userId": "123456"
	//     }
	//   ],
	//   "context": {
	//     "library": {
	//       "name": "analytics-go",
	//       "version": "3.0.0"
	//     }
	//   },
	//   "messageId": "I'm unique",
	//   "sentAt": "2009-11-10T23:00:00Z"
	// }
}

func TestEnqueue(t *testing.T) {
	tests := map[string]struct {
		ref string
		msg Message
	}{
		"alias": {
			fixture("test-enqueue-alias.json"),
			Alias{PreviousId: "A", UserId: "B"},
		},

		"group": {
			fixture("test-enqueue-group.json"),
			Group{GroupId: "A", UserId: "B"},
		},

		"identify": {
			fixture("test-enqueue-identify.json"),
			Identify{UserId: "B"},
		},

		"page": {
			fixture("test-enqueue-page.json"),
			Page{Name: "A", UserId: "B"},
		},

		"screen": {
			fixture("test-enqueue-screen.json"),
			Screen{Name: "A", UserId: "B"},
		},

		"track": {
			fixture("test-enqueue-track.json"),
			Track{
				Event:  "Download",
				UserId: "123456",
				Properties: Properties{
					"application": "Segment Desktop",
					"version":     "1.1.0",
					"platform":    "osx",
				},
			},
		},
	}

	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		BatchSize: 1,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	for name, test := range tests {
		if err := client.Enqueue(test.msg); err != nil {
			t.Error(err)
			return
		}

		if res := string(<-body); res != test.ref {
			t.Errorf("%s: invalid response:\n- expected %s\n- received: %s", name, test.ref, res)
		}
	}
}

func TestTrackWithInterval(t *testing.T) {
	const interval = 100 * time.Millisecond

	body, server := mockServer()
	defer server.Close()

	t0 := time.Now()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint: server.URL,
		Interval: interval,
		Verbose:  true,
		Logger:   t,
		now:      mockTime,
		uid:      mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: Properties{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
	})

	const ref = `{
  "batch": [
    {
      "event": "Download",
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "platform": "osx",
        "version": "1.1.0"
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    }
  ],
  "context": {
    "library": {
      "name": "analytics-go",
      "version": "3.0.0"
    }
  },
  "messageId": "I'm unique",
  "sentAt": "2009-11-10T23:00:00Z"
}`

	// Will flush in 100 milliseconds
	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}

	if t1 := time.Now(); t1.Sub(t0) < interval {
		t.Error("the flushing interval is too short:", interval)
	}
}

func TestTrackWithTimestamp(t *testing.T) {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		BatchSize: 1,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: Properties{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
		Timestamp: time.Date(2015, time.July, 10, 23, 0, 0, 0, time.UTC),
	})

	const ref = `{
  "batch": [
    {
      "event": "Download",
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "platform": "osx",
        "version": "1.1.0"
      },
      "timestamp": "2015-07-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    }
  ],
  "context": {
    "library": {
      "name": "analytics-go",
      "version": "3.0.0"
    }
  },
  "messageId": "I'm unique",
  "sentAt": "2009-11-10T23:00:00Z"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}
}

func TestTrackWithMessageId(t *testing.T) {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		BatchSize: 1,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: Properties{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
		MessageId: "abc",
	})

	const ref = `{
  "batch": [
    {
      "event": "Download",
      "messageId": "abc",
      "properties": {
        "application": "Segment Desktop",
        "platform": "osx",
        "version": "1.1.0"
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    }
  ],
  "context": {
    "library": {
      "name": "analytics-go",
      "version": "3.0.0"
    }
  },
  "messageId": "I'm unique",
  "sentAt": "2009-11-10T23:00:00Z"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}
}

func TestTrackWithContext(t *testing.T) {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		BatchSize: 1,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: Properties{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
		Context: &Context{
			Extra: map[string]interface{}{
				"whatever": "here",
			},
		},
	})

	const ref = `{
  "batch": [
    {
      "context": {
        "whatever": "here"
      },
      "event": "Download",
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "platform": "osx",
        "version": "1.1.0"
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    }
  ],
  "context": {
    "library": {
      "name": "analytics-go",
      "version": "3.0.0"
    }
  },
  "messageId": "I'm unique",
  "sentAt": "2009-11-10T23:00:00Z"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}
}

func TestTrackMany(t *testing.T) {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		BatchSize: 3,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	for i := 0; i < 5; i++ {
		client.Enqueue(Track{
			Event:  "Download",
			UserId: "123456",
			Properties: Properties{
				"application": "Segment Desktop",
				"version":     i,
			},
		})
	}

	const ref = `{
  "batch": [
    {
      "event": "Download",
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "version": 0
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    },
    {
      "event": "Download",
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "version": 1
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    },
    {
      "event": "Download",
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "version": 2
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    }
  ],
  "context": {
    "library": {
      "name": "analytics-go",
      "version": "3.0.0"
    }
  },
  "messageId": "I'm unique",
  "sentAt": "2009-11-10T23:00:00Z"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}
}

func TestTrackWithIntegrations(t *testing.T) {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		BatchSize: 1,
		now:       mockTime,
		uid:       mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: Properties{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
		Integrations: Integrations{
			"All":      true,
			"Intercom": false,
			"Mixpanel": true,
		},
	})

	const ref = `{
  "batch": [
    {
      "event": "Download",
      "integrations": {
        "All": true,
        "Intercom": false,
        "Mixpanel": true
      },
      "messageId": "I'm unique",
      "properties": {
        "application": "Segment Desktop",
        "platform": "osx",
        "version": "1.1.0"
      },
      "timestamp": "2009-11-10T23:00:00Z",
      "type": "track",
      "userId": "123456"
    }
  ],
  "context": {
    "library": {
      "name": "analytics-go",
      "version": "3.0.0"
    }
  },
  "messageId": "I'm unique",
  "sentAt": "2009-11-10T23:00:00Z"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}
}

func TestClientCloseTwice(t *testing.T) {
	client := New("0123456789")

	if err := client.Close(); err != nil {
		t.Error("closing a client should not a return an error")
	}

	if err := client.Close(); err != ErrClosed {
		t.Error("closing a client a second time should return ErrClosed:", err)
	}

	if err := client.Enqueue(Track{UserId: "1", Event: "A"}); err != ErrClosed {
		t.Error("using a client after it was closed should return ErrClosed:", err)
	}
}

func TestClientConfigError(t *testing.T) {
	client, err := NewWithConfig("0123456789", Config{
		Interval: -1 * time.Second,
	})

	if err == nil {
		t.Error("no error returned when creating a client with an invalid config")
	}

	if _, ok := err.(ConfigError); !ok {
		t.Errorf("invalid error type returned when creating a client with an invalid config: %T", err)
	}

	if client != nil {
		t.Error("invalid non-nil client object returned when creating a client with and invalid config:", client)
		client.Close()
	}
}

func TestClientEnqueueError(t *testing.T) {
	client := New("0123456789")
	defer client.Close()

	if err := client.Enqueue(testErrorMessage{}); err != testError {
		t.Error("invlaid error returned when queueing an invalid message:", err)
	}
}

var testError = errors.New("test error")

type testErrorMessage struct{}

func (m testErrorMessage) validate() error {
	return testError
}
