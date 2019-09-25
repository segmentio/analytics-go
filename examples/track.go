package main

import (
	"fmt"

	"github.com/rudderlabs/analytics-go"
)
import "time"

func main() {
	client, _ := analytics.NewWithConfig("h97jamjwbh", analytics.Config{
		Interval:  time.Second,
		BatchSize: 1,
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
			if err := client.Enqueue(analytics.Track{
				Event:  "Download",
				UserId: "444333",
				Properties: map[string]interface{}{
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
