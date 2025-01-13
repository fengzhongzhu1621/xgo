package file

import (
	"bytes"
	"io"
	"os"
)

// StatFile 判断文件类型并返回一个 io.ReadCloser，用于后续读取文件内容
func StatFile(filename string) (info os.FileInfo, reader io.ReadCloser, err error) {
	// 用于获取文件或目录的元信息（metadata），但不会跟随符号链接
	info, err = os.Lstat(filename)
	if err != nil {
		return info, nil, err
	}
	// 用于判断文件或目录是否为符号链接
	if info.Mode()&os.ModeSymlink != 0 {
		var target string
		// 读取符号链接指向的目标路径
		// 请注意，在使用 os.Readlink 之前，确保提供的路径确实是一个符号链接。否则，可能会遇到错误。
		// 可以使用 os.Lstat 函数和 FileInfo.Mode().IsSymlink() 方法来检查路径是否为符号链接。
		target, err = os.Readlink(filename)
		if err != nil {
			return info, nil, err
		}
		// 实现了 io.Closer 接口，但不会执行任何实际操作
		// 可以将任何实现了 io.Reader 接口的类型转换为 io.Closer 接口类型，而无需实际实现关闭操作。
		// bytes.NewBuffer 接受一个字节切片（[]byte）作为参数，并返回一个包含这些字节的缓冲区（*bytes.Buffer 类型）
		reader = io.NopCloser(bytes.NewBuffer([]byte(target)))
	} else if !info.IsDir() {
		// 打开文件
		reader, err = os.Open(filename)
		if err != nil {
			return info, reader, err
		}
	} else {
		reader = io.NopCloser(bytes.NewBuffer(nil))
	}

	return info, reader, err
}

// ReadFileContent 读取文件的内容
func ReadFileContent(filepath string) (string, error) {
	// 用于读取指定文件的全部内容，并将其作为一个字节切片返回。如果读取过程中发生错误，函数会返回一个非空的错误值。
	// os.ReadFile 函数会一次性读取整个文件内容到内存中，因此对于非常大的文件，可能会导致内存不足的问题。
	// 在这种情况下，建议使用 os.Open 和 io.Reader 接口逐块读取文件内容。
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
