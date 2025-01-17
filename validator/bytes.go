package validator

import "bytes"

// HasPrefixAndSuffix 判断前缀和后缀是否全部匹配.
func HasPrefixAndSuffix(s, prefix []byte, suffix []byte) bool {
	return bytes.HasPrefix(s, prefix) && bytes.HasSuffix(s, suffix)
}
