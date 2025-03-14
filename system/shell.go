package shell

import "strings"

// IsLoginShell 判断 shell 程序是否可以登录
func IsLoginShell(shell string) bool {
	if strings.HasSuffix(shell, "/nologin") {
		return false
	}

	// 匹配黑名单，从配置中获取是否可以登录标识
	value, exists := NonLoginShellMap[shell]
	if exists && !value {
		return false
	}

	return true
}

// IsPasswordDisabled 判断密码是否无效
func IsPasswordDisabled(password string) int {
	if password == "!!" || password == "!" || password == "*" {
		return 1
	}

	return 0
}
