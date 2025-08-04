package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveUploadedFile uploads the form file to specific dst.
// 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目录
	if err = os.MkdirAll(filepath.Dir(dst), 0o750); err != nil {
		return err
	}

	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 写文件
	_, err = io.Copy(out, src)
	return err
}
