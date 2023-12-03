package assert

import (
	"testing"
)

func Equals[T comparable](t *testing.T, result, expected T) {
	// when t.Errorf() is called from our Equal() function, the Go test runner will report the filename and line
	// number of the code which called our Equal() function in the output.
	t.Helper()
	if result != expected {
		t.Errorf("\nResult %T; Expected %T", result, expected)
	}
}
