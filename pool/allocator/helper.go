package allocator

// Malloc 从池中获取一个[]byte，第二个返回值用于Free操作
//
// size: 需要分配的字节大小
// returns:
//   - 第一个返回值: 实际可用的字节切片(长度=size)
//   - 第二个返回值: 内部使用的原始分配对象(用于释放)
func Malloc(size int) ([]byte, interface{}) {
	return defaultAllocator.Malloc(size)
}

// Free 将字节释放回池中
//
// bts: 由Malloc返回的第二个参数，包含原始分配对象
func Free(bts interface{}) {
	defaultAllocator.Free(bts)
}
