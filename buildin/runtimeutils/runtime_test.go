package runtimeutils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPackage(t *testing.T) {
	actual := GetPackage()
	expect := "xgo/buildin/runtimeutils"
	assert.Equal(t, expect, actual)
}

func TestGetCurrentFilename(t *testing.T) {
	// 获得当前文件路径
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	fmt.Println(dir)
}
