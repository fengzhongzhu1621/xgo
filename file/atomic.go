package file

import (
	"os"
	"path/filepath"
)

// AtomicFile behaves like os.File, but does an atomic rename operation at Close.
// 一个基本的原子文件写入机制，通过创建临时文件并在写入完成后重命名为目标文件路径，确保了文件写入的原子性。
type AtomicFile struct {
	*os.File        // 嵌入了 *os.File，因此它拥有 os.File 的所有方法和属性。
	path     string // 存储目标文件的路径，即最终替换的文件位置。
}

// NewAtomicFile 创建临时文件
func NewAtomicFile(path string, mode os.FileMode) (*AtomicFile, error) {
	// 在目标文件所在目录创建一个临时文件。临时文件的名称基于目标文件的基名。
	tmpDir := filepath.Dir(path)
	baseName := filepath.Base(path)
	f, err := os.CreateTemp(tmpDir, baseName+".tmp-*")
	if err != nil {
		return nil, err
	}
	// 设置临时文件的权限为指定的 mode
	if err := os.Chmod(f.Name(), mode); err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		// 设置权限失败，尝试关闭并删除临时文件，然后返回错误。
		if removeErr := os.Remove(f.Name()); removeErr != nil {
			return nil, removeErr
		}
		return nil, err
	}
	return &AtomicFile{File: f, path: path}, nil
}

// Close 关闭并替换文件
func (f *AtomicFile) Close() error {
	// 关闭临时文件
	if err := f.File.Close(); err != nil {
		// 关闭失败，尝试删除临时文件并返回错误
		if removeErr := os.Remove(f.Name()); removeErr != nil {
			return removeErr
		}
		return err
	}

	// 成功关闭后，将临时文件重命名为目标路径，实现原子替换。 f.Name() -> f.path
	// 用于重命名或移动文件和目录。它提供了一种简便的方法来更改文件或目录的名称，或者将其移动到不同的位置。
	if err := os.Rename(f.Name(), f.path); err != nil {
		return err
	}
	return nil
}

// Abort 中止写入并删除临时文件
func (f *AtomicFile) Abort() error {
	// 关闭临时文件
	if err := f.File.Close(); err != nil {
		if removeErr := os.Remove(f.Name()); removeErr != nil {
			return removeErr
		}
		return err
	}

	// 删除临时文件，确保不留下未完成的文件
	if err := os.Remove(f.Name()); err != nil {
		return err
	}
	return nil
}
