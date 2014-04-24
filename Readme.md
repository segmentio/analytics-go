
## analytics-go

 Segment.io analytics client for golang.

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
  client := analytics.Client("segmentio. writeKey here")

  for {
    client.Track("Download", Download{"segmentio", "1.0.0", "osx"})
    time.Sleep(50 * time.Millisecond)
  }
}
```