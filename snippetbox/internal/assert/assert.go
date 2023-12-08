package assert

import (
	"strings"
	"testing"
)

func Equals[T comparable](t *testing.T, result, expected T) {
	// when t.Errorf() is called from our Equal() function, the Go test runner will report the filename and line
	// number of the code which called our Equal() function in the output.
	t.Helper()
	if result != expected {
		t.Errorf("\nResult %v\nExpected %v", result, expected)
	}
}
func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}
func NilError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}
