package stringutils

import (
	"bytes"
	"fmt"
	"strings"
)

// -----------------------------------------------------------------
// 字符串拼接
// -----------------------------------------------------------------

// Deprecated: 低效的字符串拼接.
// 每次操作都会创建新字符串适合少量、不频繁的拼接操作
func StringPlus(p []string) string {
	var s string
	l := len(p)
	for i := 0; i < l; i++ {
		s += p[i]
	}
	return s
}

// StringFmt Deprecated: 低效的字符串拼接.
// 内部使用反射，性能较差
func StringFmt(p []interface{}) string {
	return fmt.Sprint(p...)
}

// StringBuffer Deprecated: 低效的字符串拼接.
// 转换为字符串时有额外内存分配
// 可读写缓冲区
// 线程安全
func StringBuffer(p []string) string {
	var b bytes.Buffer
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

// StringJoin 专为连接字符串切片设计（推荐）
// 需要预先构建切片
// 适合已存在字符串集合的情况
func StringJoin(p []string) string {
	return strings.Join(p, "")
}

func StringBuilder(p []string) string {
	var b strings.Builder
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

// StringBuilderEx Go 1.10+引入，专为字符串构建优化（推荐）
// 最小化内存分配和复制
// 可预分配缓冲区大小
// 线程不安全
func StringBuilderEx(p []string, n int) string {
	var b strings.Builder
	l := len(p)
	// 预分配空间：实现分配足够的内容，减少运行时的内存分配
	b.Grow(n)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}
