package pathutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// List all file names in given path.
// func ListFileNames(path string) ([]string, error)
func TestListFileNames(t *testing.T) {
	fileutil.CreateFile("./file1.txt")
	fileutil.CreateFile("./file2.txt")

	fileList, _ := fileutil.ListFileNames("./")
	os.Remove("./file1.txt")
	os.Remove("./file2.txt")

	fmt.Println(fileList)
}
