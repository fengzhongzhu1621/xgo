package rollwriter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/fengzhongzhu1621/xgo/opentelemetry/report"
	"github.com/hashicorp/go-multierror"
)

const (
	defaultLogQueueSize    = 10000
	defaultWriteLogSize    = 4 * 1024 // 4KB
	defaultLogIntervalInMs = 100
	defaultDropLog         = false
)

// AsyncRollWriter is the asynchronous rolling log writer which implements zapcore.WriteSyncer.
type AsyncRollWriter struct {
	logger io.WriteCloser
	opts   *AsyncOptions

	logQueue chan []byte
	sync     chan struct{}
	syncErr  chan error
	close    chan struct{}
	closeErr chan error
}

// NewAsyncRollWriter creates a new AsyncRollWriter.
func NewAsyncRollWriter(logger io.WriteCloser, opt ...AsyncOption) *AsyncRollWriter {
	opts := &AsyncOptions{
		LogQueueSize:     defaultLogQueueSize,
		WriteLogSize:     defaultWriteLogSize,
		WriteLogInterval: defaultLogIntervalInMs,
		DropLog:          defaultDropLog,
	}

	for _, o := range opt {
		o(opts)
	}

	w := &AsyncRollWriter{
		logger:   logger,
		opts:     opts,
		logQueue: make(chan []byte, opts.LogQueueSize),
		sync:     make(chan struct{}),
		syncErr:  make(chan error),
		close:    make(chan struct{}),
		closeErr: make(chan error),
	}

	// Start a new goroutine to write batch logs.
	go w.batchWriteLog()
	return w
}

// Write 写入日志数据，实现io.Writer接口
// 参数：data - 要写入的日志数据
// 返回值：写入的字节数、错误信息
func (w *AsyncRollWriter) Write(data []byte) (int, error) {
	log := make([]byte, len(data))
	copy(log, data)     // 复制数据以避免外部修改
	if w.opts.DropLog { // 如果启用了丢弃日志模式
		select {
		case w.logQueue <- log: // 尝试将日志放入队列
		default: // 如果队列已满
			report.LogQueueDropNum.Incr()
			return 0, errors.New("async roll writer: log queue is full") // 返回队列已满错误
		}
		return len(data), nil
	}

	w.logQueue <- log // 将日志放入队列（阻塞直到有空间）
	return len(data), nil
}

// Sync syncs logs. It implements zapcore.WriteSyncer.
func (w *AsyncRollWriter) Sync() error {
	w.sync <- struct{}{}
	return <-w.syncErr
}

// Close closes current log file. It implements io.Closer.
func (w *AsyncRollWriter) Close() error {
	err := w.Sync()
	close(w.close)
	return multierror.Append(err, <-w.closeErr).ErrorOrNil()
}

// batchWriteLog asynchronously writes logs in batches.
func (w *AsyncRollWriter) batchWriteLog() {
	buffer := bytes.NewBuffer(make([]byte, 0, w.opts.WriteLogSize*2))
	ticker := time.NewTicker(time.Millisecond * time.Duration(w.opts.WriteLogInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 定时刷新缓冲区中的数据到日志文件
			if buffer.Len() > 0 {
				_, err := w.logger.Write(buffer.Bytes())
				handleErr(err, "w.logger.Write on tick")
				buffer.Reset()
			}
		case data := <-w.logQueue:
			if len(data) >= w.opts.WriteLogSize {
				// If the length of the current data exceeds the expected maximum value,
				// we directly write it to the underlying logger instead of placing it into the buffer.
				// This prevents the buffer from being overwhelmed by excessively large data,
				// which could lead to memory leaks.
				// Prior to that, we need to write the existing data in the buffer to the underlying logger.
				_, _ = w.logger.Write(buffer.Bytes())
				buffer.Reset()
				_, _ = w.logger.Write(data)
				continue
			}

			// 如果缓冲区未满，继续读取队列中的数据，直到缓冲区满或队列为空
			// 读取队列中的数据，写入到缓冲区
			buffer.Write(data)

			// 如果缓冲区已满，将缓冲区中的数据写入到日志文件
			if buffer.Len() >= w.opts.WriteLogSize {
				_, err := w.logger.Write(buffer.Bytes())
				handleErr(err, "w.logger.Write on log queue")
				buffer.Reset()
			}
		case <-w.sync: // sync logs
			var err error
			if buffer.Len() > 0 {
				_, e := w.logger.Write(buffer.Bytes())
				err = multierror.Append(err, e).ErrorOrNil()
				buffer.Reset()
			}
			size := len(w.logQueue)
			for i := 0; i < size; i++ {
				v := <-w.logQueue
				_, e := w.logger.Write(v)
				err = multierror.Append(err, e).ErrorOrNil()
			}
			// 返回同步错误
			w.syncErr <- err
		case <-w.close: // 关闭日志写入器
			w.closeErr <- w.logger.Close()
			return
		}
	}
}

func handleErr(err error, msg string) {
	if err == nil {
		return
	}
	// Log writer has errors, so output to stdout directly.
	fmt.Printf("async roll writer err: %+v, msg: %s", err, msg)
}
