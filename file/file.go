package file

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GetGoModeName 读取go mod文件，返回模块名.
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
