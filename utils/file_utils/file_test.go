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
		name string
		filePath string
		filenameWithoutExt string
	} {
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

func TestCopy(t *testing.T) {
	var dst string = "./file_test2.go"
	err := Copy("./file_test.go", dst)
	fmt.Println(err)
	assert.Equal(t, IsFileExists(dst), true)
	err = os.Remove(dst)
	assert.NoError(t, err)
}

func TestLocateFile(t *testing.T) {
	pwd, _ := os.Getwd()
	tests := []struct {
		name     string
		filename string
		dirs   []string
		expect string
	}{
		{name: "1", filename: "file.go", dirs: []string{pwd}, expect: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := LocateFile(tt.filename, tt.dirs)
			fmt.Println(actual)
		})
	}
}

