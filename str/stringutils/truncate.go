package stringutils

// Truncate 从字符串获取指定长度的子串
func Truncate(s string, n int) string {
	if n > len(s) {
		return s
	}
	return s[:n]
}
