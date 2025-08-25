// package linkbuffer 链表缓冲区
package linkbuffer

import (
	"io"

	"github.com/fengzhongzhu1621/xgo/pool/allocator"
)

// 确保Buf类型实现了IBuffer接口
var _ IBuffer = (*Buf)(nil)

// Buf is a rich buffer to reuse underlying bytes.
// 是一个功能丰富的缓冲区，用于重用底层字节数组；使用链表结构管理数据块，支持高效的内存重用
// 支持多种读写操作，包括追加、预置、分配和读取
// 实现了标准的io.Writer和io.Reader接口
// 通过内存分配器管理内存分配和释放，减少GC压力
// 使用dirty链表管理待释放的内存块，延迟实际的内存释放操作
//
// 这种设计适合需要高效内存管理和大量I/O操作的场景，如网络编程、文件处理等
type Buf struct {
	a             allocator.IAllocator // 字节切片分配器，用于分配和释放内存
	minMallocSize int                  // 每次分配的最小字节大小

	head, tail *bytes // 双向链表的头尾指针，用于管理缓冲区数据块
	dirty      *bytes // 指向已释放但待清理的数据块链表
}

// NewLinkBuf creates a new buf. 创建一个新的缓冲区链表实例
func NewLinkBuf(a allocator.IAllocator, minMallocSize int) *Buf {
	bytes := NewBytes(nil, nil) // 创建一个空的bytes节点作为初始节点
	return &Buf{
		a:             a,             // 设置内存分配器
		minMallocSize: minMallocSize, // 设置链表中每个节点的最小分配大小
		head:          bytes,         // 初始化头节点
		tail:          bytes,         // 初始化尾节点
	}
}

// Write copies p to Buf and implements io.Writer. 将字节切片p写入缓冲区，实现io.Writer接口
func (b *Buf) Write(p []byte) (int, error) {
	// 如果尾节点没有释放函数（表示未分配），则需要分配新空间；节点未分配内存
	if b.tail.release == nil {
		bts, release := b.a.Malloc(b.minMallocSize) // 从分配器分配内存
		// bts[:0] 这行代码的作用是创建一个长度为 0，但容量保持不变的切片，指向 bts 的底层数组
		// 新切片虽然长度为 0，但其容量（capacity）与 bts 相同。它仍然指向 b.a.Malloc 分配出来的那块原始内存空间。
		// 这样做的目的是为了后续通过 append 操作向这个新切片添加数据时，可以直接复用这块已分配的内存，避免了立即再次分配内存的开销
		b.tail.next = NewBytes(bts[:0], release) // 创建新的bytes节点，追加到链表的最后
		b.tail = b.tail.next                     // 移动尾指针
	}

	available := cap(b.tail.bts) - len(b.tail.bts) // 计算当前尾节点的剩余空间
	if available >= len(p) {                       // 如果空间足够容纳所有数据
		b.tail.bts = append(b.tail.bts, p...) // 直接追加数据
		return len(p), nil
	}

	// 空间不足时，先填充剩余空间
	b.tail.bts = append(b.tail.bts, p[:available]...)
	// 分配新的内存块
	bts, release := b.a.Malloc(b.minMallocSize)
	b.tail.next = NewBytes(bts[:0], release)
	b.tail = b.tail.next

	// 递归写入剩余数据，返回写入的数据的数量
	n, err := b.Write(p[available:])

	return available + n, err
}

// Append appends a slice of bytes to Buf.
// Buf owns these bs, but won't release them to underlying allocator.
// 追加一个或多个字节切片到缓冲区，但不获取这些切片的所有权
// 不从sync.Pool中获取内存块
func (b *Buf) Append(bs ...[]byte) {
	for _, bts := range bs {
		b.append(bts)
	}
}

