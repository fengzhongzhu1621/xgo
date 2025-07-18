package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/fengzhongzhu1621/xgo/validator"
)

// CopySymlink 复制符号链接.
func CopySymlink(src string, dest string) error {
	// 通过符号链接，能获取其所指向的路径名
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	// 复制符号链接
	return os.Symlink(src, dest)
}

// CopyDir 复制目录.
func CopyDir(src string, dst string, fileInfo os.FileInfo) error {
	if fileInfo == nil {
		fileInfo2, err := os.Lstat(src)
		if err != nil {
			return err
		}
		fileInfo = fileInfo2
	}
	// 创建目录文件夹
	fileModel := fileInfo.Mode()
	if err := os.MkdirAll(dst, fileModel); err != nil {
		return err
	}
	// 读取目录dirmane 中的所有目录和文件（不包括子目录）, 返回读取到的文件的信息列表
	contents, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, content := range contents {
		name := content.Name()
		subSrc := filepath.Join(src, name)
		subDst := filepath.Join(dst, name)
		// 递归调用
		if err := _copy(subSrc, subDst, content); err != nil {
			return err
		}
	}
	return nil
}

// CopyFile 复制文件.
func CopyFile(src string, dest string, fileInfo os.FileInfo) error {
	// 如果没有 FileInfo，则获取源文件的 FileInfo
	if fileInfo == nil {
		fileInfo2, err := os.Lstat(src)
		if err != nil {
			return err
		}
		fileInfo = fileInfo2
	}
	// 递归创建目录文件夹
	fileModel := fileInfo.Mode()
	if err := os.MkdirAll(filepath.Dir(dest), fileModel); err != nil {
		return err
	}

	// 创建文件
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	// 设置文件的权限
	if err = os.Chmod(f.Name(), fileModel); err != nil {
		return err
	}

	// 打开源文件
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	// 复制文件
	_, err = io.Copy(f, s)

	return err
}

// Copy 复制文件
// 注意Lstat和stat函数的区别，两个都是返回文件的状态信息
// Lstat多了处理Link文件的功能，会返回Linked文件的信息，而state直接返回的是Link文件所指向的文件的信息.
func Copy(src, dest string) error {
	fileInfo, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return _copy(src, dest, fileInfo)
}

func _copy(src, dest string, fileInfo os.FileInfo) error {
	// 如果源文件是符号链接
	if validator.IsSymbolicLink(fileInfo) {
		return CopySymlink(src, dest)
	}
	// 如果源文件是目录
	if fileInfo.IsDir() {
		return CopyDir(src, dest, fileInfo)
	}
	// 如果源文件是普通文件
	return CopyFile(src, dest, fileInfo)
}

// 定义缓冲池大小
const bufferSize = 32 * 1024 // 32KB

// 定义全局缓冲池
var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, bufferSize)
	},
}

// bufferedCopy 使用缓冲池中的缓冲区进行数据拷贝
func BufferedCopy(dst io.Writer, src io.Reader) (int64, error) {
	// 从缓冲池中获取缓冲区
	buf := bufPool.Get().([]byte)

	defer bufPool.Put(buf)

	return io.CopyBuffer(dst, src, buf)
}
