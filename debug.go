package xgo

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// 支持的最小版本
const xGoSupportMinGoVer = 10

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
		fmt.Fprintf(DefaultWriter, "[XGO-debug] "+format, values...)
	}
}

/**
 * 获得go的最小版本
 */
func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

/**
 * 打印版本匹配信息
 */
func debugPrintWARNINGDefault() {
	// go1.14.4
	if v, e := getMinVer(runtime.Version()); e == nil && v <= xGoSupportMinGoVer {
		debugPrint(`[WARNING] Now xGo requires Go 1.11 or later and Go 1.12 will be required soon.

`)
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

func debugPrintError(err error) {
	if err != nil {
		if IsDebugging() {
			fmt.Fprintf(DefaultErrorWriter, "[XGO-debug] [ERROR] %v\n", err)
		}
	}
}
