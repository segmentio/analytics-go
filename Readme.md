# analytics-go

  Segment.io analytics client for Go.

## Installation

    $ go get github.com/segmentio/analytics-go

## Example

```go
package main

import "github.com/segmentio/analytics-go"
import "time"

type Download struct {
  Application string `json:"application"`
  Version     string `json:"version"`
  Platform    string `json:"platform"`
  UserId      string `json:"userId"`
}

func main() {
  client := analytics.New("your-writeKey-here")
  client.FlushInterval = 30 * time.Second
  client.BufferSize = 1000
  client.Debug = true

  for {
    client.Track("Download", Download{"segmentio", "1.0.0", "osx", "some-id"})
    client.Identify(User{"tobi", "tobi@ferret.com"})
    time.Sleep(100 * time.Millisecond)
  }
}
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

#### func (*Client) Track

 You must pass `UserId` or `AnonymousId`.

```go
func (c *Client) Track(event string, properties interface{})
```

#### func (*Client) Group

You must pass `UserId` or `AnonymousId`.

```go
func (c *Client) Group(id string, traits interface{})
```

#### func (*Client) Identify

You must pass `UserId` or `AnonymousId`.

```go
func (c *Client) Identify(traits interface{})
```

#### func (*Client) Page

```go
func (c *Client) Page(name string, category string, properties interface{})
```

#### func (*Client) Screen

```go
func (c *Client) Screen(name string, category string, properties interface{})
```

#### func (*Client) Alias

```go
func (c *Client) Alias(previousId string)
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
