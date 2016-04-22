package main

import (
	"fmt"

	"github.com/segmentio/analytics-go"
)
import "time"

func main() {
	client := analytics.New("h97jamjwbh")
	client.Interval = 30 * time.Second
	client.Size = 100
	client.Verbose = true
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
				UserId: "123456",
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
