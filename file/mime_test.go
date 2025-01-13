package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Get file mime type, 'file' param's type should be string or *os.File.
// func MiMeType(file any) string
func TestMiMeType(t *testing.T) {
	fname := "./test.txt"
	file, _ := os.Create(fname)
	file.WriteString("hello world")

	f, _ := os.Open(fname)
	defer f.Close()

	mimeType := fileutil.MiMeType(f)
	fmt.Println(mimeType)

	os.Remove(fname)

	// Output:
	// application/octet-stream
}
