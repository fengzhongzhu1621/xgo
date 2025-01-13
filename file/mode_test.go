package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// TestFileMode Return file mode infomation.
// func FileMode(path string) (fs.FileMode, error)
func TestFileMode(t *testing.T) {
	srcFile := "./text.txt"
	fileutil.CreateFile(srcFile)

	mode, _ := fileutil.FileMode(srcFile)

	os.Remove(srcFile)

	fmt.Println(mode)

	// Output:
	// -rw-r--r--
}
