# analytics-go

  Segment.io analytics client for Go. For additional documentation
  visit [https://segment.io/docs/libraries/go](https://segment.io/docs/libraries/go/) or view the [godocs](http://godoc.org/github.com/segmentio/analytics-go).

## Installation

    $ go get github.com/segmentio/analytics-go

## Examples

### Basic

  Full example void of `client.Track` error-handling for brevity:

```go
package main

import "github.com/segmentio/analytics-go"
import "time"

func main() {
  client := analytics.New("your-write-key-here")

  for {
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
```

### Client options

  Example with customized client:

```go
package main

import "github.com/segmentio/analytics-go"
import "time"

func main() {
  client := analytics.New("your-write-key-here")
  client.FlushAfter = 30 * time.Second
  client.FlushAt = 100

  for {
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
```

### Context

For each call a `context` map may be passed, which
is merged with the original values.

```go
client.Track(map[string]interface{}{
  "event":  "Download",
  "userId": "123456",
  "properties": map[string]interface{}{
    "application": "Segment Desktop",
    "version":     "1.1.0",
    "platform":    "osx",
  },
  "context": map[string]interface{}{
    "appVersion": "2.0.0",
    "appHostname": "some-host"
  }
})
```

### Flushing on shutdown

  The following example illustrates how `.Stop()`
  may be used to flush and wait for pending calls
  to be sent to Segment.

```go
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
  w.analytics.Close()
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
  client := analytics.New("your-write-key-here")
  client.FlushAfter = 5 * time.Second
  client.FlushAt = 25

  w := NewWorker(client)
  w.Start()
  gracefully.Shutdown()
  w.Stop()
}
```

## Debugging

 Enable debug output via the __DEBUG__ environment variable, for example `DEBUG=analytics`:

```
2014/04/23 18:56:57 buffer (110/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:57-0700}
2014/04/23 18:56:58 buffer (111/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (112/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (113/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (114/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (115/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (116/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (117/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
2014/04/23 18:56:58 buffer (118/1000) &{Track Download {segmentio 1.0.0 osx} 2014-04-23T18:56:58-0700}
```

## License

 MIT
