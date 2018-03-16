package main

import (
	"encoding/json"

	"github.com/segmentio/analytics-go"
	"github.com/segmentio/conf"
)

func main() {
	var config struct {
		WriteKey   string `conf:"writeKey"   help:"The Segment Write Key of the project to send data to"`
		Type       string `conf:"type"       help:"The type of the message to send"`
		UserID     string `conf:"userId"     help:"Unique identifier for the user"`
		GroupID    string `conf:"groupId"    help:"Unique identifier for the group"`
		Traits     string `conf:"traits"     help:"Metadata associated with the user"`
		Event      string `conf:"event"      help:"Name of the track event"`
		Properties string `conf:"properties" help:"Metadata associated with an event, page or screen call"`
		Name       string `conf:"name"       help:"Name of the page/screen"`
	}
	conf.Load(&config)

	callback := newCallback()

	client, err := analytics.NewWithConfig(config.WriteKey, analytics.Config{
		BatchSize: 1,
		Callback:  callback,
	})
	if err != nil {
		panic(err)
	}

	switch config.Type {
	case "track":
		client.Enqueue(analytics.Track{
			UserId:     config.UserID,
			Event:      config.Event,
			Properties: parseJSON(config.Properties),
		})
	case "identify":
		client.Enqueue(analytics.Identify{
			UserId: config.UserID,
			Traits: parseJSON(config.Traits),
		})
	case "group":
		client.Enqueue(analytics.Group{
			UserId:  config.UserID,
			GroupId: config.GroupID,
			Traits:  parseJSON(config.Traits),
		})
	case "page":
		client.Enqueue(analytics.Page{
			UserId:     config.UserID,
			Name:       config.Name,
			Properties: parseJSON(config.Properties),
		})
	case "screen":
		client.Enqueue(analytics.Screen{
			UserId:     config.UserID,
			Name:       config.Name,
			Properties: parseJSON(config.Properties),
		})
	}

	<-callback.success
}

// parseJSON parses a JSON formatted string into a map.
func parseJSON(v string) map[string]interface{} {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(v), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func newCallback() *callback {
	return &callback{
		success: make(chan struct{}, 1),
	}
}

// callback implements the analytics.Callback interface. It is used by the CLI
// to wait for events to be uploaded before exiting.
type callback struct {
	success chan struct{}
}

func (c *callback) Failure(_ analytics.Message, err error) {
	panic(err)
}

func (c *callback) Success(analytics.Message) {
	c.success <- struct{}{}
}
