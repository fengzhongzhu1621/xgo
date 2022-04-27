// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package fsnotify

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"xgo/utils/channel"

	"golang.org/x/sys/unix"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events  chan Event
	Errors  chan error
	mu      sync.Mutex // Map access
	fd      int        // inotify fd
	poller  *fdPoller
	watches map[string]*watch // Map of inotify watches (key: path)
	// 文件描述符和路径对象的映射关系
	paths    map[int]string // Map of watched paths (key: watch descriptor)
	done     chan struct{}  // Channel for sending a "quit message" to the reader goroutine
	doneResp chan struct{}  // Channel to respond to Close
}

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
// 创建一个监听器
func NewWatcher() (*Watcher, error) {
	// Create inotify fd
	fd, errno := unix.InotifyInit1(unix.IN_CLOEXEC)
	if fd == -1 {
		return nil, errno
	}
	// Create epoll
	// 在epoll的fd上注册了两个文件, 一个是inotify的, 另一个是其用来实现优雅退出的pipe[0].
	poller, err := newFdPoller(fd)
	if err != nil {
		unix.Close(fd)
		return nil, err
	}
	w := &Watcher{
		fd:       fd,
		poller:   poller,
		watches:  make(map[string]*watch),
		paths:    make(map[int]string),
		Events:   make(chan Event),
		Errors:   make(chan error),
		done:     make(chan struct{}),
		doneResp: make(chan struct{}),
	}

	// 启动一个协程监听文件变更事件
	go w.readEvents()
	return w, nil
}

// 判断监听器是否已经关闭
func (w *Watcher) isClosed() bool {
	return channel.IsClosed(w.done)
}

// Close removes all watches and closes the events channel.
// 关闭监听器
func (w *Watcher) Close() error {
	if w.isClosed() {
		return nil
	}

	// Send 'close' signal to goroutine, and set the Watcher to closed.
	close(w.done)

	// Wake up goroutine
	// 使用管道发送关闭事件，读协程通过epoll可以获取到此事件
	w.poller.wake()

	// Wait for goroutine to close
	// 等待读事件结束
	<-w.doneResp

	return nil
}

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error {
	name = filepath.Clean(name)
	if w.isClosed() {
		return errors.New("inotify instance already closed")
	}

	const agnosticEvents = unix.IN_MOVED_TO | unix.IN_MOVED_FROM |
		unix.IN_CREATE | unix.IN_ATTRIB | unix.IN_MODIFY |
		unix.IN_MOVE_SELF | unix.IN_DELETE | unix.IN_DELETE_SELF

	var flags uint32 = agnosticEvents

	// 添加监听文件
	w.mu.Lock()
	defer w.mu.Unlock()
	watchEntry := w.watches[name]
	if watchEntry != nil {
		flags |= watchEntry.flags | unix.IN_MASK_ADD
	}
	wd, errno := unix.InotifyAddWatch(w.fd, name, flags) // wd是监听文件的句柄
	if wd == -1 {
		return errno
	}

	if watchEntry == nil {
		w.watches[name] = &watch{wd: uint32(wd), flags: flags}
		w.paths[wd] = name
	} else {
		watchEntry.wd = uint32(wd)
		watchEntry.flags = flags
	}

	return nil
}