func (b *Buf) append(bts []byte) {
	// 如果尾节点已满或未分配，直接创建新节点添加到链表尾部
	if b.tail.release == nil || cap(b.tail.bts) == len(b.tail.bts) {
		b.tail.next = NewBytes(bts, nil) // 创建新节点，不设置释放函数
		b.tail = b.tail.next             // 移动尾指针
	} else {
		// 处理部分填充的情况
		remains := b.tail.bts[len(b.tail.bts):] // 获取当前尾节点的剩余空间
		release := b.tail.release               // 保存当前尾节点的释放函数
		b.tail.release = nil                    // 清空释放函数，表示不释放

		// 创建新节点存放追加的数据
		b.tail.next = NewBytes(bts, nil)
		b.tail = b.tail.next

		// 创建另一个新节点存放原来的剩余空间
		b.tail.next = NewBytes(remains, release)
		b.tail = b.tail.next
	}
}

// Prepend prepends a slice to bytes to Buf. Next Read starts with the first bytes of slice.
// Buf owns these bs, but won't release them to underlying allocator.
// 将一个或多个字节切片预置到缓冲区头部
func (b *Buf) Prepend(bs ...[]byte) {
	for i := len(bs) - 1; i >= 0; i-- {
		bytes := NewBytes(bs[i], nil)
		bytes.next = b.head
		b.head = bytes
	}
}

// Alloc allocates a []byte with size n.
// 在缓冲区尾部分配指定大小的字节切片
func (b *Buf) Alloc(n int) []byte {
	// 检查尾节点是否有足够空间
	if b.tail.release != nil && cap(b.tail.bts)-len(b.tail.bts) >= n {
		l := len(b.tail.bts)          // 记录当前长度
		b.tail.bts = b.tail.bts[:l+n] // 扩展切片大小
		return b.tail.bts[l : l+n]    // 返回新分配的区域
	}

	// 尾节点空间不足，分配新节点
	bts, release := b.a.Malloc(n)            // 直接从分配器分配所需大小的内存
	b.tail.next = NewBytes(bts[:n], release) // 创建新节点
	b.tail = b.tail.next                     // 移动尾指针
	return bts[:n]                           // 返回分配的切片
}

// Prelloc allocates a []byte with size n at the beginning of Buf.
// 在缓冲区头部分配指定大小的字节切片
func (b *Buf) Prelloc(n int) []byte {
	bts, release := b.a.Malloc(n)       // 从分配器分配内存
	bytes := NewBytes(bts[:n], release) // 创建新节点
	bytes.next = b.head                 // 新节点指向当前头节点
	b.head = bytes                      // 更新头指针
	return bts[:n]                      // 返回分配的切片
}

// Merge merges another Reader.
// If r is not *Buf, b does not own the bytes of r.
// If r is a *Buf, the ownership of r's bytes is changed to b, and the caller should not Release r.
// 合并另一个读取器到当前缓冲区
func (b *Buf) Merge(r IReader) {
	bb, ok := r.(*Buf) // 尝试转换为*Buf类型
	if !ok {
		// 如果不是*Buf类型，读取所有数据并追加
		for _, bts := range r.ReadAll() {
			b.Append(bts)
		}
		return
	}
	// 如果是*Buf类型，直接链接链表
	b.tail.next = bb.head // 将当前尾节点指向另一个缓冲区的头节点
	b.tail = bb.tail      // 更新尾指针为另一个缓冲区的尾节点
}

// Read copies data to p, and returns the number of byte copied and an error.
// The io.EOF is returned if Buf has no unread bytes and len(p) is not zero.
// 从缓冲区读取数据到p，实现io.Reader接口
func (b *Buf) Read(p []byte) (int, error) {
	if len(p) == 0 { // 如果目标切片为空，直接返回
		return 0, nil
	}

	defer b.ensureNotEmpty() // 确保缓冲区不为空
	var copied int           // 记录已拷贝字节数

	// 遍历所有数据块
	for b.head != nil {
		curCopied := copy(p[copied:], b.head.bts) // 从当前数据块拷贝数据
		copied += curCopied                       // 更新已拷贝字节数
		b.head.bts = b.head.bts[curCopied:]       // 更新当前数据块的剩余数据
		// 清理已空的数据块
		b.dirtyEmptyHeads()

		if copied == len(p) { // 如果已填满目标切片
			return copied, nil
		}
	}

	if copied > 0 { // 如果读取了部分数据
		return copied, nil
	}
	return copied, io.EOF // 没有数据可读，返回EOF
}

