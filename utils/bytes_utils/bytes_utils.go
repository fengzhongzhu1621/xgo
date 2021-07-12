package bytes_utils

import "bytes"

/**
 * 判断前缀和后缀是否全部匹配
 */
func HasPrefixAndSuffix(s, prefix []byte, suffix []byte) bool {
	return bytes.HasPrefix(s, prefix) && bytes.HasSuffix(s, suffix)
}


/**
 * 去掉结尾换行符
 * trimEOL cuts unixy style \n and windowsy style \r\n suffix from the string
 */
func TrimEOL(b []byte) []byte {
	lns := len(b)
	if lns > 0 && b[lns-1] == '\n' {
		lns--
		if lns > 0 && b[lns-1] == '\r' {
			lns--
		}
	}
	return b[:lns]
}
