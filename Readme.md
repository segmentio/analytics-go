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
}

func main() {
  client := analytics.New("your-writeKey-here")
  client.FlushInterval = 30 * time.Second
  client.BufferSize = 1000
  client.Debug = true

  for {
    client.Track("Download", Download{"segmentio", "1.0.0", "osx"})
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

#### func (*Client) Alias

```go
func (c *Client) Alias(previousId string) error
```

#### func (*Client) Group

```go
func (c *Client) Group(id string, traits interface{}) error
```

#### func (*Client) Identify

```go
func (c *Client) Identify(traits interface{}) error
```

#### func (*Client) Page

```go
func (c *Client) Page(name string, category string, properties interface{}) error
```

#### func (*Client) Screen

```go
func (c *Client) Screen(name string, category string, properties interface{}) error
```

#### func (*Client) Track

```go
func (c *Client) Track(event string, properties interface{}) error
```
