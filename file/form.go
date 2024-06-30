package file

import (
	"io"
	"mime/multipart"
	"os"
)

// SaveUploadedFile uploads the form file to specific dst.
// 保存上传的文件 .
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
