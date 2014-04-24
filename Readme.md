# analytics
--
    import "github.com/segmentio/analytics-go"


## Usage

```go
const Version = "0.0.1"
```

#### type Client

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
