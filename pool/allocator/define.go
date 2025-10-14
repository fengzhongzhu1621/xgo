package allocator

// Allocator is the interface to Malloc or Free bytes.
type IAllocator interface {
	// Malloc mallocs a []byte with specific size.
	// The second return value is the consequence for go's escape analysis.
	// See ClassAllocator and https://github.com/golang/go/issues/8618 for details.
	Malloc(int) ([]byte, interface{})
	// Free frees the allocated bytes. It accepts the second return value of Malloc.
	Free(interface{})
}
