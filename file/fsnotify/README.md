# 简介
用于监控文件系统中的变化，包括文件和目录的修改、删除、创建等。

* 文件和目录监控：监控单个文件或目录的变化。
* 多平台支持：支持 Linux、macOS 和 Windows 操作系统。对于 Linux，使用inotify，对于 macOS，使用kqueue，对于 Windows，使用ReadDirectoryChangesW。
* 高效的事件通知：基于操作系统提供的文件系统通知机制。
* 文件监控粒度：支持多种文件事件，支持文件和目录级别的事件监听。

```
go get github.com/fsnotify/fsnotify
```


# 事件类型
* fsnotify.Create：文件或目录被创建。
* fsnotify.Remove：文件或目录被删除。
* fsnotify.Rename：文件或目录被重命名。
* fsnotify.Write：文件内容被写入。
* fsnotify.Chmod：文件权限被修改。


# 注意
* 事件丢失：操作系统的文件系统通知机制并不是完全可靠，在高频繁的文件操作时，可能会错过一些事件。为了避免丢失事件，可以结合定时器定期扫描文件状态。
* 文件句柄限制：操作系统对每个进程可打开的文件句柄数量是有限制的。监控大量文件，可能会遇到这个限制。Linux 上，你可以通过 `ulimit -n` 查看并设置文件句柄数量。
