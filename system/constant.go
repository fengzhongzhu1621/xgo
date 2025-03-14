package shell

var NonLoginShellMap = map[string]bool{
	"/bin/false":     true,
	"/sbin/nologin":  true,
	"/sbin/shutdown": true,
}

// LoginShellMap shell 程序是否可以登录配置
var LoginShellMap = map[string]bool{
	"/bin/sh":           true,
	"/bin/bash":         true,
	"/sbin/nologin":     false,
	"/usr/bin/sh":       true,
	"/usr/bin/bash":     true,
	"/usr/sbin/nologin": false,
	"/bin/tcsh":         true,
	"/bin/csh":          true,
	"/bin/ksh":          true,
	"/bin/rksh":         true,
	"/bin/zsh":          true,
}
