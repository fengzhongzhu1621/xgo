package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Returns file modified time(unix timestamp).
// func MTime(filepath string) (int64, error)
func TestMTime(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0777)
	defer f.Close()

	f.WriteString("hello\nworld")

	mtime, err := fileutil.MTime(fname)

	fmt.Println(mtime)
	fmt.Println(err)

	os.Remove(fname)

}
