package main

import (
	"github.com/rudderlabs/analytics-go"
)

func main() {
	// Instantiates a client to use send messages to the Rudder API.
	// User your WRITE KEY in below placeholder "RUDDER WRITE KEY"
	client := analytics.New("1aUR9IELHp6jqOW8HWkrYvMYHWy", "https://218da72a.ngrok.io")

	// Enqueues a track event that will be sent asynchronously.
	client.Enqueue(analytics.Track{
		UserId: "test-user",
		Event:  "test-snippet",
	})

	// Flushes any queued messages and closes the client.
	client.Close()
}
