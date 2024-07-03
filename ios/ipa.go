package ios

import (
	"archive/zip"
	"errors"
	"io"
	"regexp"
)

// 读取压缩包中的图标文件
func ParseIpaIcon(path string) (data []byte, err error) {
	iconPattern := regexp.MustCompile(`(?i)^Payload/[^/]*/icon\.png$`)
	// 打开压缩包
	r, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer r.Close()

	// 编译压缩文件匹配图标文件
	var zfile *zip.File
	for _, file := range r.File {
		if iconPattern.MatchString(file.Name) {
			zfile = file
			break
		}
	}
	if zfile == nil {
		err = errors.New("icon.png file not found")
		return
	}

	// 打开图标文件
	plreader, err := zfile.Open()
	if err != nil {
		return
	}
	defer plreader.Close()
	// 用于读取 io.Reader 接口的所有数据到一个字节切片（[]byte）中。
	// 如果读取成功，函数会返回读取到的数据和 nil 错误；如果读取失败，则会返回一个非空的错误。
	return io.ReadAll(plreader)
}
