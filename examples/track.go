package main

import "github.com/segmentio/analytics-go"
import "time"

func main() {
	client := analytics.New("h97jamjwbh")
	client.Interval = 30 * time.Second
	client.Verbose = true
	client.Size = 25

	for {
		client.Track(&analytics.Track{
			Event:  "Download",
			UserId: "123456",
			Properties: map[string]interface{}{
				"application": "Segment Desktop",
				"version":     "1.1.0",
				"platform":    "osx",
			},
		})

		time.Sleep(50 * time.Millisecond)
	}
}
