package file

import (
	"fmt"
	"os"
)

// CreateDir 创建目录.
func CreateDir(dirname string) error {
	var err error
	fin, err := os.Lstat(dirname)
	if err != nil {
		if os.IsExist(err) {
			return err
		}
		// 目录不存在则创建
		return os.MkdirAll(dirname, os.ModePerm)
	}

	if !fin.IsDir() {
		return fmt.Errorf("directory %s already exists", dirname)
	}

	return nil
}
