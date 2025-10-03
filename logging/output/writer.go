package output

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fengzhongzhu1621/xgo/plugin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	writers = make(map[string]plugin.IFactory)
)

// RegisterWriter registers log output writer. Writer may have multiple implementations.
func RegisterWriter(name string, writer plugin.IFactory) {
	writers[name] = writer
}

// GetWriter gets log output writer, returns nil if not exist.
func GetLogWriter(name string) plugin.IFactory {
	return writers[name]
}

// NewSlogWriter 返回日志的输出
func NewSlogWriter(writerType string, settings map[string]string) (io.Writer, error) {
	switch writerType {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	case "file":
		return GetFileWriter(settings)
	}

	return nil, fmt.Errorf("[%s] writer not supported", writerType)
}

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
	// 获得日志写目录
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

	// 支持日志滚动和压缩功能
	writer := &lumberjack.Logger{
		// 日志文件的位置
		Filename: logPath,
		// 每个日志文件的大小限制（以MB为单位）
		MaxSize: size,
		// 保留的最大日志文件数量
		MaxBackups: backups,
		// 保留的最大日志文件天数
		MaxAge: age,
		// true 日志文件名将使用本地时间
		// false（默认值）时，日志文件名将使用 UTC 时间
		LocalTime: true,
		// 当日志文件滚动（即创建新的日志文件，并且旧的日志文件需要被归档或删除）时，旧的日志文件将被自动压缩。
		// 压缩通常使用常见的压缩算法，如gzip、bzip2等。lumberjack包默认使用gzip进行压缩。
		// 压缩后的文件通常会有一个.gz扩展名，例如app.log.1.gz。
		Compress: true,
	}

	return writer, nil
}
