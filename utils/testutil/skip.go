package testutil

import (
	"runtime"
	"testing"
)

// SkipWindows Skip some tests on Windows that kept failing when Windows was added to the CI as a target.
func SkipWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip test on Windows")
	}
}
