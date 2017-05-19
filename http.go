package analytics

import (
	"io/ioutil"
	"net/http"
	"strings"
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
	client := http.Client{
		Transport: roundTripper{},
		Timeout:   10 * time.Second,
	}

	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	_, err := client.Do(req)
	_, ok := transport.(requestCanceler)

	return err == nil || ok
}

type requestCanceler interface {
	CancelRequest(*http.Request)
}

type roundTripper struct{}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Proto:      r.Proto,
		ProtoMajor: r.ProtoMajor,
		ProtoMinor: r.ProtoMinor,
		Body:       ioutil.NopCloser(strings.NewReader("")),
		Request:    r,
	}, nil
}
