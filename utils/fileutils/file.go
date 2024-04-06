package fileutils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const currentProcessFd = "/proc/self/fd"

// 删除所有的子目录.
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

// 获得当前进程的所有文件描述符名称.
func GetCurrentProcessAllFdName() ([]string, error) {
	// 打开当前进程的文件描述符
	fd, err := os.Open(currentProcessFd)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// 获得所有的目录名
	names, err := fd.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", fd.Name(), err)
	}

	return names, nil
}

// 获得当前进程文件描述符的数量.
func GetCurrentProcessFdsLen() (int, error) {
	fdNames, err := GetCurrentProcessAllFdName()
	if err != nil {
		return 0, err
	}
	return len(fdNames), nil
}

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

// 判断文件是否是符号链接.
func IsSymbolicLink(fileInfo os.FileInfo) bool {
	return fileInfo.Mode()&os.ModeSymlink != 0
}

// 复制符号链接.
func CopySymlink(src string, dest string) error {
	// 通过符号链接，能获取其所指向的路径名
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	// 复制符号链接
	return os.Symlink(src, dest)
}

func _copy(src, dest string, fileInfo os.FileInfo) error {
	// 如果源文件是符号链接
	if IsSymbolicLink(fileInfo) {
		return CopySymlink(src, dest)
	}
	// 如果源文件是目录
	if fileInfo.IsDir() {
		return CopyDir(src, dest, fileInfo)
	}
	// 如果源文件是普通文件
	return CopyFile(src, dest, fileInfo)
}

// 复制目录.
func CopyDir(src string, dst string, fileInfo os.FileInfo) error {
	if fileInfo == nil {
		fileInfo2, err := os.Lstat(src)
		if err != nil {
			return err
		}
		fileInfo = fileInfo2
	}
	// 创建目录文件夹
	fileModel := fileInfo.Mode()
	if err := os.MkdirAll(dst, fileModel); err != nil {
		return err
	}
	// 读取目录dirmane 中的所有目录和文件（不包括子目录）, 返回读取到的文件的信息列表
	contents, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, content := range contents {
		name := content.Name()
		subSrc := filepath.Join(src, name)
		subDst := filepath.Join(dst, name)
		// 递归调用
		if err := _copy(subSrc, subDst, content); err != nil {
			return err
		}
	}
	return nil
}

// 复制文件.
func CopyFile(src string, dest string, fileInfo os.FileInfo) error {
	// 如果没有 FileInfo，则获取源文件的 FileInfo
	if fileInfo == nil {
		fileInfo2, err := os.Lstat(src)
		if err != nil {
			return err
		}
		fileInfo = fileInfo2
	}
	// 递归创建目录文件夹
	fileModel := fileInfo.Mode()
	if err := os.MkdirAll(filepath.Dir(dest), fileModel); err != nil {
		return err
	}

	// 创建文件
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	// 设置文件的权限
	if err = os.Chmod(f.Name(), fileModel); err != nil {
		return err
	}

	// 打开源文件
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	// 复制文件
	_, err = io.Copy(f, s)

	return err
}

// 复制文件
// 注意Lstat和stat函数的区别，两个都是返回文件的状态信息
// Lstat多了处理Link文件的功能，会返回Linked文件的信息，而state直接返回的是Link文件所指向的文件的信息.
func Copy(src, dest string) error {
	fileInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return _copy(src, dest, fileInfo)
}

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
	path_list := strings.Split(envPath, ":")
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
			return filepath, nil
		}
	}
	return "", err
}
