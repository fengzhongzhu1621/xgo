package pathutils

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/file/homedir"
	"github.com/fengzhongzhu1621/xgo/tests"
)

func TestAbsPathify(t *testing.T) {
	tests.SkipWindows(t)

	home := homedir.UserHomeDir()
	homer := filepath.Join(home, "homer")
	wd, _ := os.Getwd()

	tests.Setenv(t, "HOMER_ABSOLUTE_PATH", homer)
	tests.Setenv(t, "VAR_WITH_RELATIVE_PATH", "relative")

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
