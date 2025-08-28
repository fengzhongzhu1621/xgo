//go:build windows

package file

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Get the version of exe or dll file on windows.
// func GetExeOrDllVersion(filePath string) (string, error)
func TestGetExeOrDllVersion(t *testing.T) {
	v, err := fileutil.GetExeOrDllVersion(`C:\Program Files\Tencent\WeChat\WeChat.exe`)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(v)

	// Output:
	// 3.9.10.19
}
