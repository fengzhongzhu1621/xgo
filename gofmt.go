package xgo

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 格式化go代码.
func GoFmtFile(fpath string) error {
	// 读源码文件
	in, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	// 格式化源代码
	out, err := format.Source(in)
	if err != nil {
		return err
	}
	// 使用格式化后的代码替换
	err = ioutil.WriteFile(fpath, out, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// 格式化go代码目录.
func GoFmtDir(dir string) error {
	err := filepath.Walk(dir, func(fpath string, fileInfo os.FileInfo, err error) error {
		if strings.HasSuffix(fpath, ".go") && !fileInfo.IsDir() {
			err := GoFmtFile(fpath)
			if err != nil {
				return fmt.Errorf("gofmt file %s %v", fpath, err)
			}
		}
		return nil
	})
	return err
}
