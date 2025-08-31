//go:build !windows
// +build !windows

package rollwriter

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// backupFile backs this file up and reopen a new one if file size is too large.
func (w *RollWriter) backupFile() {
	if !(w.opts.MaxSize > 0 && atomic.LoadInt64(&w.currSize) >= w.opts.MaxSize) {
		return
	}
	atomic.StoreInt64(&w.currSize, 0)

	// Rename the old file.
	backup := w.currPath + "." + time.Now().Format(backupTimeFormat)
	if _, err := w.os.Stat(w.currPath); !os.IsNotExist(err) {
		if err := w.os.Rename(w.currPath, backup); err != nil {
			fmt.Printf("os.Rename from %s to backup %s err: %+v\n", w.currPath, backup, err)
		}
	}

	// Reopen a new one.
	if err := w.doReopenFile(w.currPath, ""); err != nil {
		fmt.Printf("w.doReopenFile %s err: %+v\n", w.currPath, err)
	}

	// 异步清理文件
	w.notify()
}
