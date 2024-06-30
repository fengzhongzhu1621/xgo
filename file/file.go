package file

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 判断文件是否被修改.
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

// 判断文件是否存在.
func IsFileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

// 判断是否是目录.
func IsDirType(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断是否是文件.
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

// 读取go mod文件，返回模块名.
func GetGoModeName() (mod string, err error) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p := filepath.Join(d, "go.mod")
	_, err = os.Lstat(p)
	if err != nil {
		return
	}
	fin, err := os.Open(p)
	if err != nil {
		return
	}
	// 读文件
	sc := bufio.NewScanner(fin)
	// 逐行扫描，当扫描因为抵达输入流结尾或者遇到错误而停止时，本方法会返回false
	for sc.Scan() {
		// 获得行内容
		l := sc.Text()
		if strings.HasPrefix(l, "module ") {
			return strings.Split(l, " ")[1], nil
		}
	}
	return
}

// 实现 unix whtich 命令功能
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
