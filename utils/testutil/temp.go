package testutil

import (
	"io/ioutil"
	"testing"
)

// TempMkdir makes a temporary directory
// 需要调用 defer os.RemoveAll(testDirToMoveFiles) 进行删除
func TempMkdir(t *testing.T, pattern string) string {
	dir, err := ioutil.TempDir("", pattern)
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir
}

// TempMkFile makes a temporary file.
// 需要调用 defer os.Remove(name) 删除临时文件
func TempMkFile(t *testing.T, dir string, prefix string) string {
	f, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer f.Close()
	return f.Name()
}
