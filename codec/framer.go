package codec

// IFramer defines how to read a data frame.
type IFramer interface {
	// ReadFrame 读取来自网络的的二进制数据
	ReadFrame() ([]byte, error)
}

// ISafeFramer is a special framer, provides an isSafe() method
// to describe if it is safe when concurrent read.
type ISafeFramer interface {
	IFramer
	// IsSafe returns if this framer is safe when concurrent read.
	IsSafe() bool
}

// IsSafeFramer returns if this framer is safe when concurrent read. The input
// parameter f should implement ISafeFramer interface. If not , this method will return false.
func IsSafeFramer(f interface{}) bool {
	framer, ok := f.(ISafeFramer)
	if ok && framer.IsSafe() {
		return true
	}
	return false
}
