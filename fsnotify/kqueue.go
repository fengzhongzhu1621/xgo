// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || openbsd || netbsd || dragonfly || darwin
// +build freebsd openbsd netbsd dragonfly darwin

package fsnotify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/datetime"
	"github.com/fengzhongzhu1621/xgo/file/pathutils"

	"golang.org/x/sys/unix"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events chan Event
	Errors chan error
	done   chan struct{} // Channel for sending a "quit message" to the reader goroutine

	kq int // File descriptor (as returned by the kqueue() syscall).

	mu              sync.Mutex      // Protects access to watcher data
	watches         map[string]int  // Map of watched file descriptors (key: path). 路径和文件描述符的映射关系
	externalWatches map[string]bool // Map of watches added by user of the library.
	// Map of watched directories to fflags used in kqueue.
	// 监听的目录和监听的事件的映射关系
	dirFlags map[string]uint32
	// Map file descriptors to path names for processing kqueue events. 文件描述符和路径对象的映射关系
	paths map[int]pathutils.PathInfo
	// Keep track of if we know this file exists (to stop duplicate create events).
	fileExists map[string]bool
	isClosed   bool // Set to true when Close() is first called
}

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) {
	// creates a new kernel event queue and returns a descriptor.
	kq, err := kqueue()
	if err != nil {
		return nil, err
	}
	// 创建监听器对象
	w := &Watcher{
		kq:              kq,
		watches:         make(map[string]int), // 文件描述符
		dirFlags:        make(map[string]uint32),
		paths:           make(map[int]pathutils.PathInfo),
		fileExists:      make(map[string]bool),
		externalWatches: make(map[string]bool),
		Events:          make(chan Event),
		Errors:          make(chan error),
		done:            make(chan struct{}),
	}

	go w.readEvents()
	return w, nil
}

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error {
	// 防止重复关闭
	w.mu.Lock()
	if w.isClosed {
		w.mu.Unlock()
		return nil
	}
	// 标记已经关闭
	w.isClosed = true

	// copy paths to remove while locked
	// 获得监听的文件描述符
	var pathsToRemove = make([]string, 0, len(w.watches))
	for name := range w.watches {
		pathsToRemove = append(pathsToRemove, name)
	}
	w.mu.Unlock()
	// unlock before calling Remove, which also locks

	// 删除监听的文件描述符
	for _, name := range pathsToRemove {
		w.Remove(name)
	}

	// send a "quit" message to the reader goroutine
	// 发送关闭命令
	close(w.done)

	return nil
}

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error {
	w.mu.Lock()
	w.externalWatches[name] = true
	w.mu.Unlock()
	_, err := w.addWatch(name, noteAllEvents)
	return err
}

// Remove stops watching the the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error {
	name = filepath.Clean(name)
	// 获得监听的文件或目录的描述符
	w.mu.Lock()
	watchfd, ok := w.watches[name]
	w.mu.Unlock()
	if !ok {
		return fmt.Errorf("can't remove non-existent kevent watch for: %s", name)
	}

	// 内核取消监听事件
	const registerRemove = unix.EV_DELETE
	if err := register(w.kq, []int{watchfd}, registerRemove, 0); err != nil {
		return err
	}

	// 关闭需要监听的文件或目录
	unix.Close(watchfd)

	w.mu.Lock()
	isDir := w.paths[watchfd].IsDir
	delete(w.watches, name)
	delete(w.paths, watchfd)
	delete(w.dirFlags, name)
	w.mu.Unlock()

	// Find all watched paths that are in this directory that are not external.
	if isDir {
		// 获得目录下监听的子目录或文件
		var pathsToRemove []string
		w.mu.Lock()
		for _, path := range w.paths {
			wdir, _ := filepath.Split(path.Name)
			if filepath.Clean(wdir) == name {
				if !w.externalWatches[path.Name] {
					pathsToRemove = append(pathsToRemove, path.Name)
				}
			}
		}
		w.mu.Unlock()
		// 删除目录下监听的子目录或文件的监听事件
		for _, name := range pathsToRemove {
			// Since these are internal, not much sense in propagating error
			// to the user, as that will just confuse them with an error about
			// a path they did not explicitly watch themselves.
			w.Remove(name)
		}
	}

	return nil
}

// WatchList returns the directories and files that are being monitered.
// 获得监听的文件或目录列表.
func (w *Watcher) WatchList() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	entries := make([]string, 0, len(w.watches))
	for pathname := range w.watches {
		entries = append(entries, pathname)
	}

	return entries
}

// Watch all events (except NOTE_EXTEND, NOTE_LINK, NOTE_REVOKE)
// 默认监听的事件.
const noteAllEvents = unix.NOTE_DELETE | unix.NOTE_WRITE | unix.NOTE_ATTRIB | unix.NOTE_RENAME

// keventWaitTime to block on each read from kevent.
var keventWaitTime = datetime.DurationToTimespec(100 * time.Millisecond)

