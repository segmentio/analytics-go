Journify client for Go.

## Installation

The package can be simply installed via go get, we recommend that you use a
package version management system like the Go vendor directory or a tool like
Godep to avoid issues related to API breaking changes introduced between major
versions of the library.

To install it in the GOPATH:
```
go get https://github.com/journifyio/journify-go-sdk
```

## Documentation

The links bellow should provide all the documentation needed to make the best
use of the library and the Journify API:

- [Documentation](https://segment.com/docs/libraries/go/)
- [godoc](https://godoc.org/gopkg.in/journifyio/journify-go-sdk)
- [API](https://segment.com/docs/libraries/http/)
- [Specs](https://segment.com/docs/spec/)

## Usage

```go
package main

import (
    "os"

    journify "github.com/journifyio/journify-go-sdk"
)

func main() {
    // Instantiates a client to use send messages to the Journify API.
    client := journify.New(os.Getenv("JOURNIFY_WRITE_KEY"))

    // Enqueues a track event that will be sent asynchronously.
    client.Enqueue(journify.Track{
        UserId: "test-user",
        Event:  "test-snippet",
    })

    // Flushes any queued messages and closes the client.
    client.Close()
}
```

## License

The library is released under the [MIT license](License.md).