// ReadN tries best to read all size into one []byte.
// The second return value may be smaller than size if underlying bytes is not continuous.
// 尝试读取指定大小的数据到一个连续的字节切片中
func (b *Buf) ReadN(size int) ([]byte, int) {
	defer b.ensureNotEmpty() // 确保缓冲区不为空
	b.dirtyEmptyHeads()      // 清理已空的数据块

	for b.head != nil {
		if size >= len(b.head.bts) { // 如果需要的大小大于等于当前块大小
			bts := b.dirtyHead() // 获取整个数据块
			b.dirtyEmptyHeads()  // 清理已空的数据块
			return bts, len(bts) // 返回整个数据块
		}

		// 只需要当前数据块的一部分
		bts := b.head.bts[:size]       // 获取前size字节
		b.head.bts = b.head.bts[size:] // 更新当前数据块的剩余数据
		return bts, size               // 返回部分数据
	}
	return nil, 0 // 没有数据可读
}

// ReadAll returns all underlying []byte in [][]byte.
// 读取所有底层数据块，返回二维字节切片
func (b *Buf) ReadAll() [][]byte {
	defer b.ensureNotEmpty() // 确保缓冲区不为空
	var all [][]byte         // 存储所有数据块

	for b.head != nil {
		if bts := b.dirtyHead(); len(bts) != 0 { // 获取并清理数据块
			all = append(all, bts) // 添加到结果切片
		}
	}
	return all
}

// ReadNext returns the next continuous []byte.
// 读取下一个连续的数据块
func (b *Buf) ReadNext() []byte {
	defer b.ensureNotEmpty() // 确保缓冲区不为空

	for b.head != nil {
		if bts := b.dirtyHead(); len(bts) != 0 { // 获取并清理数据块
			return bts
		}
	}
	return nil // 没有数据可读
}

// Release releases the read bytes to allocator.
// 释放已读取的字节到底层分配器
func (b *Buf) Release() {
	// 遍历dirty链表，释放所有待释放的节点
	for b.dirty != nil {
		b.a.Free(b.dirty.release) // 调用分配器的Free方法释放内存

		dirty := b.dirty       // 保存当前节点
		b.dirty = b.dirty.next // 移动dirty指针

		// 清空当前节点属性
		dirty.release = nil
		dirty.bts = nil
		dirty.next = nil

		BytesPool.Put(dirty) // 将节点放回对象池
	}
}

// Len returns the total len of underlying bytes.
// 返回缓冲区中未读取字节的总长度
func (b *Buf) Len() int {
	var l int // 记录总长度

	// 遍历所有数据块
	for bytes := b.head; bytes != nil; bytes = bytes.next {
		l += len(bytes.bts) // 累加每个数据块的长度
	}
	return l
}

// dirtyEmptyHeads 清理所有已空的数据块
func (b *Buf) dirtyEmptyHeads() {
	// 循环处理头节点为空的情况
	for b.head != nil && len(b.head.bts) == 0 {
		// 头指针移动到下一个节点
		b.dirtyHead() // 将空节点移到dirty链表
	}
}

// dirtyHead 将当前头节点移动到dirty链表以待释放
func (b *Buf) dirtyHead() []byte {
	bts := b.head.bts // 保存当前头节点的数据

	head := b.head     // 保存当前头节点
	b.head = head.next // 移动头指针到下一个节点

	if head.release == nil { // 如果节点没有分配器释放函数
		// 清空节点属性并放回对象池
		head.bts = nil
		head.next = nil
		BytesPool.Put(head)
	} else {
		// 将节点添加到dirty链表头部
		head.next = b.dirty
		b.dirty = head
	}

	return bts // 返回获取的数据
}

// ensureNotEmpty 确保缓冲区始终有一个空节点
func (b *Buf) ensureNotEmpty() {
	if b.head == nil { // 如果头节点为空
		bytes := NewBytes(nil, nil) // 创建新的空节点
		b.head = bytes              // 设置头节点
		b.tail = bytes              // 设置尾节点
	}
}
