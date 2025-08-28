package file

import (
	"os"
	"path/filepath"
	"strings"
)

// RemoveFileExt 根据路径获得文件名，并去掉文件名的后缀.
func RemoveFileExt(filePath string) string {
	// 根据路径获得文件名
	filename := filepath.Base(filePath)
	// 获得 . 出现的位置
	idx := strings.LastIndex(filename, ".")
	if idx < 0 {
		return filename
	}
	return filename[:idx]
}

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
