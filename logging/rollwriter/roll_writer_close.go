package rollwriter

import (
	"fmt"
	"time"
)

// Close closes the current log file. It implements io.Closer.
func (w *RollWriter) Close() error {
	if w.getCurrFile() == nil {
		return nil
	}
	err := w.getCurrFile().Close()
	w.setCurrFile(nil)

	if w.notifyCh != nil {
		close(w.notifyCh)
		w.notifyCh = nil
	}

	if w.closeCh != nil {
		close(w.closeCh)
		w.closeCh = nil
	}

	return err
}

// runCloseFiles delays closing file in a new goroutine.
func (w *RollWriter) runCloseFiles() {
	for f := range w.closeCh {
		time.Sleep(20 * time.Millisecond)

		// Close the file.
		if err := f.file.Close(); err != nil {
			fmt.Printf("f.file.Close err: %+v, filename: %s\n", err, f.file.Name())
		}
		if f.rename == "" || f.file.Name() == f.rename {
			continue
		}
		if err := w.os.Rename(f.file.Name(), f.rename); err != nil {
			fmt.Printf("os.Rename from %s to %s err: %+v\n", f.file.Name(), f.rename, err)
		}
		w.notify()
	}
}
