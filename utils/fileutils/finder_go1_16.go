//go:build go1.16 && finder
// +build go1.16,finder

package fileutils

import (
	"errors"
	"fmt"
	"io/fs"
	"path"

	"github.com/spf13/afero"
)

type Finder struct {
	paths      []string // 文件的搜索路径
	fileNames  []string // 需要搜索的文件名，不包含后缀
	extensions []string // 需要搜索的文件后缀

	withoutExtension bool // 是否搜索不带后缀的文件，默认不搜索
}

// 搜索文件，只返回一个匹配的文件路径，如果没有找到，返回的路径为空字符串
func (f Finder) Find(fsys fs.FS) (string, error) {
	for _, searchPath := range f.paths {
		for _, fileName := range f.fileNames {
			for _, extension := range f.extensions {
				// 构造需要搜索的文件路径
				filePath := path.Join(searchPath, fileName+"."+extension)
				// 判读文件是否存在，文件不存在 err 为 nil，继续下个循环
				ok, err := FileExists(fsys, filePath)
				if err != nil {
					return "", err
				}

				if ok {
					return filePath, nil
				}
			}

			// 搜索不带后缀的文件
			if f.withoutExtension {
				filePath := path.Join(searchPath, fileName)

				ok, err := FileExists(fsys, filePath)
				if err != nil {
					return "", err
				}

				if ok {
					return filePath, nil
				}
			}
		}
	}

	return "", nil
}

// FileExists 判读文件是否存在.
func FileExists(fsys fs.FS, filePath string) (bool, error) {
	// 获得文件的基础信息
	fileInfo, err := fs.Stat(fsys, filePath)
	if err == nil {
		// 判断文件是否为目录
		return !fileInfo.IsDir(), nil
	}

	// 判断错误类型，文件不存在
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// FindConfigFile Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func FindConfigFile(fs afero.Fs, configPaths []string, configName string,
	supportedExts []string, configType string) (string, error) {
	finder := Finder{
		paths:            configPaths,
		fileNames:        []string{configName},
		extensions:       supportedExts,
		withoutExtension: configType != "",
	}

	file, err := finder.Find(afero.NewIOFS(fs))
	if err != nil {
		return "", err
	}

	if file == "" {
		return "", ConfigFileNotFoundError{configName, fmt.Sprintf("%s", configPaths)}
	}

	return file, nil
}
