package rollwriter

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/lestrrat-go/strftime"
)

const (
	backupTimeFormat = "bk-20060102-150405.00000"
	compressSuffix   = ".gz"
)

// RollWriter 是一个支持按大小或时间滚动的文件日志写入器
// 实现了 io.WriteCloser 接口
type RollWriter struct {
	filePath string // 文件路径
	opts     *Options

	pattern  *strftime.Strftime // 时间模式
	currPath string             // 当前文件路径
	currSize int64              // 当前文件大小
	currDir  string             // 当前目录
	currFile atomic.Value       // 当前打开的日志文件句柄
	openTime int64              // 上次打开文件的时间戳

	mu         sync.Mutex               // 互斥锁
	notifyOnce sync.Once                // 通知互斥锁
	notifyCh   chan bool                // 通知通道
	closeOnce  sync.Once                // 关闭通道的互斥锁
	closeCh    chan *closeAndRenameFile // 关闭文件句柄的通道

	os customizedOS // 操作系统接口
}

type closeAndRenameFile struct {
	file   *os.File
	rename string
}

type customizedOS interface {
	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
	MkdirAll(path string, perm fs.FileMode) error
	Rename(oldpath string, newpath string) error
	Stat(name string) (fs.FileInfo, error)
	Remove(name string) error
}

// NewRollWriter 创建一个新的RollWriter实例
// 参数：filePath - 文件路径，opt - 配置选项
// 返回值：RollWriter指针、错误信息
func NewRollWriter(filePath string, opt ...Option) (*RollWriter, error) {
	opts := &Options{
		MaxSize:    0,     // 默认不按文件大小滚动
		MaxAge:     0,     // 默认不清理过期日志
		MaxBackups: 0,     // 默认不清理冗余日志
		Compress:   false, // 默认不压缩
	}

	// opt具有最高优先级，会覆盖原始配置
	for _, o := range opt {
		o(opts)
	}

	if filePath == "" {
		return nil, errors.New("invalid file path")
	}

	pattern, err := strftime.New(filePath + opts.TimeFormat)
	if err != nil {
		return nil, errors.New("invalid time pattern")
	}

	w := &RollWriter{
		filePath: filePath,               // 文件路径
		opts:     opts,                   // 配置选项
		pattern:  pattern,                // 时间模式
		currDir:  filepath.Dir(filePath), // 当前目录
		os:       defaultCustomizedOS,    // 操作系统接口
	}

	// 创建日志目录
	if err := w.os.MkdirAll(w.currDir, 0755); err != nil { // 创建目录
		return nil, err
	}

	return w, nil
}

// getCurrFile 返回当前的日志文件
func (w *RollWriter) getCurrFile() *os.File {
	if file, ok := w.currFile.Load().(*os.File); ok {
		return file
	}
	return nil
}

// setCurrFile 设置当前的日志文件
func (w *RollWriter) setCurrFile(file *os.File) {
	w.currFile.Store(file)
}
