// Package allocator implements byte slice pooling management
// to reduce the pressure of memory allocation.
//
// 该包实现了字节切片池化管理，用于减少内存分配压力。

package allocator

import (
	"fmt"
	"sync"

	"github.com/fengzhongzhu1621/xgo/numpy/math"
)

// maxPowerToRoundUpInt 定义了最大支持的2的幂指数
// 63表示支持的最大size是1<<63，这是64位系统上最大的可能值
const maxPowerToRoundUpInt = 63

var _ IAllocator = (*ClassAllocator)(nil)

// ClassAllocator 是一个字节池，管理的字节大小满足2^n
//
// pools: 存储不同大小的sync.Pool数组，索引i对应大小2^i
type ClassAllocator struct {
	pools [maxPowerToRoundUpInt]*sync.Pool
}

// defaultAllocator 是默认的全局分配器实例
// 使用NewClassAllocator初始化，提供线程安全的字节池
var defaultAllocator = NewClassAllocator()

// NewClassAllocator 创建一个新的ClassAllocator实例
//
// returns: 初始化好的ClassAllocator指针
func NewClassAllocator() *ClassAllocator {
	// 创建包含maxPowerToRoundUpInt个sync.Pool的数组
	// 每个Pool对应一个特定大小的字节池(1<<0, 1<<1, ..., 1<<62)
	// 1<<62 的字节值：1<<62 表示将数字1左移62位，计算结果为 4,611,686,018,427,387,904 字节
	// 1 GB = 2^30 字节，2^62 = 2^32GB = 4,294,967,296 GB = 4,194,304 TB（太字节）= 4,096 PB（拍字节) = 4 EB（艾字节）
	var pools [maxPowerToRoundUpInt]*sync.Pool

	// 初始化每个Pool
	for i := range pools {
		size := 1 << i // 计算该Pool管理的切片大小(2^i)
		// 注意：下面的代码不会分配内存，仅在使用 Get() 方法后才会分配内存
		//
		// 为每个Pool创建一个新的字节切片，并将其放入Pool中
		// 当Pool为空时，调用New函数创建新的字节切片
		// 当Pool不为空时，从Pool中获取一个字节切片
		// 注意：New函数返回的字节切片长度为0，但容量为size
		pools[i] = &sync.Pool{
			New: func() interface{} {
				// 当Pool为空时，创建新的字节切片
				return make([]byte, size)
			},
		}
	}

	// 返回初始化好的ClassAllocator
	return &ClassAllocator{pools: pools}
}

// Malloc 从池中获取一个[]byte，第二个返回值用于Free操作
//
// 注意：虽然可以使用第一个返回值来释放字节，但这会导致额外的堆分配
// 详见 https://github.com/golang/go/issues/8618
//
// size: 需要分配的字节大小
// returns:
//   - 第一个返回值: 实际可用的字节切片(长度=size)
//   - 第二个返回值: 内部使用的原始分配对象(用于释放)
func (a *ClassAllocator) Malloc(size int) ([]byte, interface{}) {
	// 检查size有效性
	if size <= 0 {
		panic(fmt.Sprintf("invalid alloc size %d", size))
	}

	// 计算所需大小对应的Pool索引
	// 找到最小的 power 使得 2^power >= size
	power := math.PowerToRoundUp(size)

	// 从对应Pool获取或创建字节切片（此时会分配内存）
	v := a.pools[power].Get()

	// 返回切片和原始对象（可能比实际需要的大）
	return v.([]byte)[:size], v
}

// Free 将字节释放回池中
//
// bts: 由Malloc返回的第二个参数，包含原始分配对象
func (a *ClassAllocator) Free(bts interface{}) {
	// 获取切片的实际容量
	cap := cap(bts.([]byte))

	// 检查容量有效性
	if cap == 0 {
		panic("free an empty bytes")
	}

	// 计算容量对应的Pool索引
	power := math.PowerToRoundUp(cap)

	// 验证容量是否为2的幂
	if 1<<power != cap {
		panic(fmt.Sprintf("cap %d of bts must be power of two", cap))
	}

	// 将对象放回对应Pool
	a.pools[power].Put(bts)
}
