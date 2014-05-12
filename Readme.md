# analytics-go

  Segment.io analytics client for Go. For additional documentation
  visit [https://segment.io/docs/tracking-api/](https://segment.io/docs/tracking-api/).

## Installation

    $ go get github.com/segmentio/analytics-go

## Example

  Full example void of `client.Track` error-handling for brevity:

```go
package main

import "github.com/segmentio/analytics-go"
import "time"

func main() {
  client := analytics.New("h97jamjw3h")

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

With customized client:

```go
package main

import "github.com/segmentio/analytics-go"
import "time"

func main() {
  client := analytics.New("h97jamjw3h")
  client.FlushInterval = 5 * time.Second
  client.BufferSize = 20
  client.Debug = true

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

## API

```go
const Version = "0.0.1"
```

#### type Client

 By default messages are flushed in batches of __500__ or after
 the default flush interval of __30__ seconds.

```go
type Client struct {
	Debug         bool
	BufferSize    int
	FlushInterval time.Duration
	Endpoint      string
	Key           string
}
```

#### func  New

```go
func New(key string) (c *Client)
```

#### func (*Client) Alias

```go
func (c *Client) Alias(msg Message) error
```

#### func (*Client) Group

```go
func (c *Client) Group(msg Message) error
```

#### func (*Client) Identify

```go
func (c *Client) Identify(msg Message) error
```

#### func (*Client) Page

```go
func (c *Client) Page(msg Message) error
```

#### func (*Client) Screen

```go
func (c *Client) Screen(msg Message) error
```

#### func (*Client) Track

```go
func (c *Client) Track(msg Message) error
```

#### type Message

```go
type Message map[string]interface{}
```

## Debugging

 Enable `.Debug` to output verbose debugging info:

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
