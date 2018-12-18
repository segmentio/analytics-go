package analytics

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	datadog "github.com/zorkian/go-datadog-api"
)

var successCounters = newCounters("success")
var failureCounters = newCounters("failure")

// Reporter provides a function to reporting metrics.
type Reporter interface {
	Report(metricName string, value interface{}, tags []string, ts time.Time)
	AddTags(tags []string)
}

func splitTags(name string) (string, []string) {
	tokens := strings.Split(name, ".")
	if len(tokens) <= 1 {
		return name, []string{}
	}
	names := []string{}
	tags := []string{}
	for _, token := range tokens {
		if strings.Contains(token, ":") {
			tags = append(tags, token)
		} else {
			names = append(names, token)
		}
	}
	return strings.Join(names, "."), tags
}

func reportAll(prefix string, r Reporter) {
	ts := time.Now()
	metrics := metrics.DefaultRegistry.GetAll()
	go func() {
		for key, metric := range metrics {
			for measure, value := range metric {
				name, tags := splitTags(key)
				name = prefix + "." + name
				r.Report(name+"."+measure, value, tags, ts)
			}
		}
	}()
}

var hostname = func() string {
	h, err := os.Hostname()
	if err != nil {
		return "localhost"
	}
	return h
}()

// LogReporter report metrics as a log.
type LogReporter struct {
	Logger Logger
	tags   []string
}

// Report reports metrics.
func (r LogReporter) Report(metricName string, value interface{}, tags []string, ts time.Time) {
	allTags := append(tags, r.tags...)
	r.Logger.Logf("%s[%s] = %v", metricName, strings.Join(allTags, ", "), value)
}

// AddTags adds tags to be added to each metric reported.
func (r *LogReporter) AddTags(tags []string) {
	r.tags = append(r.tags, tags...)
}

// NewDatadogReporter is a factory method to create Datadog reporter
// with sane defaults.
func NewDatadogReporter(apiKey, appKey string) *DatadogReporter {
	dr := DatadogReporter{
		Client: datadog.NewClient(apiKey, appKey),
	}
	dr.Client.HttpClient = &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 30 * time.Second,
		},
	}
	dr.logger = newDefaultLogger()
	dr.tags = []string{"transport:http", "sdk:go", "version:" + Version}
	return &dr
}

// WithLogger sets logger to DatadogReporter.
func (dd *DatadogReporter) WithLogger(logger Logger) *DatadogReporter {
	dd.logger = logger
	return dd
}

// DatadogReporter reports metrics to DataDog.
type DatadogReporter struct {
	Client *datadog.Client
	logger Logger
	tags   []string
}

// AddTags adds tags to be added to each metric reported.
func (dd *DatadogReporter) AddTags(tags []string) {
	dd.tags = append(dd.tags, tags...)
}

// Report sends provided metric to Datadog.
func (dd DatadogReporter) Report(metricName string, value interface{}, tags []string, ts time.Time) {
	metricType := "gauge"
	metricValue, err := func() (float64, error) {
		switch v := value.(type) {
		case float64:
			return v, nil
		case int64:
			return float64(v), nil
		case int:
			return float64(v), nil
		}
		return 0, fmt.Errorf("can't handle value %+v", value)
	}()
	if err != nil {
		dd.logger.Errorf("Serializing value for metric %s(%+v) failed: %s", metricName, value, err)
		return
	}
	metricTimestamp := float64(ts.Truncate(time.Minute).Unix())
	allTags := append(tags, "hostname:"+hostname)
	allTags = append(allTags, dd.tags...)
	metric := datadog.Metric{
		Metric: &metricName,
		Type:   &metricType,
		Tags:   allTags,
		Points: []datadog.DataPoint{{&metricTimestamp, &metricValue}},
	}

	if err := dd.Client.PostMetrics([]datadog.Metric{metric}); err != nil {
		dd.logger.Errorf("Reporting metrics failed: %s", err)
	}
}

func resetMetrics() {
	registry := metrics.DefaultRegistry
	for name := range registry.GetAll() {
		metric := registry.Get(name)
		switch m := metric.(type) {
		case metrics.Counter:
			m.Clear()
		case metrics.Gauge:
			m.Update(0)
		case metrics.Histogram:
			// do nothing as Histogram has it's own internal cleanup
		}
	}
}

// newCounters returns factory for tagged counters.
func newCounters(name string) func(tags ...string) metrics.Counter {
	counters := make(map[string]metrics.Counter)
	mu := &sync.Mutex{}

	return func(tags ...string) metrics.Counter {
		fullName := strings.Join(append([]string{name}, tags...), ".")

		mu.Lock()
		defer mu.Unlock()

		counter, ok := counters[fullName]
		if !ok {
			counter = metrics.GetOrRegister(
				fullName,
				metrics.NewCounter(),
			).(metrics.Counter)
			counters[fullName] = counter
		}
		return counter
	}
}

func (c *client) loopMetrics() {
	var reporter = c.Config.Reporter
	reporter.AddTags([]string{
		"key:" + c.key,
		"endpoint:" + c.Config.Endpoint,
	})

	for range time.Tick(60 * time.Second) {
		reportAll("evas.submitted", reporter)
		resetMetrics()
	}
}
