package rollwriter

import (
	"fmt"
	"sync/atomic"
	"time"
)

// reopenFile reopens the file regularly. It notifies the scavenger if file path has changed.
func (w *RollWriter) reopenFile() {
	if w.getCurrFile() == nil || time.Now().Unix()-atomic.LoadInt64(&w.openTime) > 10 {
		// 记录文件打开时间
		atomic.StoreInt64(&w.openTime, time.Now().Unix())
		// 获得最新的文件路径
		oldPath := w.currPath
		currPath := w.pattern.FormatString(time.Now())
		if w.currPath != currPath {
			w.currPath = currPath
			// 异步清理文件
			w.notify()
		}
		// 重新打开文件
		if err := w.doReopenFile(w.currPath, oldPath); err != nil {
			fmt.Printf("w.doReopenFile %s err: %+v\n", w.currPath, err)
		}
	}
}

// delayCloseAndRenameFile delays closing and renaming the given file.
// 延迟关闭旧的文件句柄
func (w *RollWriter) delayCloseAndRenameFile(f *closeAndRenameFile) {
	w.closeOnce.Do(func() {
		w.closeCh = make(chan *closeAndRenameFile, 100)
		go w.runCloseFiles()
	})
	w.closeCh <- f
}
