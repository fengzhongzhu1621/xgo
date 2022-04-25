package pathutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"xgo/utils/testutil"
)

// Skip some tests on Windows that kept failing when Windows was added to the CI as a target.
func skipWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip test on Windows")
	}
}

func TestAbsPathify(t *testing.T) {
	skipWindows(t)

	home := userHomeDir()
	homer := filepath.Join(home, "homer")
	wd, _ := os.Getwd()

	testutil.Setenv(t, "HOMER_ABSOLUTE_PATH", homer)
	testutil.Setenv(t, "VAR_WITH_RELATIVE_PATH", "relative")

	tests := []struct {
		input  string
		output string
	}{
		{"", wd},
		{"sub", filepath.Join(wd, "sub")},
		{"./", wd},
		{"./sub", filepath.Join(wd, "sub")},
		{"$HOME", home},
		{"$HOME/", home},
		{"$HOME/sub", filepath.Join(home, "sub")},
		{"$HOMER_ABSOLUTE_PATH", homer},
		{"$HOMER_ABSOLUTE_PATH/", homer},
		{"$HOMER_ABSOLUTE_PATH/sub", filepath.Join(homer, "sub")},
		{"$VAR_WITH_RELATIVE_PATH", filepath.Join(wd, "relative")},
		{"$VAR_WITH_RELATIVE_PATH/", filepath.Join(wd, "relative")},
		{"$VAR_WITH_RELATIVE_PATH/sub", filepath.Join(wd, "relative", "sub")},
	}

	for _, test := range tests {
		got := AbsPathify(test.input)
		if got != test.output {
			t.Errorf("Got %v\nexpected\n%q", got, test.output)
		}
	}
}