// Remove stops watching the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error {
	name = filepath.Clean(name)

	// Fetch the watch.
	w.mu.Lock()
	defer w.mu.Unlock()
	watch, ok := w.watches[name]

	// Remove it from inotify.
	if !ok {
		return fmt.Errorf("can't remove non-existent inotify watch for: %s", name)
	}

	// We successfully removed the watch if InotifyRmWatch doesn't return an
	// error, we need to clean up our internal state to ensure it matches
	// inotify's kernel state.
	delete(w.paths, int(watch.wd))
	delete(w.watches, name)

	// 移除监听文件
	// inotify_rm_watch will return EINVAL if the file has been deleted;
	// the inotify will already have been removed.
	// watches and pathes are deleted in ignoreLinux() implicitly and asynchronously
	// by calling inotify_rm_watch() below. e.g. readEvents() goroutine receives IN_IGNORE
	// so that EINVAL means that the wd is being rm_watch()ed or its file removed
	// by another thread and we have not received IN_IGNORE event.
	success, errno := unix.InotifyRmWatch(w.fd, watch.wd)
	if success == -1 {
		// TODO: Perhaps it's not helpful to return an error here in every case.
		// the only two possible errors are:
		// EBADF, which happens when w.fd is not a valid file descriptor of any kind.
		// EINVAL, which is when fd is not an inotify descriptor or wd is not a valid watch descriptor.
		// Watch descriptors are invalidated when they are removed explicitly or implicitly;
		// explicitly by inotify_rm_watch, implicitly when the file they are watching is deleted.
		return errno
	}

	return nil
}

// WatchList returns the directories and files that are being monitered.
func (w *Watcher) WatchList() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	entries := make([]string, 0, len(w.watches))
	for pathname := range w.watches {
		entries = append(entries, pathname)
	}

	return entries
}

type watch struct {
	wd    uint32 // Watch descriptor (as returned by the inotify_add_watch() syscall)
	flags uint32 // inotify flags of this watch (see inotify(7) for the list of valid flags)
}

// readEvents reads from the inotify file descriptor, converts the
// received events into Event objects and sends them via the Events channel
func (w *Watcher) readEvents() {
	var (
		buf   [unix.SizeofInotifyEvent * 4096]byte // Buffer for a maximum of 4096 raw events
		n     int                                  // Number of bytes read with read()
		errno error                                // Syscall errno
		ok    bool                                 // For poller.wait
	)

	defer close(w.doneResp)
	defer close(w.Errors)
	defer close(w.Events)
	defer unix.Close(w.fd)
	defer w.poller.close()

	for {
		// See if we have been closed.
		if w.isClosed() {
			return
		}

		// 程序阻塞在这行, 直到epoll监听到相关事件为止
		ok, errno = w.poller.wait()
		if errno != nil {
			select {
			case w.Errors <- errno:
			case <-w.done: // 不过监听器关闭事件
				return
			}
			continue
		}

		if !ok {
			continue
		}

		// 读出事件到buffer里, 放到下面处理
		n, errno = unix.Read(w.fd, buf[:])
		// If a signal interrupted execution, see if we've been asked to close, and try again.
		// http://man7.org/linux/man-pages/man7/signal.7.html :
		// "Before Linux 3.8, reads from an inotify(7) file descriptor were not restartable"
		if errno == unix.EINTR {
			continue
		}

		// unix.Read might have been woken up by Close. If so, we're done.
		if w.isClosed() {
			return
		}

		// 当读到的事件小于16字节(一个事件结构体的单位大小), 异常处理逻辑
		if n < unix.SizeofInotifyEvent {
			var err error
			if n == 0 {
				// If EOF is received. This should really never happen.
				err = io.EOF
			} else if n < 0 {
				// If an error occurred while reading.
				err = errno
			} else {
				// Read was too short.
				err = errors.New("notify: short read in readEvents()")
			}
			select {
			case w.Errors <- err:
			case <-w.done:
				return
			}
			continue
		}

		var offset uint32
		// 此时我们也不知道读了几个事件到buffer里
		// 所以我们就用offset记录下当前所读到的位置偏移量, 直到读完为止
		// 这个for循环结束条件是: offset累加到了某个值, 以至于剩余字节数不够读取出一整个inotify event结构体
		// We don't know how many events we just read into the buffer
		// While the offset points to at least one whole event...
		for offset <= uint32(n-unix.SizeofInotifyEvent) {
			// Point "raw" to the event in the buffer
			// 强制把地址值转换成inotify结构体
			raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

			// 所发生的事件以掩码形式表示
			mask := uint32(raw.Mask)
			// 当监听的是个目录时, 目录中发生事件的文件名会包含在结构体中, 这里的len就是文件名的长度
			nameLen := uint32(raw.Len)

			// mask格式错误, 向Errors chan发送事件
			if mask&unix.IN_Q_OVERFLOW != 0 {
				select {
				case w.Errors <- ErrEventOverflow:
				case <-w.done:
					return
				}
			}

			// If the event happened to the watched directory or the watched file, the kernel
			// doesn't append the filename to the event, but we would like to always fill the
			// the "Name" field with a valid filename. We retrieve the path of the watch from
			// the "paths" map.
			w.mu.Lock()
			// 取出这个文件描述符所对应的文件名
			name, ok := w.paths[int(raw.Wd)]
			// IN_DELETE_SELF occurs when the file/directory being watched is removed.
			// This is a sign to clean up the maps, otherwise we are no longer in sync
			// with the inotify kernel state which has already deleted the watch
			// automatically.
			// 如果发生删除事件, 也一并在上下文中删掉这个文件名
			if ok && mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF {
				delete(w.paths, int(raw.Wd))
				delete(w.watches, name)
			}
			w.mu.Unlock()

			if nameLen > 0 {
				// 当watch是一个目录的时候, 其下面的文件发生事件时, 就会导致这个nameLen大于0
				// 此时读取文件名字(文件名就在inotify event结构体的后面), 强制把地址值转换成长度4096的byte数组
				// Point "bytes" at the first byte of the filename
				bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))[:nameLen:nameLen]
				// The filename is padded with NULL bytes. TrimRight() gets rid of those.
				// 拼接路径(文件名会以\000为结尾表示, 所以要去掉)
				name += "/" + strings.TrimRight(string(bytes[0:nameLen]), "\000")
			}

			// 生成一个event
			event := newEvent(name, mask)

			// Send the events that are not ignored on the events channel
			// 如果这个事件没有被忽略, 那么发送到Events chan
			if !event.ignoreLinux(mask) {
				select {
				case w.Events <- event: // 发送事件
				case <-w.done:
					return
				}
			}

			// 移动offset偏移量到下个inotify event结构体
			// Move to the next event in the buffer
			offset += unix.SizeofInotifyEvent + nameLen
		}
	}
}

