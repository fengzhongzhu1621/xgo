package file

// Clear the file content, write empty string to the file.
// func ClearFile(path string) error
import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

func TestClearFile(t *testing.T) {
	err := fileutil.ClearFile("./test.txt")
	if err != nil {
		fmt.Println(err)
	}
}
