package file

import (
	"fmt"
	"os"
	"testing"
)

func TestStat(t *testing.T) {
	// 符号链接的路径
	symlinkPath := "stat_test.go"

	// 使用 os.Stat() 获取目标文件的信息
	statInfo, err := os.Stat(symlinkPath)
	if err != nil {
		fmt.Println("Error using os.Stat():", err)
	} else {
		fmt.Printf("os.Stat() - Target file info: %+v\n", statInfo)
	}

	// 使用 os.Lstat() 获取符号链接本身的信息
	lstatInfo, err := os.Lstat(symlinkPath)
	if err != nil {
		fmt.Println("Error using os.Lstat():", err)
	} else {
		fmt.Printf("os.Lstat() - Symlink info: %+v\n", lstatInfo)
	}
}
