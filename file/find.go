package file

import (
	"fmt"
	"path/filepath"
)

// 从指定的目录中查询文件，返回查找到的文件的绝对路径
// - 如果文件没有找到，返回err
// - 如果查找到多个文件，返回err.
func LocateFile(filename string, dirs []string) (string, error) {
	// 在当前目录下查询文件
	if len(dirs) == 0 || (len(dirs) == 1 && (dirs)[0] == ".") {
		abs, _ := filepath.Abs(".")
		return filepath.Join(abs, filename), nil
	}
	// 在指定目录下查询文件
	filepaths := []string{}
	for _, dir := range dirs {
		filepath := filepath.Join(dir, filename)
		if IsFileType(filepath) {
			filepaths = append(filepaths, filepath)
		}
	}
	// 判断查询到的文件的数量
	if len(filepaths) == 0 {
		return "", fmt.Errorf("%s not found in %v", filename, dirs)
	} else if len(filepaths) > 1 {
		return "", fmt.Errorf("%s was found in more than one directory: %v", filename, dirs)
	}

	absPath, err := filepath.Abs(filepaths[0])
	if err != nil {
		return "", err
	}
	return absPath, nil
}
