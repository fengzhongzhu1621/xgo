package file

import (
	"os"
	"time"
)

// IsFileModified 判断文件是否被修改.
func IsFileModified(filePath string, lastModifyTime time.Time) bool {
	baseFile, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if baseFile.ModTime().UnixNano() > lastModifyTime.UnixNano() {
		return true
	}
	return false
}

// FileExists 判断文件/文件夹是否存在（返回错误原因）.
func FileOrDirExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// IsFileExists 判断文件/文件夹是否存在（不返回错误）.
func IsFileOrDirExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDirType 判断是否是目录.
func IsDirType(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFileType 判断是否是文件.
func IsFileType(path string) bool {
	return !IsDirType(path)
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsDir()
}
