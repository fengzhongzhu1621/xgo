package validator

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/fengzhongzhu1621/xgo/crypto/randutils"
	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	tests := []struct {
		name string
		f    func()
	}{
		{
			name: "file exist",
			f: func() {
				// 创建临时文件
				f, err := os.CreateTemp("", "test")
				assert.NoError(t, err)
				defer os.Remove(f.Name())
				// 判读文件是否存在
				exist, err := FileOrDirExists(f.Name())
				assert.NoError(t, err)
				assert.True(t, exist)
			},
		},
		{
			name: "file not exist",
			f: func() {
				// 创建一个随机文件名，确保文件不存在
				fileName := randutils.RandomString(10)
				exist, err := FileOrDirExists(fileName)
				assert.NoError(t, err)
				assert.False(t, exist)
			},
		},
		{
			name: "is dir",
			f: func() {
				// 创建临时目录
				dirPath, err := os.MkdirTemp("", "example")
				assert.NoError(t, err)
				defer os.RemoveAll(dirPath) // clean up
				// 判断文件夹是否存在
				exist, err := FileOrDirExists(dirPath)
				assert.NoError(t, err)
				assert.True(t, exist)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f()
		})
	}
}

// Checks if a file or directory exists.
// func IsExist(path string) bool
func TestIsExist(t *testing.T) {
	result1 := fileutil.IsExist("./")
	result2 := fileutil.IsExist("./xxx.go")

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// true
	// false
}

// Checks if a file is symbol link or not.
// func IsLink(path string) bool
func TestIsLink(t *testing.T) {
	srcFile := "./test.txt"
	fileutil.CreateFile(srcFile)

	linkFile := "./test.link"
	os.Symlink(srcFile, linkFile)

	result1 := fileutil.IsLink(srcFile)
	result2 := fileutil.IsLink(linkFile)

	os.Remove(srcFile)
	os.Remove(linkFile)

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// false
	// true
}

// Checks if the path is directy or not.
// func IsDir(path string) bool
func TestIsDir(t *testing.T) {
	result1 := fileutil.IsDir("./")
	result2 := fileutil.IsDir("./xxx.go")

	fmt.Println(result1)
	fmt.Println(result2)

	// Output:
	// true
	// false
}
