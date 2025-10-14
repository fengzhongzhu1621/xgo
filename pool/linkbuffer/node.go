package linkbuffer

import "sync"

// bytesNode 实现了一个基于 sync.Pool 的字节缓冲区池（链表结构中的节点定义）
// 表示一个可重用的字节缓冲区
type bytesNode struct {
	bts     []byte      // 节点存放的内容
	release interface{} // 已分配的字节切片地址（后续需要主动释放）
	next    *bytesNode  // 指向下一个分配的字节切片，实现链表结构
}

// BytesPool 全局字节缓冲区池实例
// 使用 sync.Pool 实现，每次从池中获取一个 bytes 实例，并初始化创建一个空的节点
var BytesPool = sync.Pool{New: func() interface{} { return &bytesNode{} }}

// NewBytesNode 创建或从池中获取一个Bytes实例
// bts 分配的内存字节切片
// release 分配的内存字节切片的地址，根据此指针进行释放
func NewBytesNode(bts []byte, release interface{}) *bytesNode {
	bytes := BytesPool.Get().(*bytesNode)
	// 初始化
	bytes.bts = bts
	bytes.release = release
	bytes.next = nil

	return bytes
}

// Bytes 返回底层字节切片
func (b *bytesNode) Bytes() []byte {
	return b.bts
}

// SetBytes 设置字节切片
func (b *bytesNode) SetBytes(bts []byte) {
	b.bts = bts
}
