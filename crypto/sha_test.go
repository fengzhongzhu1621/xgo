package crypto

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// returns file sha value, param `shaType` should be 1, 256 or 512.
// func Sha(filepath string, shaType ...int) (string, error)
func TestSha(t *testing.T) {
	fname := "./test.txt"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0777)
	defer f.Close()

	f.WriteString("hello\nworld")

	sha1, err := fileutil.Sha(fname, 1)
	sha256, _ := fileutil.Sha(fname, 256)
	sha512, _ := fileutil.Sha(fname, 512)

	fmt.Println(sha1)
	fmt.Println(sha256)
	fmt.Println(sha512)
	fmt.Println(err)

	os.Remove(fname)
}