// addWatch adds name to the watched file set.
// The flags are interpreted as described in kevent(2).
// Returns the real path to the file which was added, if any, which may be different from the one passed in the case of symlinks.
func (w *Watcher) addWatch(name string, flags uint32) (string, error) {
	var isDir bool
	// Make ./name and name equivalent
	name = filepath.Clean(name)

	w.mu.Lock()
	if w.isClosed {
		w.mu.Unlock()
		return "", errors.New("kevent instance already closed")
	}

	watchfd, alreadyWatching := w.watches[name]
	// We already have a watch, but we can still override flags.
	if alreadyWatching {
		isDir = w.paths[watchfd].IsDir
	}
	w.mu.Unlock()

	if !alreadyWatching {
		fi, err := os.Lstat(name)
		if err != nil {
			return "", err
		}

		// Don't watch sockets.
		if fi.Mode()&os.ModeSocket == os.ModeSocket {
			return "", nil
		}

		// Don't watch named pipes.
		if fi.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
			return "", nil
		}

		// Follow Symlinks
		// Unfortunately, Linux can add bogus symlinks to watch list without
		// issue, and Windows can't do symlinks period (AFAIK). To  maintain
		// consistency, we will act like everything is fine. There will simply
		// be no file events for broken symlinks.
		// Hence the returns of nil on errors.
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			name, err = filepath.EvalSymlinks(name)
			if err != nil {
				return "", nil
			}

			w.mu.Lock()
			_, alreadyWatching = w.watches[name]
			w.mu.Unlock()

			if alreadyWatching {
				return name, nil
			}

			fi, err = os.Lstat(name)
			if err != nil {
				return "", nil
			}
		}

		// 打开文件
		watchfd, err = unix.Open(name, openMode, 0700)
		if watchfd == -1 {
			return "", err
		}

		isDir = fi.IsDir()
	}

	// 内核注册监听事件
	const registerAdd = unix.EV_ADD | unix.EV_CLEAR | unix.EV_ENABLE
	if err := register(w.kq, []int{watchfd}, registerAdd, flags); err != nil {
		unix.Close(watchfd)
		return "", err
	}

	// 构造文件描述符和路径的映射关系
	if !alreadyWatching {
		w.mu.Lock()
		w.watches[name] = watchfd
		w.paths[watchfd] = pathutils.PathInfo{Name: name, IsDir: isDir}
		w.mu.Unlock()
	}

	if isDir {
		// Watch the directory if it has not been watched before,
		// or if it was watched before, but perhaps only a NOTE_DELETE (watchDirectoryFiles)
		w.mu.Lock()

		watchDir := (flags&unix.NOTE_WRITE) == unix.NOTE_WRITE &&
			(!alreadyWatching || (w.dirFlags[name]&unix.NOTE_WRITE) != unix.NOTE_WRITE)
		// Store flags so this watch can be updated later
		// 构造监听的目录和监听事件的映射关系
		w.dirFlags[name] = flags
		w.mu.Unlock()

		if watchDir {
			// 监听目录下的所有文件
			if err := w.watchDirectoryFiles(name); err != nil {
				return "", err
			}
		}
	}
	return name, nil
}

// readEvents reads from kqueue and converts the received kevents into
// Event values that it sends down the Events channel.
func (w *Watcher) readEvents() {
	eventBuffer := make([]unix.Kevent_t, 10)

loop:
	for {
		// See if there is a message on the "done" channel
		select {
		case <-w.done:
			break loop
		default:
		}

		// Get new events
		kevents, err := read(w.kq, eventBuffer, &keventWaitTime)
		// EINTR is okay, the syscall was interrupted before timeout expired.
		if err != nil && err != unix.EINTR {
			select {
			case w.Errors <- err:
			case <-w.done:
				break loop
			}
			continue
		}

		// Flush the events we received to the Events channel
		for len(kevents) > 0 {
			kevent := &kevents[0]
			watchfd := int(kevent.Ident)
			mask := uint32(kevent.Fflags)
			w.mu.Lock()
			path := w.paths[watchfd]
			w.mu.Unlock()
			event := newEvent(path.Name, mask)

			if path.IsDir && !(event.Op&Remove == Remove) {
				// Double check to make sure the directory exists. This can happen when
				// we do a rm -fr on a recursively watched folders and we receive a
				// modification event first but the folder has been deleted and later
				// receive the delete event
				if _, err := os.Lstat(event.Name); os.IsNotExist(err) {
					// mark is as delete event
					event.Op |= Remove
				}
			}

			if event.Op&Rename == Rename || event.Op&Remove == Remove {
				w.Remove(event.Name)
				w.mu.Lock()
				delete(w.fileExists, event.Name)
				w.mu.Unlock()
			}

			if path.IsDir && event.Op&Write == Write && !(event.Op&Remove == Remove) {
				w.sendDirectoryChangeEvents(event.Name)
			} else {
				// Send the event on the Events channel.
				select {
				case w.Events <- event:
				case <-w.done:
					break loop
				}
			}

			if event.Op&Remove == Remove {
				// Look for a file that may have overwritten this.
				// For example, mv f1 f2 will delete f2, then create f2.
				if path.IsDir {
					fileDir := filepath.Clean(event.Name)
					w.mu.Lock()
					_, found := w.watches[fileDir]
					w.mu.Unlock()
					if found {
						// make sure the directory exists before we watch for changes. When we
						// do a recursive watch and perform rm -fr, the parent directory might
						// have gone missing, ignore the missing directory and let the
						// upcoming delete event remove the watch from the parent directory.
						if _, err := os.Lstat(fileDir); err == nil {
							w.sendDirectoryChangeEvents(fileDir)
						}
					}
				} else {
					filePath := filepath.Clean(event.Name)
					if fileInfo, err := os.Lstat(filePath); err == nil {
						w.sendFileCreatedEventIfNew(filePath, fileInfo)
					}
				}
			}

			// Move to next event
			kevents = kevents[1:]
		}
	}

	// cleanup
	err := unix.Close(w.kq)
	if err != nil {
		// only way the previous loop breaks is if w.done was closed so we need to async send to w.Errors.
		select {
		case w.Errors <- err:
		default:
		}
	}
	close(w.Events)
	close(w.Errors)
}

