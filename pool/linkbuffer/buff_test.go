package linkbuffer

import (
	stdbytes "bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/fengzhongzhu1621/xgo/pool/allocator"
	"github.com/stretchr/testify/require"
)

var _ allocator.IAllocator = (*wrappedAllocator)(nil)

// 定义wrappedAllocator结构体，用于包装一个 allocator.IAllocator 并进行测试监控
type wrappedAllocator struct {
	t           *testing.T           // 测试对象，用于报告测试失败和错误
	a           allocator.IAllocator // 被包装的实际内存分配器
	malloced    map[*byte]struct{}   // 记录当前已分配（尚未释放）的内存块起始地址的映射。使用map来快速查找，值用空结构体节省空间。
	mallocTimes int                  // 记录总的内存分配次数
}

func newWrappedAllocator(t *testing.T, a allocator.IAllocator) *wrappedAllocator {
	return &wrappedAllocator{t: t, a: a, malloced: make(map[*byte]struct{})}
}

func (a *wrappedAllocator) Malloc(size int) ([]byte, interface{}) {
	// 调用被包装分配器的Malloc方法实际分配内存
	// bts: 分配得到的字节切片，free: 与该内存关联的释放句柄（通常是一个释放函数或标识）
	bts, free := a.a.Malloc(size)
	// 记录本次分配的内存的起始地址（通过取切片第一个元素的地址）到malloced映射中，表示该内存已分配
	a.malloced[&bts[0]] = struct{}{}
	// 分配次数加1
	a.mallocTimes++
	// 返回分配的字节切片和释放句柄
	return bts, free
}

func (a *wrappedAllocator) Free(free interface{}) {
	// 先将释放请求委托给被包装的分配器
	a.a.Free(free)
	if _, ok := a.malloced[&free.([]byte)[0]]; !ok {
		require.FailNow(a.t, "free unknown bytes")
	}
	// 从已分配记录map中删除该内存块的记录
	delete(a.malloced, &free.([]byte)[0])
}

// MustMallocTimes 断言到当前为止的内存分配次数等于预期值n
func (a *wrappedAllocator) MustMallocTimes(n int) {
	require.Equal(a.t, n, a.mallocTimes)
}

// MustAllFreed 断言所有分配的内存都已被正确释放（即malloced映射为空）
func (a *wrappedAllocator) MustAllFreed() {
	require.Empty(a.t, a.malloced)
}

func newBytesAllocator() *bytesAllocator {
	return &bytesAllocator{pool: make(map[int][]interface{})}
}

// /////////////////////////////////////////////////////////////////////////////////
var _ allocator.IAllocator = (*bytesAllocator)(nil)

type bytesAllocator struct {
	// pool是一个映射，用于按大小（int）缓存可重用的内存块（[]interface{} 实际存储的是 []byte）
	pool map[int][]interface{}
}

// Malloc 分配指定大小的字节切片，优先从池中获取，池中没有则新建
func (a *bytesAllocator) Malloc(n int) ([]byte, interface{}) {
	// 根据索引获得内存块池（索引值就是内存块容量）
	bsPool, ok := a.pool[n]
	if ok && len(bsPool) != 0 {
		// 从该大小对应的切片末尾取出一个内存块（后进先出，LIFO）
		bts := bsPool[len(bsPool)-1]
		// 从池中移除该内存块
		a.pool[n] = bsPool[:len(bsPool)-1]
		// 返回该内存块（转换为[]byte）和其本身作为释放句柄（因为这里释放就是放回池子）
		return bts.([]byte), bts
	}

	// 如果池中没有该大小的可用内存块，则使用make新建一个指定大小n的字节切片
	bts := make([]byte, n)
	return bts, bts
}

// Free 将内存块放回池中以便重用
func (a *bytesAllocator) Free(v interface{}) {
	bts := v.([]byte)
	// 将切片重置为其原始容量（可能大于当前长度），以便后续重用时可使用全部容量
	bts = bts[:cap(bts)]
	a.pool[len(bts)] = append(a.pool[len(bts)], v)
}

