package process

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fengzhongzhu1621/xgo/file"
	log "github.com/sirupsen/logrus"
)

var pidFile string

// init 包初始化时，设置 pidFile 的默认路径。
func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("get current path failed. Error:%s", err.Error())
	}
	// 存储 PID 文件的路径，初始化为当前工作目录下的 pid/<executable_name>.pid。
	pidFile = cwd + "/pid/" + filepath.Base(os.Args[0]) + ".pid"
}

// SavePid 保存 PID
func SavePid() error {
	if err := WritePid(); err != nil {
		return fmt.Errorf("write pid file failed. err:%s", err.Error())
	}

	return nil
}

// SetPidFilePath 允许外部设置 pidFile 的路径。
func SetPidFilePath(p string) {
	pidFile = p
}

// WritePad 写进程 ID 到文件中
func WritePid() error {
	if pidFile == "" {
		return fmt.Errorf("pidFile is not set")
	}

	// 创建 pid 文件目录
	if err := os.MkdirAll(filepath.Dir(pidFile), os.FileMode(0o755)); err != nil {
		return err
	}

	file, err := file.NewAtomicFile(pidFile, os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("error opening pidFile %s: %s", pidFile, err)
	}
	defer file.Close()

	// 写文件
	_, err = fmt.Fprintf(file, "%d", os.Getpid())
	if err != nil {
		return err
	}

	return nil
}

// ReadPid 读取 PID 文件获得进程ID
func ReadPid() (int, error) {
	if pidFile == "" {
		return 0, fmt.Errorf("pidFile is empty")
	}

	d, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))
	if err != nil {
		return 0, fmt.Errorf("error parsing pid from %s: %s", pidFile, err)
	}

	return pid, nil
}
