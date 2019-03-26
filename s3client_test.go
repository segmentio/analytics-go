package analytics

import (
	"encoding/json"
	"testing"
)

func TestTargetMessageMarshalling(t *testing.T) {
	m := Track{
		Event:  "FooBared",
		UserId: "tuna",
		Properties: map[string]interface{}{
			"index": 1,
			"qwer":  3424,
		},
	}
	tm, err := makeTargetMessage(m, 10000, nil, func() Time { return Time{} })
	if err != nil {
		t.Error(err)
	}
	b, err := json.Marshal(tm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("json: %s", string(b))

	expected := `{"event":{"userId":"tuna","event":"FooBared","timestamp":0,"properties":{"index":1,"qwer":3424}},"sentAt":0,"receivedAt":0}`

	if string(b) != expected {
		t.Errorf("Expected: %s, Actual: %s", expected, string(b))
	}
}

func TestS3Client(t *testing.T) {
	c, err := NewS3ClientWithConfig(
		S3ClientConfig{
			Stream: "tuna",
			Stage:  "pavel",
		},
		Config{
			Verbose: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		m := Track{
			Event:  "FooBared",
			UserId: "tuna",
			Properties: map[string]interface{}{
				"index": i,
				"qwer":  3424,
			},
		}
		if err := c.Enqueue(m); err != nil {
			t.Error(err)
		}
	}
	if err := c.Close(); err != nil {
		t.Error(err)
	}

	t.FailNow()
}
