package zipfile

import (
	"archive/zip"
	"io"

	"github.com/fengzhongzhu1621/xgo/file"
)

type Zip struct {
	*zip.Writer
}

func (z *Zip) Add(relpath, abspath string) error {
	// 判断文件类型并返回一个 io.ReadCloser，用于后续读取文件内容
	info, rdc, err := file.StatFile(abspath)
	if err != nil {
		return err
	}
	defer rdc.Close()

	// 创建一个 zip.FileInfoHeader 实例
	/*
		zip.FileInfoHeader 结构体包含以下字段：

		Name：文件名（包括路径），以 UTF-8 编码的字符串。
		ModTime：文件的最后修改时间，使用 time.Time 类型表示。
		Mode：文件的权限和模式，使用 os.FileMode 类型表示。
		Size：文件的大小，以字节为单位，使用 int64 类型表示。
		Sys：文件的底层系统特定信息，通常为 *syscall.Stat_t 类型。
	*/
	hdr, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// 格式化文件名
	hdr.Name = file.SanitizedName(relpath)
	if info.IsDir() {
		hdr.Name += "/"
	}

	// 设置压缩方法为 Deflate
	hdr.Method = zip.Deflate // compress method

	// 将文件头添加到 ZIP 文件
	writer, err := z.CreateHeader(hdr)
	if err != nil {
		return err
	}

	// 将源文件内容复制到 ZIP 文件
	_, err = io.Copy(writer, rdc)

	return err
}
