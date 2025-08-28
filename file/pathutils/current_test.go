package pathutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

func TestCurrentPath(t *testing.T) {
	absPath := fileutil.CurrentPath()
	fmt.Println(absPath)
}
