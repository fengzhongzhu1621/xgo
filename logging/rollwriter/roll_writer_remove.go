package rollwriter

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// runCleanFiles 在新的goroutine中清理冗余或过期的（压缩的）日志文件
func (w *RollWriter) runCleanFiles() {
	for range w.notifyCh { // 监听清理通知通道
		if w.opts.MaxBackups == 0 && w.opts.MaxAge == 0 && !w.opts.Compress {
			// 如果不需要任何清理操作
			continue
		}
		// 异步清理文件
		w.cleanFiles()
	}
}

// cleanFiles cleans redundant or expired (compressed) logs.
func (w *RollWriter) cleanFiles() {
	// Get the file list of current log.
	files, err := w.getOldLogFiles()
	if err != nil {
		fmt.Printf("w.getOldLogFiles err: %+v\n", err)
		return
	}
	if len(files) == 0 {
		return
	}

	// Find the oldest files to scavenge.
	var compress, remove []logInfo
	files = filterByMaxBackups(files, &remove, w.opts.MaxBackups)

	// Find the expired files by last modified time.
	files = filterByMaxAge(files, &remove, w.opts.MaxAge)

	// Find files to compress by file extension .gz.
	filterByCompressExt(files, &compress, w.opts.Compress)

	// 删除过期或冗余文件
	w.removeFiles(remove)

	// 压缩日志文件
	w.compressFiles(compress)
}

// getOldLogFiles 返回按修改时间排序的日志文件列表
func (w *RollWriter) getOldLogFiles() ([]logInfo, error) {
	entries, err := os.ReadDir(w.currDir) // 读取当前目录的所有条目
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory %s :%w", w.currDir, err)
	}

	var logFiles []logInfo

	// 获取基础文件名
	filename := filepath.Base(w.filePath)
	for _, e := range entries {
		// 跳过目录
		if e.IsDir() {
			continue
		}

		// 跳过非日志文件，获得日志的最新修改时间
		if modTime, err := w.matchLogFile(e.Name(), filename); err == nil { // 匹配日志文件
			logFiles = append(logFiles, logInfo{modTime, e}) // 添加到日志文件列表
		}
	}

	// 按格式时间排序
	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// matchLogFile checks whether current log file matches all relative log files, if matched, returns
// the modified time.
func (w *RollWriter) matchLogFile(filename, filePrefix string) (time.Time, error) {
	// Exclude current log file.
	// a.log
	// a.log.20200712
	if filepath.Base(w.currPath) == filename {
		return time.Time{}, errors.New("ignore current logfile")
	}

	// Match all log files with current log file.
	// a.log -> a.log.20200712-1232/a.log.20200712-1232.gz
	// a.log.20200712 -> a.log.20200712.20200712-1232/a.log.20200712.20200712-1232.gz
	if !strings.HasPrefix(filename, filePrefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}

	st, err := w.os.Stat(filepath.Join(w.currDir, filename))
	if err != nil {
		return time.Time{}, fmt.Errorf("file stat fail: %w", err)
	}
	return st.ModTime(), nil
}

// removeFiles deletes expired or redundant log files.
func (w *RollWriter) removeFiles(remove []logInfo) {
	// Clean expired or redundant files.
	for _, f := range remove {
		file := filepath.Join(w.currDir, f.Name())
		if err := w.os.Remove(file); err != nil {
			fmt.Printf("remove file %s err: %+v\n", file, err)
		}
	}
}

// compressFiles compresses demanded log files.
func (w *RollWriter) compressFiles(compress []logInfo) {
	// Compress log files.
	for _, f := range compress {
		fn := filepath.Join(w.currDir, f.Name())
		w.compressFile(fn, fn+compressSuffix)
	}
}

// compressFile compresses file src to dst, and removes src on success.
func (w *RollWriter) compressFile(src, dst string) (err error) {
	f, err := w.os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	gzf, err := w.os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		f.Close()
		return fmt.Errorf("failed to open compressed file: %v", err)
	}

	gz := gzip.NewWriter(gzf)
	defer func() {
		gz.Close()
		// Make sure files are closed before removing, or else the removal
		// will fail on Windows.
		f.Close()
		gzf.Close()
		if err != nil {
			w.os.Remove(dst)
			err = fmt.Errorf("failed to compress file: %v", err)
			return
		}
		w.os.Remove(src)
	}()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	return nil
}
