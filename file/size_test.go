package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Returns file size in bytes.
// func FileSize(path string) (int64, error)
func TestFileSize(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0o777)
	defer f.Close()

	f.WriteString("hello\nworld")

	size, err := fileutil.FileSize(fname)

	fmt.Println(size)
	fmt.Println(err)

	os.Remove(fname)

	// Output:
	// 11
	// <nil>
}
