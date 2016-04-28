package analytics

import (
	"net/http"
	"testing"
)

type cancelableRoundTripper struct {
	http.RoundTripper
}

func (c cancelableRoundTripper) CancelRequest(req *http.Request) {}

func TestHttpTransportIsCancelableTrue(t *testing.T) {
	if !httpTransportIsCancelable(cancelableRoundTripper{testTransportOK}) {
		t.Error("cancelable transport wasn't properly detected")
	}
}

func TestHttpTransportIsCancelableFalse(t *testing.T) {
	if httpTransportIsCancelable(testTransportOK) {
		t.Error("non-cancelable transport wasn't properly detected")
	}
}
