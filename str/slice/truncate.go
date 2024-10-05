package slice

// TruncateBytes 截断字节切片
func TruncateBytes(content []byte, length int) []byte {
	if len(content) > length {
		return content[:length]
	}
	return content
}

// TruncateBytesToString 截断字节切片，并将字节切片转换为字符串
func TruncateBytesToString(content []byte, length int) string {
	s := TruncateBytes(content, length)
	return string(s)
}
