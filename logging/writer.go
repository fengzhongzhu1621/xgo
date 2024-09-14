package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

func GetWriter(writerType string, settings map[string]string) (io.Writer, error) {
	switch WriteType(writerType) {
	case WriterTypeOs:
		return GetOSWriter(settings)
	case WriterTypeFile:
		return GetFileWriter(settings)
	default:
		return GetOSWriter(map[string]string{"name": "stdout"})
	}
}

// GetOSWriter 获得终端输出 Writer
func GetOSWriter(settings map[string]string) (io.Writer, error) {
	switch settings["name"] {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return os.Stdout, nil
	}
}

// GetFileWriter 获得文件输出 Writer
func GetFileWriter(settings map[string]string) (io.Writer, error) {
	// 判断日志写目录
	path, ok := settings["path"]
	if ok {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("file path %s not exists", path)
		}
	} else {
		return nil, errors.New("log file path should not be empty")
	}

	// 获得日志文件名
	filename := settings["name"]

	// 日志备份文件数量
	backups := 10
	backupsStr, ok := settings["backups"]
	if ok {
		backupsInt, err := strconv.Atoi(backupsStr)
		if err != nil {
			return nil, errors.New("backups should be integer")
		}
		backups = backupsInt
	}

	// 日志文件大小
	size := 50
	sizeStr, ok := settings["size"]
	if ok {
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil {
			return nil, errors.New("size should be integer")
		}
		size = sizeInt
	}

	// 日志过期删除天数
	age := 7
	ageStr, ok := settings["age"]
	if ok {
		ageInt, err := strconv.Atoi(ageStr)
		if err != nil {
			return nil, errors.New("age should be integer")
		}
		age = ageInt
	}

	logPath := filename
	if path != "" {
		rawPath := strings.TrimSuffix(path, "/")
		logPath = filepath.Join(rawPath, filename)
	}

	writer := &lumberjack.Logger{
		Filename: logPath,
		// megabytes
		MaxSize:    size,
		MaxBackups: backups,
		// days
		MaxAge:    age,
		LocalTime: true,
	}

	return writer, nil
}
