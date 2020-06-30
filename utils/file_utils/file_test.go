package file_utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveFileExt(t *testing.T) {
	tests := []struct {
		filePath string
		filenameWithoutExt string
	} {
		{filePath: "On Unix:", filenameWithoutExt: "On Unix:"},
		{filePath: "/foo/bar/baz.js", filenameWithoutExt: "baz"},
		{filePath: "/foo/bar/baz", filenameWithoutExt: "baz"},
		{filePath: "/foo/bar/baz/", filenameWithoutExt: "baz"},
		{filePath: "dev.txt", filenameWithoutExt: "dev"},
		{filePath: "../todo.txt", filenameWithoutExt: "todo"},
		{filePath: "..", filenameWithoutExt: "."},
		{filePath: ".", filenameWithoutExt: ""},
		{filePath: "/", filenameWithoutExt: filepath.FromSlash("/")},
		{filePath: "", filenameWithoutExt: ""},
	}
	for _, test := range tests {
		actual := RemoveFileExt(test.filePath)
		assert.Equal(t, actual, test.filenameWithoutExt)
	}

}

func TestCopy(t *testing.T) {
	var dst string = "./file_test2.go"
	err := Copy("./file_test.go", dst)
	fmt.Println(err)
	assert.Equal(t, IsFileExists(dst), true)
	err = os.Remove(dst)
	assert.NoError(t, err)
}