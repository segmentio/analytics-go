package main

import "github.com/ninnemana/analytics-go"
import "time"
import "log"

// run with DEBUG=analytics to view
// analytics-specific debug output
func main() {
	client := analytics.New("h97jamjw3h")
	client.FlushAfter = 30 * time.Second
	client.FlushAt = 25

	for {
		log.Println("send track")

		client.Track(map[string]interface{}{
			"event":  "Download",
			"userId": "123456",
			"properties": map[string]interface{}{
				"application": "Segment Desktop",
				"version":     "1.1.0",
				"platform":    "osx",
			},
		})

		time.Sleep(50 * time.Millisecond)
	}
}
