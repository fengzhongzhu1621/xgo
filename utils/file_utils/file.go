package file_utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const currentProcessFd = "/proc/self/fd"


/**
 * 删除所有的子目录
 */
func RemoveContents(dir string) error {
	// 获得目录下所有的文件名
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * 获得当前进程的所有文件描述符名称
 */
func GetCurrentProcessAllFdName() ([]string, error) {
	// 打开当前进程的文件描述符
	fd, err := os.Open(currentProcessFd)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// 获得所有的目录名
	names, err := fd.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	return names, nil
}

/**
 * 获得当前进程文件描述符的数量
 */
func GetCurrentProcessFdsLen() (int, error) {
	fdNames, err := GetCurrentProcessAllFdName()
	if err != nil {
		return 0, err
	}
	return len(fdNames), nil
}


/**
 * 判断文件是否被修改
 */
func IsFileModified(filePath string, lastModifyTime time.Time) bool {
	baseFile, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if baseFile.ModTime().UnixNano() > lastModifyTime.UnixNano() {
		return true
	}
	return false
}