// /////////////////////////////////////////////////////////////////////////////////
func BenchmarkBuf(b *testing.B) {
	bigBts := make([]byte, 1<<10)
	b.Run("link_buffer_bigBytes", func(b *testing.B) {
		b.ReportAllocs()
		bb := NewLinkBuf(allocator.NewClassAllocator(), 1<<9)
		for i := 0; i < b.N; i++ {
			bb.Append(bigBts)
			bb.ReadNext()
			bb.Release()
		}
	})
	b.Run("copy_each_bigBytes", func(b *testing.B) {
		b.ReportAllocs()
		var bb []byte
		for i := 0; i < b.N; i++ {
			bb = append(bb, bigBts...)
		}
	})
	b.Run("link_buffer_reuse", func(b *testing.B) {
		b.ReportAllocs()
		r := rand.New(rand.NewSource(1))
		bb := NewLinkBuf(allocator.NewClassAllocator(), 1<<10)
		for i := 0; i < b.N; i++ {
			copy(bb.Alloc(16), bigBts)
			copy(bb.Alloc(int(r.Int31()%1<<20+1)), bigBts)
			bb.ReadNext()
			bb.ReadNext()
			bb.Release()
		}
	})
	b.Run("std_buffer", func(b *testing.B) {
		b.ReportAllocs()
		r := rand.New(rand.NewSource(1))
		for i := 0; i < b.N; i++ {
			bb := stdbytes.Buffer{}
			bb.Write(bigBts[:16])
			bb.Write(bigBts[:r.Int31()%1<<20+1])
		}
	})
	b.Run("bytes_cannot_reuse", func(b *testing.B) {
		b.ReportAllocs()
		r := rand.New(rand.NewSource(1))
		for i := 0; i < b.N; i++ {
			bts := make([]byte, r.Int31()%1<<20+16)
			copy(bts[:16], bigBts)
			copy(bts[16:], bigBts)
		}
	})
}

