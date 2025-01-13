package file

// Create file in path. return true if create succeed.
//func CreateFile(path string) bool
import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

func TestCreateFile(t *testing.T) {
	isCreatedSucceed := fileutil.CreateFile("./test.txt")
	fmt.Println(isCreatedSucceed)
}

// Create directory in absolute path. param `absPath` like /a, /a/b.
// func CreateDir(absPath string) error
func TestCreateDir(t *testing.T) {
	err := fileutil.CreateDir("/a/b") // will create folder /a/b
	fmt.Println(err)
}
