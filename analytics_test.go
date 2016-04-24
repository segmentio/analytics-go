package analytics

import (
	"net/http/httptest"
	"testing"
)
import "encoding/json"
import "net/http"
import "bytes"
import "time"
import "fmt"
import "io"

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
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 1,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
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
	//       "timestamp": "2009-11-10T23:00:00+0000",
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
	//   "sentAt": "2009-11-10T23:00:00+0000"
	// }
}

func TestTrack(t *testing.T) {
	body, server := mockServer()
	defer server.Close()

	client, _ := NewWithConfig("h97jamjwbh", Config{
		Endpoint:  server.URL,
		Verbose:   true,
		Logger:    t,
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 1,
	})
	defer client.Close()

	if err := client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
	}); err != nil {
		t.Error(err)
		return
	}

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
      "timestamp": "2009-11-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
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
		Now:      mockTime,
		UID:      mockId,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
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
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 1,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
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
      "timestamp": "2015-07-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
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
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 1,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
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
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 1,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
		Context: Context{
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
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
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 3,
	})
	defer client.Close()

	for i := 0; i < 5; i++ {
		client.Enqueue(Track{
			Event:  "Download",
			UserId: "123456",
			Properties: map[string]interface{}{
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
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
		Now:       mockTime,
		UID:       mockId,
		BatchSize: 1,
	})
	defer client.Close()

	client.Enqueue(Track{
		Event:  "Download",
		UserId: "123456",
		Properties: map[string]interface{}{
			"application": "Segment Desktop",
			"version":     "1.1.0",
			"platform":    "osx",
		},
		Integrations: map[string]interface{}{
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
      "timestamp": "2009-11-10T23:00:00+0000",
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
  "sentAt": "2009-11-10T23:00:00+0000"
}`

	if res := string(<-body); ref != res {
		t.Errorf("invalid response:\n- expected %s\n- received: %s", ref, res)
	}
}

func TestCloseTwice(t *testing.T) {
	client := New("0123456789")

	if err := client.Close(); err != nil {
		t.Error("closing a client should not a return an error")
	}

	if err := client.Close(); err != io.EOF {
		t.Error("closing a client a second time should return io.EOF:", err)
	}

	if err := client.Enqueue(Track{UserId: "1", Event: "A"}); err != io.EOF {
		t.Error("using a client after it was closed should return io.EOF:", err)
	}
}
