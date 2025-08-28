package file

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/stretchr/testify/assert"
)

func TestRemoveFileExt(t *testing.T) {
	tests := []struct {
		name               string
		filePath           string
		filenameWithoutExt string
	}{
		{name: "1", filePath: "On Unix:", filenameWithoutExt: "On Unix:"},
		{name: "2", filePath: "/foo/bar/baz.js", filenameWithoutExt: "baz"},
		{name: "3", filePath: "/foo/bar/baz", filenameWithoutExt: "baz"},
		{name: "4", filePath: "/foo/bar/baz/", filenameWithoutExt: "baz"},
		{name: "5", filePath: "dev.txt", filenameWithoutExt: "dev"},
		{name: "6", filePath: "../todo.txt", filenameWithoutExt: "todo"},
		{name: "7", filePath: "..", filenameWithoutExt: "."},
		{name: "8", filePath: ".", filenameWithoutExt: ""},
		{name: "9", filePath: "/", filenameWithoutExt: filepath.FromSlash("/")},
		{name: "10", filePath: "", filenameWithoutExt: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := RemoveFileExt(tt.filePath)
			assert.Equal(t, actual, tt.filenameWithoutExt)
		})
	}
}

// Remove the file of path.
// func RemoveFile(path string) error
func TestRemoveFile(t *testing.T) {
	fname := "./test.txt"

	fileutil.CreateFile(fname)

	result1 := fileutil.IsExist(fname)

	fileutil.RemoveFile(fname)

	result2 := fileutil.IsExist(fname)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// true
	// false
}
