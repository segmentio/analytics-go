package analytics

import (
	"net/http"
	"reflect"
	"time"
)

// Makes an HTTP client using the transport given as argument and setting a
// default timeout.
// The request cancellation API changed between go 1.5 and 1.6, a new `Cancel`
// field was added to the http.Request type and the historical `CancelRequest`
// method of the transports was deprecated. The function takes care of setting
// the timeout only if it's gonna be supported by the transport (in go 1.5 it
// would generate an error to have a timeout if no `CancelRequest` method was
// available).
// The check is done at runtime because there are no way to do conditional
// compilation based on the go version.
func makeHttpClient(transport http.RoundTripper) (client http.Client) {
	if httpClientCanTimeout(transport) {
		client.Timeout = 10 * time.Second
	}
	client.Transport = transport
	return
}

func httpClientCanTimeout(transport http.RoundTripper) bool {
	return httpRequestIsCancelable() || httpTransportIsCancelable(transport)
}

func httpRequestIsCancelable() bool {
	// When this condition is true the go runtime is in version 1.6+ and a
	// timeout can always be set on the client.
	return reflect.ValueOf(http.Request{}).FieldByName("Cancel").IsValid()
}

func httpTransportIsCancelable(transport http.RoundTripper) bool {
	// When the runtime is in version 1.5.x or lower there is no `Cancel`
	// field on the request object and a timeout can only be set if the
	// transport has a `CancelRequest` method.
	_, ok := transport.(requestCanceler)
	return ok
}

type requestCanceler interface {
	CancelRequest(*http.Request)
}
