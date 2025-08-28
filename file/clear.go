package file

import (
	"io"
	"os"
)

func ClearFile(f *os.File) error {
	// 将文件大小截断为 0
	err := f.Truncate(0)
	if err != nil {
		return err
	}

	// 重新设置文件指针
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}
