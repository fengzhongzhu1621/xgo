package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// RemoveContents 删除所有的子目录.
func RemoveContents(dir string) error {
	// 获得目录下所有的文件名
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

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
