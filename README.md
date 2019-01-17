# analytics-go [![go-doc](https://godoc.org/github.com/FindHotel/analytics-go?status.svg)](https://godoc.org/github.com/FindHotel/analytics-go) [![Build Status](https://travis-ci.com/FindHotel/analytics-go.svg?branch=master)](https://travis-ci.com/FindHotel/analytics-go)

Segment analytics client for Go.

This is a fork of [segmentio/analytics-go](https://github.com/segmentio/analytics-go) library.

Notable changes:

- Use fully customised endpoint without `/v1/batch` postfix
- Use `X-API-Key` header for authorisation
- Ability to report usage metrics to DataDog

Latest stable branch is [master](https://github.com/FindHotel/analytics-go/tree/master). Releases are presented on [Releases page](https://github.com/FindHotel/analytics-go/releases).

## Installation

If you use [dep](https://github.com/golang/dep) then add these lines to your `Gopkg.toml`:

    [[constraint]]
      name = "github.com/FindHotel/analytics-go"
      version = "3.4.0"  # get the latest version from Releases

## Documentation

The links bellow should provide all the documentation needed to make the best
use of the library and the Segment API:

- [Documentation](https://segment.com/docs/libraries/go/)
- [godoc](https://godoc.org/github.com/FindHotel/analytics-go)
- [API](https://segment.com/docs/libraries/http/)
- [Specs](https://segment.com/docs/spec/)

## Usage

```go
package main

import (
    "os"

    analytics "github.com/FindHotel/analytics-go"
)

func main() {
    // Instantiates a client to use send messages to the segment API.
    client, err := analytics.NewWithConfig(
        os.Getenv("SEGMENT_WRITE_KEY"),
        analytics.Config{
            Endpoint: os.Getenv("SEGMENT_ENDPOINT"),
        },
    )
    if err != nil { // ALWAYS check for errors!
        panic(err)
    }

    // Enqueues a track event that will be sent asynchronously.
    client.Enqueue(analytics.Track{
        UserId: "test-user",
        Event:  "test-snippet",
    })

    // Flushes any queued messages and closes the client.
    client.Close()
}
```

## Reporting SDK metrics to DataDog

If you want to have SDK metrics (number of events succeeded, failed etc.)
to be delivered to DataDog use the following example:

```go
package main

import (
    "os"

    analytics "github.com/FindHotel/analytics-go"
)

func main() {
    // Instantiates a client to use send messages to the segment API.
    reporter := analytics.NewDatadogReporter(os.Getenv("DD_API_KEY"), os.Getenv("DD_APP_KEY"))
    // if you don't need metrics use
    // reporter := &analytics.DiscardReporter{}

    client, err := analytics.NewWithConfig(
        os.Getenv("SEGMENT_WRITE_KEY"),
        analytics.Config{
            Endpoint: os.Getenv("SEGMENT_ENDPOINT"),
            Reporters: []analytics.Reporter{reporter},
        },
    )
    if err != nil {
        panic(err)
    }

    // Enqueues a track event that will be sent asynchronously.
    client.Enqueue(analytics.Track{
        UserId: "test-user",
        Event:  "test-snippet",
    })

    // Flushes any queued messages and closes the client.
    client.Close()
}
```

## License

The library is released under the [MIT license](License.md).
