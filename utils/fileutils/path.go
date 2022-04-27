package fileutils

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

type PathInfo struct {
	Name  string // 路径名称
	IsDir bool   // 是否是目录
}

// 获得最后一个字符
func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// JoinPaths 路径合并.
func JoinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

// PathExists 判断文件路径是否存在.
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetWd 获得应用程序当前路径.
func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

// GetHomeDir 获得当前用户的$HOME目录.
func GetHomeDir() string {
	home, _ := homedir.Dir()
	return home
}
