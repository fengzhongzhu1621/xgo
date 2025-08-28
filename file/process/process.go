package process

import (
	"fmt"
	"os"
)

const currentProcessFd = "/proc/self/fd"

// 获得当前进程的所有文件描述符名称.
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
		return nil, fmt.Errorf("could not read %s: %s", fd.Name(), err)
	}

	return names, nil
}

// 获得当前进程文件描述符的数量.
func GetCurrentProcessFdsLen() (int, error) {
	fdNames, err := GetCurrentProcessAllFdName()
	if err != nil {
		return 0, err
	}
	return len(fdNames), nil
}
