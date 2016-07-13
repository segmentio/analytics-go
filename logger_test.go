package analytics

import (
	"bytes"
	"errors"
	"log"
	"testing"
)

// This test ensures the standard logger shim to the Logger interface is working
// as expected.
func TestStdLogger(t *testing.T) {
	var buffer bytes.Buffer
	var logger = StdLogger(log.New(&buffer, "test ", 0), DebugLevel)

	logger.Debugf("Hello World!")
	logger.Infof("The answer is %d", 42)
	logger.Warnf("oops!")
	logger.Errorf("%s", errors.New("something went wrong!"))

	const ref = `test - DEBUG - Hello World!
test - INFO - The answer is 42
test - WARN - oops!
test - ERROR - something went wrong!
`

	if res := buffer.String(); ref != res {
		t.Errorf("invalid logs from standard logger:\n- expected: %s\n- found: %s", ref, res)
	}
}
