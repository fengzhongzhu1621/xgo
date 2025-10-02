package output

// WriteMode is the log write mode, one of 1, 2, 3.
type WriteMode int

const (
	// WriteSync writes synchronously.
	WriteSync = 1
	// WriteAsync writes asynchronously.
	WriteAsync = 2
	// WriteFast writes fast(may drop logs asynchronously).
	WriteFast = 3
)