// Certain types of events can be "ignored" and not sent over the Events
// channel. Such as events marked ignore by the kernel, or MODIFY events
// against files that do not exist.
func (e *Event) ignoreLinux(mask uint32) bool {
	// Ignore anything the inotify API says to ignore
	if mask&unix.IN_IGNORED == unix.IN_IGNORED {
		return true
	}

	// If the event is not a DELETE or RENAME, the file must exist.
	// Otherwise the event is ignored.
	// *Note*: this was put in place because it was seen that a MODIFY
	// event was sent after the DELETE. This ignores that MODIFY and
	// assumes a DELETE will come or has come if the file doesn't exist.
	if !(e.Op&Remove == Remove || e.Op&Rename == Rename) {
		_, statErr := os.Lstat(e.Name)
		return os.IsNotExist(statErr)
	}
	return false
}

// newEvent returns an platform-independent Event based on an inotify mask.
func newEvent(name string, mask uint32) Event {
	e := Event{Name: name}
	if mask&unix.IN_CREATE == unix.IN_CREATE || mask&unix.IN_MOVED_TO == unix.IN_MOVED_TO {
		e.Op |= Create
	}
	if mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF || mask&unix.IN_DELETE == unix.IN_DELETE {
		e.Op |= Remove
	}
	if mask&unix.IN_MODIFY == unix.IN_MODIFY {
		e.Op |= Write
	}
	if mask&unix.IN_MOVE_SELF == unix.IN_MOVE_SELF || mask&unix.IN_MOVED_FROM == unix.IN_MOVED_FROM {
		e.Op |= Rename
	}
	if mask&unix.IN_ATTRIB == unix.IN_ATTRIB {
		e.Op |= Chmod
	}
	return e
}
