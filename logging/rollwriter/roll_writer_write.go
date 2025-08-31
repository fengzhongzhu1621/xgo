package rollwriter

import (
	"errors"
	"sync/atomic"
	"time"
)

// Write writes logs. It implements io.Writer.
func (w *RollWriter) Write(v []byte) (n int, err error) {
	// Reopen file every 10 seconds.
	if w.getCurrFile() == nil || time.Now().Unix()-atomic.LoadInt64(&w.openTime) > 10 {
		w.mu.Lock()
		w.reopenFile()
		w.mu.Unlock()
	}

	// Return when failed to open the file.
	if w.getCurrFile() == nil {
		return 0, errors.New("open file fail")
	}

	// Write logs to file.
	n, err = w.getCurrFile().Write(v)
	atomic.AddInt64(&w.currSize, int64(n))

	// Rolling on full.
	if w.opts.MaxSize > 0 && atomic.LoadInt64(&w.currSize) >= w.opts.MaxSize {
		w.mu.Lock()
		w.backupFile()
		w.mu.Unlock()
	}

	return n, err
}
