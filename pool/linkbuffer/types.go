package linkbuffer

import "io"

// Buffer is the interface of link buffer.
type IBuffer interface {
	IReader
	IWriter

	Release()
}

// Reader is the interface to read from link buffer.
type IReader interface {
	io.Reader
	ReadN(size int) ([]byte, int)
	ReadAll() [][]byte
	ReadNext() []byte
}

// Writer is the interface to write to link buffer.
type IWriter interface {
	io.Writer
	Append(...[]byte)
	Prepend(...[]byte)
	Alloc(size int) []byte
	Prelloc(size int) []byte
	Len() int
	Merge(IReader)
}
