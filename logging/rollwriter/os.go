package rollwriter

import (
	"io/fs"
	"os"
)

var defaultCustomizedOS = stdOS{}

type stdOS struct{}

func (stdOS) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (stdOS) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (stdOS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (stdOS) Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (stdOS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (stdOS) Remove(name string) error {
	return os.Remove(name)
}
