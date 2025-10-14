package bufferpool

import "bytes"

// IBufferPool 缓冲区池接口定义
type IBufferPool interface {
	Put(*bytes.Buffer)
	Get() *bytes.Buffer
}
