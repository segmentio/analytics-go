package analytics

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"testing"
)

type testLogger struct {
	logs   string
	errors string
}

func (t *testLogger) Logf(format string, args ...interface{}) {
	t.logs += fmt.Sprintf("INFO "+format+"\n", args...)
}

func (t *testLogger) Errorf(format string, args ...interface{}) {
	t.errors += fmt.Sprintf("ERROR "+format+"\n", args...)
}

// This test ensures that the interface doesn't get changed and stays compatible
// with the testLogger type.
// If someone were to modify the interface in backward incompatible manner this
// test would break.
func TestDummyLogger(t *testing.T) {
	var tester testLogger
	var logger Logger = &tester

	logger.Logf("Hello World!")
	logger.Logf("The answer is %d", 42)
	logger.Errorf("%s", errors.New("something went wrong!"))

	if tester.logs != "INFO Hello World!\nINFO The answer is 42\n" {
		t.Error("invalid logs:", tester.logs)
	}

	if tester.errors != "ERROR something went wrong!\n" {
		t.Error("invalid errors:", tester.errors)
	}
}

// This test ensures the standard logger shim to the Logger interface is working
// as expected.
func TestStdLogger(t *testing.T) {
	var buffer bytes.Buffer
	var logger = StdLogger(log.New(&buffer, "test ", 0))

	logger.Logf("Hello World!")
	logger.Logf("The answer is %d", 42)
	logger.Errorf("%s", errors.New("something went wrong!"))

	const ref = `test INFO Hello World!
test INFO The answer is 42
test ERROR something went wrong!
`

	if res := buffer.String(); ref != res {
		t.Errorf("invalid logs from standard logger:\n- expected: %s\n- found: %s", ref, res)
	}
}
