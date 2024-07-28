package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// 从指定的多个目录中查询文件（查询深度为 1），返回查找到的文件的绝对路径
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

	// 获得文件的绝对路径
	absPath, err := filepath.Abs(filepaths[0])
	if err != nil {
		return "", err
	}
	return absPath, nil
}

// DeepPath 在目录下查询子目录，获得子目录的路径（如果子目录只有一个，则继续查询子目录）
func DeepPath(basedir, name string, maxDepth int) string {
	// loop max 5, incase of for loop not finished
	for depth := 0; depth <= maxDepth; depth += 1 {
		// 读取指定目录中的文件和子目录。它返回一个 []os.FileInfo 类型的切片，其中包含目录中文件和子目录的信息。
		finfos, err := os.ReadDir(filepath.Join(basedir, name))
		if err != nil || len(finfos) != 1 {
			break
		}
		if finfos[0].IsDir() {
			// 如果目录下只有一个子目录，则继续查询子目录
			name = filepath.ToSlash(filepath.Join(name, finfos[0].Name()))
		} else {
			break
		}
	}
	return name
}
