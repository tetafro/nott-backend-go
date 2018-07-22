package testutils

import (
	"bytes"
	"encoding/json"
	"testing"
)

// AssertResponse compacts expected JSON response body and
// asserts it with the actual one.
func AssertResponse(t *testing.T, actual, expected string) {
	buf := new(bytes.Buffer)
	err := json.Compact(buf, []byte(expected))
	if err != nil {
		t.Fatalf("Fail to compact json: %v", err)
	}
	expectedBody := buf.String() + "\n"
	if actual != expectedBody {
		t.Fatalf("Expected response %s, but got %s", expectedBody, actual)
	}
}
