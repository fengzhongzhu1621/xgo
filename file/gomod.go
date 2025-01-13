package file

import (
	"bufio"
	"os"
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
