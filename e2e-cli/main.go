package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	analytics "github.com/segmentio/analytics-go/v3"
)

// Input JSON structures

type InputConfig struct {
	FlushAt       int `json:"flushAt"`
	FlushInterval int `json:"flushInterval"` // milliseconds
	MaxRetries    int `json:"maxRetries"`
	Timeout       int `json:"timeout"` // seconds
}

type InputEvent struct {
	Type        string                 `json:"type"`
	Event       string                 `json:"event"`
	UserId      string                 `json:"userId"`
	AnonymousId string                 `json:"anonymousId"`
	MessageId   string                 `json:"messageId"`
	PreviousId  string                 `json:"previousId"`
	GroupId     string                 `json:"groupId"`
	Name        string                 `json:"name"`
	Timestamp   string                 `json:"timestamp"`
	Properties  map[string]interface{} `json:"properties"`
	Traits      map[string]interface{} `json:"traits"`
	Integrations map[string]interface{} `json:"integrations"`
}

type InputSequence struct {
	DelayMs int          `json:"delayMs"`
	Events  []InputEvent `json:"events"`
}

type Input struct {
	WriteKey  string          `json:"writeKey"`
	ApiHost   string          `json:"apiHost"`
	Sequences []InputSequence `json:"sequences"`
	Config    InputConfig     `json:"config"`
}

// Output JSON structure

type Output struct {
	Success    bool   `json:"success"`
	SentBatches int   `json:"sentBatches"`
	Error      string `json:"error,omitempty"`
}

// Callback implementation

type myCallback struct {
	mu           sync.Mutex
	failures     []string
	successCount int
}

func (cb *myCallback) Success(msg analytics.Message) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.successCount++
	fmt.Fprintf(os.Stderr, "[callback] success: %T\n", msg)
}

func (cb *myCallback) Failure(msg analytics.Message, err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	errMsg := fmt.Sprintf("message %T failed: %v", msg, err)
	cb.failures = append(cb.failures, errMsg)
	fmt.Fprintf(os.Stderr, "[callback] failure: %s\n", errMsg)
}

// stderr logger that satisfies analytics.Logger

type stderrLogger struct {
	logger *log.Logger
}

func (l *stderrLogger) Logf(format string, args ...interface{}) {
	l.logger.Printf("INFO: "+format, args...)
}

func (l *stderrLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf("ERROR: "+format, args...)
}

func parseTimestamp(ts string) time.Time {
	if ts == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[warn] could not parse timestamp %q: %v\n", ts, err)
		return time.Time{}
	}
	return t
}

func enqueueEvent(client analytics.Client, ev InputEvent) error {
	ts := parseTimestamp(ev.Timestamp)
	integrations := analytics.Integrations(ev.Integrations)

	switch ev.Type {
	case "track":
		msg := analytics.Track{
			UserId:       ev.UserId,
			AnonymousId:  ev.AnonymousId,
			MessageId:    ev.MessageId,
			Event:        ev.Event,
			Properties:   analytics.Properties(ev.Properties),
			Integrations: integrations,
		}
		if !ts.IsZero() {
			msg.Timestamp = ts
		}
		return client.Enqueue(msg)

	case "identify":
		msg := analytics.Identify{
			UserId:       ev.UserId,
			AnonymousId:  ev.AnonymousId,
			MessageId:    ev.MessageId,
			Traits:       analytics.Traits(ev.Traits),
			Integrations: integrations,
		}
		if !ts.IsZero() {
			msg.Timestamp = ts
		}
		return client.Enqueue(msg)

	case "page":
		msg := analytics.Page{
			UserId:       ev.UserId,
			AnonymousId:  ev.AnonymousId,
			MessageId:    ev.MessageId,
			Name:         ev.Name,
			Properties:   analytics.Properties(ev.Properties),
			Integrations: integrations,
		}
		if !ts.IsZero() {
			msg.Timestamp = ts
		}
		return client.Enqueue(msg)

	case "screen":
		msg := analytics.Screen{
			UserId:       ev.UserId,
			AnonymousId:  ev.AnonymousId,
			MessageId:    ev.MessageId,
			Name:         ev.Name,
			Properties:   analytics.Properties(ev.Properties),
			Integrations: integrations,
		}
		if !ts.IsZero() {
			msg.Timestamp = ts
		}
		return client.Enqueue(msg)

	case "alias":
		msg := analytics.Alias{
			UserId:       ev.UserId,
			PreviousId:   ev.PreviousId,
			MessageId:    ev.MessageId,
			Integrations: integrations,
		}
		if !ts.IsZero() {
			msg.Timestamp = ts
		}
		return client.Enqueue(msg)

	case "group":
		msg := analytics.Group{
			UserId:       ev.UserId,
			AnonymousId:  ev.AnonymousId,
			MessageId:    ev.MessageId,
			GroupId:      ev.GroupId,
			Traits:       analytics.Traits(ev.Traits),
			Integrations: integrations,
		}
		if !ts.IsZero() {
			msg.Timestamp = ts
		}
		return client.Enqueue(msg)

	default:
		return fmt.Errorf("unknown event type: %q", ev.Type)
	}
}

