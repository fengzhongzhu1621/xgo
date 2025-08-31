//go:build windows
// +build windows

package rollwriter

import (
	"fmt"
	"sync/atomic"
	"time"
)

// backupFile backs this file up and reopen a new one if file size is too large.
func (w *RollWriter) backupFile() {
	if !(w.opts.MaxSize > 0 && atomic.LoadInt64(&w.currSize) >= w.opts.MaxSize) {
		return
	}
	atomic.StoreInt64(&w.currSize, 0)
	backup := w.currPath + "." + time.Now().Format(backupTimeFormat)
	if err := w.doReopenFile(w.currPath, backup); err != nil {
		fmt.Printf("w.doReopenFile %s err: %+v\n", w.currPath, err)
	}
	w.notify()
}
