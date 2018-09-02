package testutils

import (
	"bytes"
	"encoding/json"
	"testing"
)

// AssertJSON compacts expected JSON and asserts it with the actual one.
func AssertJSON(t *testing.T, actual, expected string) {
	buf := new(bytes.Buffer)
	err := json.Compact(buf, []byte(expected))
	if err != nil {
		t.Fatalf("Fail to compact json: %v", err)
	}
	compacted := buf.String() + "\n"
	if actual != compacted {
		t.Fatalf("Expected response %s, but got %s", compacted, actual)
	}
}
