package main

import "github.com/visionmedia/go-gracefully"
import "github.com/segmentio/analytics-go"
import "time"
import "log"

type Worker struct {
	analytics *analytics.Client
	exit      chan struct{}
}

func (w *Worker) Start() {
	println("starting")

	go func() {
		for {
			select {
			case <-w.exit:
				return
			case <-time.Tick(50 * time.Millisecond):
				log.Println("send track")

				w.analytics.Track(map[string]interface{}{
					"event":  "Download",
					"userId": "123456",
					"properties": map[string]interface{}{
						"application": "Segment Desktop",
						"version":     "1.1.0",
						"platform":    "osx",
					},
				})
			}
		}
	}()
}

func (w *Worker) Stop() {
	println("stopping")
	close(w.exit)
	println("flushing analytics")
	w.analytics.Stop()
	println("bye :)")
}

func NewWorker(client *analytics.Client) *Worker {
	return &Worker{
		analytics: client,
		exit:      make(chan struct{}),
	}
}

// run with DEBUG=analytics to view
// analytics-specific debug output
func main() {
	client := analytics.New("h97jamjw3h")
	client.FlushAfter = 5 * time.Second
	client.FlushAt = 25

	w := NewWorker(client)
	w.Start()
	gracefully.Shutdown()
	w.Stop()
}
