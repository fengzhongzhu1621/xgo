package bytesutils

// TruncateBytes 截断字节切片
func TruncateBytes(content []byte, length int) []byte {
	if len(content) > length {
		return content[:length]
	}
	return content
}
