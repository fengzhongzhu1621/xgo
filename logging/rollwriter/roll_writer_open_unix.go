//go:build !windows
// +build !windows

package rollwriter

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// doReopenFile reopens the file.
func (w *RollWriter) doReopenFile(path, _ string) error {
	// 记录文件打开时间
	atomic.StoreInt64(&w.openTime, time.Now().Unix())

	// 打开文件
	f, err := w.os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile %s err: %w", path, err)
	}
	// 设置最新的打开文件句柄
	last := w.getCurrFile()
	w.setCurrFile(f)

	if last != nil {
		// 延迟关闭旧的文件句柄
		w.delayCloseAndRenameFile(&closeAndRenameFile{file: last})
	}

	// 计算文件的大小
	st, err := w.os.Stat(path)
	if err != nil {
		return fmt.Errorf("os.Stat %s err: %w", path, err)
	}
	atomic.StoreInt64(&w.currSize, st.Size())

	return nil
}
