package buffer

import (
	"bufio"
	"io"
)

// DefaultReaderSize is the default size of reader in bit.
const DefaultReaderSize = 4 * 1024

// readerSizeConfig is the default size of buffer when framer read package.
var readerSizeConfig = DefaultReaderSize

// NewReaderSize returns a reader with read buffer. Size <= 0 means no buffer.
func NewReaderSize(r io.Reader, size int) io.Reader {
	if size <= 0 {
		return r
	}
	return bufio.NewReaderSize(r, size)
}

// NewReader returns reader with the default buffer size.
func NewReader(r io.Reader) io.Reader {
	return bufio.NewReaderSize(r, readerSizeConfig)
}

// GetReaderSize returns size of read buffer in bit.
func GetReaderSize() int {
	return readerSizeConfig
}

// SetReaderSize sets the size of read buffer in bit.
func SetReaderSize(size int) {
	readerSizeConfig = size
}
