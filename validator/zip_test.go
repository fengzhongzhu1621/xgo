package validator

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Checks if file is zip file or not.
// func IsZipFile(filepath string) bool
func TestIsZipFile(t *testing.T) {
	isZip := fileutil.IsZipFile("./zipfile.zip")
	fmt.Println(isZip)
}