// newEvent returns a platform-independent Event based on kqueue Fflags.
func newEvent(name string, mask uint32) Event {
	e := Event{Name: name}
	if mask&unix.NOTE_DELETE == unix.NOTE_DELETE {
		e.Op |= Remove
	}
	if mask&unix.NOTE_WRITE == unix.NOTE_WRITE {
		e.Op |= Write
	}
	if mask&unix.NOTE_RENAME == unix.NOTE_RENAME {
		e.Op |= Rename
	}
	if mask&unix.NOTE_ATTRIB == unix.NOTE_ATTRIB {
		e.Op |= Chmod
	}
	return e
}

func newCreateEvent(name string) Event {
	return Event{Name: name, Op: Create}
}

// watchDirectoryFiles to mimic inotify when adding a watch on a directory.
func (w *Watcher) watchDirectoryFiles(dirPath string) error {
	// Get all files
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		filePath := filepath.Join(dirPath, fileInfo.Name())
		filePath, err = w.internalWatch(filePath, fileInfo)
		if err != nil {
			return err
		}

		w.mu.Lock()
		w.fileExists[filePath] = true
		w.mu.Unlock()
	}

	return nil
}

// sendDirectoryEvents searches the directory for newly created files
// and sends them over the event channel. This functionality is to have
// the BSD version of fsnotify match Linux inotify which provides a
// create event for files created in a watched directory.
func (w *Watcher) sendDirectoryChangeEvents(dirPath string) {
	// Get all files
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		select {
		case w.Errors <- err:
		case <-w.done:
			return
		}
	}

	// Search for new files
	for _, fileInfo := range files {
		filePath := filepath.Join(dirPath, fileInfo.Name())
		err := w.sendFileCreatedEventIfNew(filePath, fileInfo)

		if err != nil {
			return
		}
	}
}

// sendFileCreatedEvent sends a create event if the file isn't already being tracked.
func (w *Watcher) sendFileCreatedEventIfNew(filePath string, fileInfo os.FileInfo) (err error) {
	w.mu.Lock()
	_, doesExist := w.fileExists[filePath]
	w.mu.Unlock()
	if !doesExist {
		// Send create event
		select {
		case w.Events <- newCreateEvent(filePath):
		case <-w.done:
			return
		}
	}

	// like watchDirectoryFiles (but without doing another ReadDir)
	filePath, err = w.internalWatch(filePath, fileInfo)
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.fileExists[filePath] = true
	w.mu.Unlock()

	return nil
}

func (w *Watcher) internalWatch(name string, fileInfo os.FileInfo) (string, error) {
	if fileInfo.IsDir() {
		// mimic Linux providing delete events for subdirectories
		// but preserve the flags used if currently watching subdirectory
		w.mu.Lock()
		flags := w.dirFlags[name]
		w.mu.Unlock()

		flags |= unix.NOTE_DELETE | unix.NOTE_RENAME
		return w.addWatch(name, flags)
	}

	// watch file to mimic Linux inotify
	return w.addWatch(name, noteAllEvents)
}

// kqueue creates a new kernel event queue and returns a descriptor.
func kqueue() (kq int, err error) {
	kq, err = unix.Kqueue()
	if kq == -1 {
		return kq, err
	}
	return kq, nil
}

// register events with the queue.
func register(kq int, fds []int, flags int, fflags uint32) error {
	changes := make([]unix.Kevent_t, len(fds))

	for i, fd := range fds {
		// SetKevent converts int to the platform-specific types:
		unix.SetKevent(&changes[i], fd, unix.EVFILT_VNODE, flags)
		changes[i].Fflags = fflags
	}

	// register the events
	success, err := unix.Kevent(kq, changes, nil, nil)
	if success == -1 {
		return err
	}
	return nil
}

// read retrieves pending events, or waits until an event occurs.
// A timeout of nil blocks indefinitely, while 0 polls the queue.
func read(kq int, events []unix.Kevent_t, timeout *unix.Timespec) ([]unix.Kevent_t, error) {
	n, err := unix.Kevent(kq, nil, events, timeout)
	if err != nil {
		return nil, err
	}
	return events[0:n], nil
}
