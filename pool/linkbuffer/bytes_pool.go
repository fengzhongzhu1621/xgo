package linkbuffer

import "sync"

// bytes 实现了一个基于 sync.Pool 的字节缓冲区池，
// 表示一个可重用的字节缓冲区
type bytes struct {
	bts     []byte
	release interface{} // 已分配的字节切片地址（后续需要主动释放）
	next    *bytes      // 指向下一个分配的字节切片，实现链表结构
}

var BytesPool = sync.Pool{New: func() interface{} { return &bytes{} }}

// NewBytes 创建或从池中获取一个Bytes实例
// bts 分配的内存字节切片
// release 分配的内存字节切片的地址，根据此指针进行释放
func NewBytes(bts []byte, release interface{}) *bytes {
	bytes := BytesPool.Get().(*bytes)
	// 初始化
	bytes.bts = bts
	bytes.release = release
	bytes.next = nil

	return bytes
}

// Bytes 返回底层字节切片
func (b *bytes) Bytes() []byte {
	return b.bts
}

// SetBytes 设置字节切片
func (b *bytes) SetBytes(bts []byte) {
	b.bts = bts
}
