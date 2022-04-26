// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

// Package fsnotify provides a platform-independent interface for file system notifications.
package fsnotify

// 定义事件类型

import (
	"bytes"
	"errors"
	"fmt"
)

// Event represents a single file system notification.
type Event struct {
	Name string // Relative path to the file or directory. 发送变化的文件或目录
	Op   Op     // File operation that triggered the event. 具体的变化
}

// Op describes a set of file operations.
type Op uint32

// 文件或目录变化的类型
// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod // 文件或目录的属性发生变化时触发，linux系统可以通过chmod命令改变文件或目录属性
)

// 按位判断触发的类型，如果类型为0，返回空字符串
func (op Op) String() string {
	// Use a buffer for efficient string concatenation
	var buffer bytes.Buffer

	if op&Create == Create {
		buffer.WriteString("|CREATE")
	}
	if op&Remove == Remove {
		buffer.WriteString("|REMOVE")
	}
	if op&Write == Write {
		buffer.WriteString("|WRITE")
	}
	if op&Rename == Rename {
		buffer.WriteString("|RENAME")
	}
	if op&Chmod == Chmod {
		buffer.WriteString("|CHMOD")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:] // Strip leading pipe
}

// String returns a string representation of the event in the form
// "file: REMOVE|WRITE|..."
func (e Event) String() string {
	return fmt.Sprintf("%q: %s", e.Name, e.Op.String())
}

// Common errors that can be reported by a watcher
var (
	ErrEventOverflow = errors.New("fsnotify queue overflow")
)
