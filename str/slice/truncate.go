package slice

func TruncateBytes(content []byte, length int) []byte {
	if len(content) > length {
		return content[:length]
	}
	return content
}

func TruncateBytesToString(content []byte, length int) string {
	s := TruncateBytes(content, length)
	return string(s)
}
