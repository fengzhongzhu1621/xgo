package bytes

import "bytes"

func HasPrefixAndSuffix(s, prefix []byte, suffix []byte) bool {
	return bytes.HasPrefix(s, prefix) && bytes.HasSuffix(s, suffix)
}