func TestBuf_Write(t *testing.T) {
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	// 创建一个链表
	b := NewLinkBuf(wa, 4)

	// 因为链表为空，所以会创建一个节点（容量为4），写入123后，容量变为1
	n, err := b.Write([]byte("123"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	// 链表最后一个节点剩余1个字节的空间，不够写入，所以在链表结尾添加了一个新的节点
	// 4 写入到第一个节点
	// 5 写入到第二个节点
	n, err = b.Write([]byte("45"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	// 内存分配了两次
	wa.MustMallocTimes(2)
}

func TestBuf_Read(t *testing.T) {
	// 创建空链表
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 分配了3个节点（内存池分配）
	n, err := b.Write([]byte("1234567890"))
	require.Nil(t, err)
	require.Equal(t, 10, n)

	// 读取了3个字节，第一个节点剩余1个字节 123
	bts := make([]byte, 3)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 3, n)
	require.Equal(t, "123", string(bts))

	// 在第一个节点读取剩余的1个字节 4
	// 在第二个字节读取2个字节 56
	_, err = b.Read(bts)
	require.Nil(t, err)

	// 在第二个字节读取剩余的2个字节 78
	// 在第三个字节读取1个字节 9
	_, err = b.Read(bts)
	require.Nil(t, err)

	// 在第四个字节读取剩余的3个字节 0
	n, err = b.Read(bts)
	require.Nil(t, err)

	require.Equal(t, 1, n)
	require.Equal(t, "089", string(bts))

	n, err = b.Read(bts)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 0, n)

	// 释放节点中的缓存
	// 释放节点
	b.Release()

	wa.MustMallocTimes(3) // 分配了3个节点
	wa.MustAllFreed()
}

func TestBuf_Append(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 创建第一个节点并写入（内存池分配）
	n, err := b.Write([]byte("12"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	// 第一个节点剩余2个字节未填充，将此节点的release 设置为 nil cap=4
	// 创建第二个节点，内容是 []byte("345")， realease = nil len=3 cap=3（非内存池分配）
	// 创建第三个节点，内容是第一个节点剩余的两个未填充节点，release是第一个节点的缓冲区内存地址，共享第一个节点分配的内存块 len=0 cap=2（内存池不分配）
	b.Append([]byte("345"))

	// the remaining two bytes is available, no need to malloc new bytes.
	// 将 67 写入到第三个节点
	n, err = b.Write([]byte("67"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	bts := make([]byte, 7)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 7, n)
	require.Equal(t, "1234567", string(bts))

	b.Release()
	wa.MustMallocTimes(1) // 只从sync.Pool分配了一次内存（第一个节点）
	wa.MustAllFreed()
}

func TestBuf_Prepend(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 创建第一个节点并写入（内存池分配）
	n, err := b.Write([]byte("12"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	// 创建一个新的节点，追加到header的前面（非内存池分配）
	// PrependNode(345,len=3,cap=3,release=nil,next=head) -> head(len=0,cap=0,release=nil,next=第一个节点) -> 第一个节点
	b.Prepend([]byte("345"))

	bts := make([]byte, 5)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, "34512", string(bts))

	b.Release()
	wa.MustMallocTimes(1) // 只从sync.Pool分配了一次内存（第一个节点）
	wa.MustAllFreed()
}

func TestBuf_Alloc(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 创建第一个节点并写入，剩余2个空闲字节（内存池分配）
	n, err := b.Write([]byte("12"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	// 在第一个节点的剩余空间预留一个字节，剩余1个空闲字节（非内存池分配）
	bts := b.Alloc(1)
	require.Len(t, bts, 1)
	// 在第一个节点写入一个字节 4
	// 创建第二个节点，写入字节 56，剩余2个空闲字节（内存池分配）
	n, err = b.Write([]byte("456"))
	require.Nil(t, err)
	require.Equal(t, 3, n)
	// 预留的字节写入值3
	bts[0] = '3'

	// 因为第二个节点只剩余2个空闲字节，不够分配
	// 创建第三个节点（内存池分配）
	bts = b.Alloc(3)
	require.Len(t, bts, 3)
	// 填充预留节点
	copy(bts, "789")

	bts = make([]byte, 9)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 9, n)
	require.Equal(t, "123456789", string(bts))

	b.Release()
	wa.MustMallocTimes(3)
	wa.MustAllFreed()
}

func TestBuf_Prelloc(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 创建第一个节点（内存池分配）
	n, err := b.Write([]byte("12"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	// 创建第-1个节点，预留1个字节 len=1, cap=1 （内存池分配）
	bts := b.Prelloc(1)
	require.Len(t, bts, 1)
	// 创建第-2个字节，追加 456 （非内存池分配）
	b.Prepend([]byte("456"))
	bts[0] = '3' // 填充预留节点

	// 创建-3个节点，预留3个字节 len=3,cap=3（内存池分配）
	bts = b.Prelloc(3)
	require.Len(t, bts, 3)
	copy(bts, "789") // 填充预留节点

	bts = make([]byte, 9)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 9, n)
	require.Equal(t, "789456312", string(bts))

	b.Release()
	wa.MustMallocTimes(3)
	wa.MustAllFreed()
}

func TestBuf_Merge(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b1 := NewLinkBuf(wa, 4)

	// 创建第1个节点，内容是 []byte("123")， realease = nil len=3, cap=3（非内存池分配）
	b1.Append([]byte("123"))
	// 创建第2个字节（内存池分配），填充后剩余1个空闲字节
	n, err := b1.Write([]byte("456"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	// 创建第二个空链表，每个节点最多容纳2个字节
	b2 := NewLinkBuf(wa, 2)
	// 创建第一个节点（内存池分配）, len=2, cap=2
	// 创建第二个节点（内存池分配），len=1, cap=2
	n, err = b2.Write([]byte("567"))
	require.Nil(t, err)
	require.Equal(t, 3, n)
	// 创建第三个节点（非内存池分配），len=2, cap=2
	// 创建第四个字节（内存池不分配），len=0, cap=1 共享第一个字节的内存空间
	b2.Append([]byte("89"))

	// 读取第一个节点，将其放到dirty中
	bts := make([]byte, 2)
	n, err = b2.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, "56", string(bts))
	// 释放了第一个节点到内存池中
	b2.Release()

	b1.Merge(b2)
	bts = make([]byte, 9)
	n, err = b1.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 9, n)
	require.Equal(t, "123456789", string(bts))

	b1.Release()
	wa.MustMallocTimes(1 + 2)
	wa.MustAllFreed()
}

func TestBuf_ReadN(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 创建第1个节点，len=3, cap=4（内存池分配）
	n, err := b.Write([]byte("123"))
	require.Nil(t, err)
	require.Equal(t, 3, n)
	// 创建第2个节点 len=3, cap=3（非内存池分配）
	// 创建第3个字节 len=0, cap=1（内存池不分配）共享第一个字节的内存空间
	b.Append([]byte("456"))

	// 第一个节点读取1个字节
	bts, n := b.ReadN(1)
	require.Equal(t, 1, n)
	require.Equal(t, "1", string(bts))

	// 第一个节点读取剩余的2个字节
	bts, n = b.ReadN(3)
	require.Equal(t, 2, n)
	require.Equal(t, "23", string(bts))

	// 第二个节点读取3个字节
	bts, n = b.ReadN(3)
	require.Equal(t, 3, n)
	require.Equal(t, "456", string(bts))

	bts, n = b.ReadN(3)
	require.Equal(t, 0, n)
	require.Nil(t, bts)

	b.Release()
	wa.MustMallocTimes(1)
	wa.MustAllFreed()
}

func TestBuf_ReadAll(t *testing.T) {
	// 创建空链表，每个节点最多容纳4个字节
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	// 创建第1个节点，len=3, cap=4（内存池分配） 123
	n, err := b.Write([]byte("123"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	// 创建第2个节点 len=2, cap=2（非内存池分配） 45
	// 创建第3个节点 len=0, cap=1 (内存池不分配) 和第一个节点共享内存空间
	b.Append([]byte("45"))

	// 将 6 写入到第三个节点
	// 创建第4个节点，len=4, cap=4（内存池分配）7890
	n, err = b.Write([]byte("67890"))
	require.Nil(t, err)
	require.Equal(t, 5, n)

	bs := b.ReadAll()
	require.Len(t, bs, 4)
	require.Equal(t, "123", string(bs[0]))
	require.Equal(t, "45", string(bs[1]))
	require.Equal(t, "6", string(bs[2]))
	require.Equal(t, "7890", string(bs[3]))

	b.Release()
	wa.MustMallocTimes(2)
	wa.MustAllFreed()
}

func TestBuf_ReadNext(t *testing.T) {
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	n, err := b.Write([]byte("123"))
	require.Nil(t, err)
	require.Equal(t, 3, n)
	b.Append([]byte("45"))
	n, err = b.Write([]byte("678"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	require.Equal(t, "123", string(b.ReadNext()))
	require.Equal(t, "45", string(b.ReadNext()))
	require.Equal(t, "6", string(b.ReadNext()))
	require.Equal(t, "78", string(b.ReadNext()))
	require.Equal(t, "", string(b.ReadNext()))

	b.Release()
	wa.MustMallocTimes(2)
	wa.MustAllFreed()
}

func TestBuf_WriteBigData(t *testing.T) {
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 2)

	n, err := b.Write([]byte("123"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	copy(b.Alloc(3), "456")

	n, err = b.Write([]byte("78"))
	require.Nil(t, err)
	require.Equal(t, 2, n)

	bts := make([]byte, 8)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 8, n)
	require.Equal(t, "12345678", string(bts))

	b.Release()
	wa.MustMallocTimes(2 + 1 + 1)
	wa.MustAllFreed()
}

func TestBuf_PrependAfterRead(t *testing.T) {
	wa := newWrappedAllocator(t, allocator.NewClassAllocator())
	b := NewLinkBuf(wa, 4)

	b.Append([]byte("123"))

	n, err := b.Write([]byte("456"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	bts := make([]byte, 4)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, "1234", string(bts))

	b.Prepend([]byte("34"))
	copy(b.Prelloc(2), "12")

	bts = make([]byte, 6)
	n, err = b.Read(bts)
	require.Nil(t, err)
	require.Equal(t, 6, n)
	require.Equal(t, "123456", string(bts))

	b.Release()
	wa.MustMallocTimes(1 + 1)
	wa.MustAllFreed()
}

func TestBuf_UseBytesAfterRelease(t *testing.T) {
	bytesAllocator := newBytesAllocator()

	b1 := NewLinkBuf(bytesAllocator, 4)
	n, err := b1.Write([]byte("123"))
	require.Nil(t, err)
	require.Equal(t, 3, n)
	b1.Append([]byte("45"))
	n, err = b1.Write([]byte("678"))
	require.Nil(t, err)
	require.Equal(t, 3, n)

	bts := b1.ReadNext()
	require.Equal(t, "123", string(bts))
	b1.Release()

	b2 := NewLinkBuf(bytesAllocator, 4)
	n, err = b2.Write([]byte("1234"))
	require.Nil(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, "123", string(bts), "bts of b1 is not released now")

	b1.ReadNext() // read "45"
	b1.Release()  // bts's underlying buffer is still not released

	n, err = b2.Write([]byte("5678"))
	require.Nil(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, "123", string(bts), "bts of b1 is not released now")

	b1.ReadNext() // read "6"
	b1.Release()  // bts's underlying buffer has been released

	n, err = b2.Write([]byte("5678"))
	require.Nil(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, "567", string(bts), "bts of b1 has been changed by b2")
}
