package xgo

import (
	"fmt"
	"strings"
)

/**
 * 判断是否是调试模式
 */
func IsDebugging() bool {
	return xGoMode == debugCode
}


/**
 * 调试模式打印函数
 */
func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		// 添加换行符
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		// 默认打印调试信息到标准输出
		fmt.Fprintf(DefaultWriter, "[GIN-debug] "+format, values...)
	}
}

/**
 * 如果当前是调试模式，则打印警告信息
 */
func debugPrintWARNINGNew() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export XGO_MODE=release
 - using code:	xgo.SetMode(xgo.ReleaseMode)

`)
}
