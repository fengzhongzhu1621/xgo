package tests

import (
	"os"
	"testing"
)

// 创建临时目录，返回临时目录的路径
// TempMkdir makes a temporary directory
// 需要调用 defer os.RemoveAll(testDirToMoveFiles) 进行删除.
func TempMkdir(t *testing.T, pattern string) string {
	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir
}

// 创建临时文件，返回临时文件的名称
// TempMkFile makes a temporary file.
// 需要调用 defer os.Remove(name) 删除临时文件.
func TempMkFile(t *testing.T, dir string, prefix string) string {
	f, err := os.CreateTemp(dir, prefix)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer f.Close()
	return f.Name()
}
