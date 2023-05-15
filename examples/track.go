package main

import (
	"fmt"

	journify "github.com/journifyio/journify-go-sdk"
)
import "time"

func main() {
	client, _ := journify.NewWithConfig("h97jamjwbh", journify.Config{
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
			if err := client.Enqueue(journify.Track{
				Event:  "Download",
				UserId: "123456",
				Properties: map[string]interface{}{
					"application": "Journify Desktop",
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
