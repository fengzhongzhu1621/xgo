package pool

import (
	"bytes"
	"sync"
)

// 确保 defaultPool 类型实现了 IBufferPool 接口（编译时检查）
var _ IBufferPool = (*defaultPool)(nil)

// DefaultBufferPool 全局默认的缓冲区池实例
var DefaultBufferPool IBufferPool

// IBufferPool 缓冲区池接口定义
type IBufferPool interface {
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
}

type defaultPool struct {
	pool *sync.Pool
}

func (p *defaultPool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

func (p *defaultPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// SetBufferPool allows to replace the default logrus buffer pool
// to better meets the specific needs of an application.
// 允许替换默认的日志缓冲区池
func SetBufferPool(bp IBufferPool) {
	DefaultBufferPool = bp
}

// bytes.buffer是一个缓冲byte类型的缓冲器，使用 sync.Pool 减少 bytes.Buffer 创建的成本
// 保存和复用临时对象，减少内存分配，降低 GC 压力。
func init() {
	SetBufferPool(&defaultPool{
		pool: &sync.Pool{
			// New 函数在池中没有可用对象时创建新的 bytes.Buffer
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	})
}