func run(input Input) Output {
	cb := &myCallback{}

	cfg := analytics.Config{
		Callback: cb,
		Logger: &stderrLogger{
			logger: log.New(os.Stderr, "[analytics] ", log.LstdFlags),
		},
		Verbose: true,
	}

	if input.ApiHost != "" {
		cfg.Endpoint = input.ApiHost
	}

	if input.Config.FlushInterval > 0 {
		cfg.Interval = time.Duration(input.Config.FlushInterval) * time.Millisecond
	}

	if input.Config.FlushAt > 0 {
		cfg.BatchSize = input.Config.FlushAt
	}

	if input.Config.MaxRetries > 0 {
		cfg.MaxRetries = input.Config.MaxRetries
	}

	client, err := analytics.NewWithConfig(input.WriteKey, cfg)
	if err != nil {
		return Output{
			Success: false,
			Error:   fmt.Sprintf("failed to create analytics client: %v", err),
		}
	}

	var enqueueErrors []string

	for i, seq := range input.Sequences {
		if seq.DelayMs > 0 {
			fmt.Fprintf(os.Stderr, "[info] sequence %d: delaying %dms\n", i, seq.DelayMs)
			time.Sleep(time.Duration(seq.DelayMs) * time.Millisecond)
		}

		for j, ev := range seq.Events {
			fmt.Fprintf(os.Stderr, "[info] sequence %d, event %d: enqueuing %s\n", i, j, ev.Type)
			if err := enqueueEvent(client, ev); err != nil {
				msg := fmt.Sprintf("sequence %d, event %d (%s): enqueue error: %v", i, j, ev.Type, err)
				fmt.Fprintf(os.Stderr, "[error] %s\n", msg)
				enqueueErrors = append(enqueueErrors, msg)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "[info] closing client (flushing)...\n")
	if err := client.Close(); err != nil {
		enqueueErrors = append(enqueueErrors, fmt.Sprintf("close error: %v", err))
	}

	cb.mu.Lock()
	failures := cb.failures
	successCount := cb.successCount
	cb.mu.Unlock()

	allErrors := append(enqueueErrors, failures...)

	if len(allErrors) > 0 {
		return Output{
			Success:     false,
			SentBatches: countBatches(successCount, input.Config.FlushAt),
			Error:       fmt.Sprintf("%v", allErrors),
		}
	}

	return Output{
		Success:     true,
		SentBatches: countBatches(successCount, input.Config.FlushAt),
	}
}

// countBatches estimates the number of batches sent based on total successful
// messages and batch size. Returns at least 1 if any messages succeeded.
func countBatches(successCount, batchSize int) int {
	if successCount == 0 {
		return 0
	}
	if batchSize <= 0 {
		batchSize = analytics.DefaultBatchSize
	}
	batches := successCount / batchSize
	if successCount%batchSize != 0 {
		batches++
	}
	return batches
}

func main() {
	inputFlag := flag.String("input", "", "JSON input describing the events to send")
	flag.Parse()

	if *inputFlag == "" {
		fmt.Fprintf(os.Stderr, "error: --input flag is required\n")
		writeOutput(Output{Success: false, Error: "--input flag is required"})
		os.Exit(1)
	}

	var input Input
	if err := json.Unmarshal([]byte(*inputFlag), &input); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to parse --input JSON: %v\n", err)
		writeOutput(Output{Success: false, Error: fmt.Sprintf("failed to parse input JSON: %v", err)})
		os.Exit(1)
	}

	result := run(input)
	writeOutput(result)

	if !result.Success {
		os.Exit(1)
	}
}

func writeOutput(out Output) {
	b, err := json.Marshal(out)
	if err != nil {
		// Fallback if marshalling fails
		fmt.Printf(`{"success":false,"sentBatches":0,"error":"failed to marshal output: %v"}`, err)
		return
	}
	fmt.Println(string(b))
}
