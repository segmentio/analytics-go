package main

import (
	"fmt"

	"time"

	"github.com/segmentio/analytics-go"
)

func main() {
	client, _ := analytics.NewWithConfig("h97jamjwbh", analytics.Config{
		Interval:  30 * time.Second,
		BatchSize: 100,
		Verbose:   true,
	})
	defer client.Close()

	done := time.After(3 * time.Second)
	tick := time.Tick(50 * time.Millisecond)

	for {
		select {
		case <-done:
			fmt.Println("exiting")
			return

		case <-tick:
			if err := client.Enqueue(analytics.Message{
				"event":  "Download",
				"userId": "123456",
				"properties": map[string]interface{}{
					"application": "Segment Desktop",
					"version":     "1.1.0",
					"platform":    "osx",
				},
			}); err != nil {
				fmt.Println("error:", err)
				return
			}
		}
	}
}
