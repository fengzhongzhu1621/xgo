package codec

// Framer defines how to read a data frame.
type Framer interface {
	ReadFrame() ([]byte, error)
}

// SafeFramer is a special framer, provides an isSafe() method
// to describe if it is safe when concurrent read.
type SafeFramer interface {
	Framer
	// IsSafe returns if this framer is safe when concurrent read.
	IsSafe() bool
}

// IsSafeFramer returns if this framer is safe when concurrent read. The input
// parameter f should implement SafeFramer interface. If not , this method will return false.
func IsSafeFramer(f interface{}) bool {
	framer, ok := f.(SafeFramer)
	if ok && framer.IsSafe() {
		return true
	}
	return false
}
