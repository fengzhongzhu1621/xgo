package tests

import (
	"path/filepath"
	"testing"
)

// 获得文件的绝对路径
// AbsFilePath calls filepath.Abs on path.
func AbsFilePath(t *testing.T, path string) string {
	t.Helper()

	s, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}

	return s
}
