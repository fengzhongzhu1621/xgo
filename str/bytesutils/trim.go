package bytesutils

// trimEOL 去掉结尾换行符 cuts unixy style \n and windowsy style \r\n suffix from the string.
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
