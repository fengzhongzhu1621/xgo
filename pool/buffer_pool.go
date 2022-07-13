package pool

import (
	"bytes"
	"sync"
)

var (
	bufferPool BufferPool
)

type BufferPool interface {
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
}

var _ BufferPool = (*defaultPool)(nil)

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
func SetBufferPool(bp BufferPool) {
	bufferPool = bp
}

// bytes.buffer是一个缓冲byte类型的缓冲器，使用 sync.Pool 减少 bytes.Buffer 创建的成本
// 保存和复用临时对象，减少内存分配，降低 GC 压力。
func init() {
	SetBufferPool(&defaultPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	})
}

func GetDefaultBufferPool() BufferPool {
	return bufferPool
}
