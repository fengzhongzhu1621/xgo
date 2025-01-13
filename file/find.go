package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fengzhongzhu1621/xgo/validator"
)

// LocateFile 从指定的多个目录中查询文件（查询深度为 1），返回查找到的文件的绝对路径
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
		if validator.IsFileType(filepath) {
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

// Which 实现 unix whtich 命令功能
func Which(cmd string) (filepath string, err error) {
	// 获得当前PATH环境变量
	envPath := os.Getenv("PATH")
	// 分割为多个路径
	path_list := strings.Split(envPath, string(os.PathListSeparator))
	for _, dirpath := range path_list {
		// 判断环境变量路径是否是目录
		dirInfo, err := os.Stat(dirpath)
		if err != nil {
			return "", err
		}
		if !dirInfo.IsDir() {
			continue
		}
		// 判断命令所在的路径是否存在
		filepath := path.Join(dirpath, cmd)
		_, err = os.Stat(filepath)
		if err == nil || os.IsExist(err) {
			return filepath, err
		}
	}
	return "", err
}
