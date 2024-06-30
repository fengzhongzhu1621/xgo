package file

import (
	"path/filepath"
	"strings"
)

// 根据路径获得文件名，并去掉文件名的后缀.
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