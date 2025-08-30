package writev

// QuitHandler Buffer goroutine exit handler.
// 队列满时是否丢弃请求
type QuitHandler func(*Buffer)

// Options is the Buffer configuration.
// Options 是缓冲区的配置结构
type Options struct {
	handler    QuitHandler // Set the goroutine to exit the cleanup function. 设置goroutine退出时的清理函数
	bufferSize int         // Set the length of each connection request queue. 设置每个连接请求队列的长度
	dropFull   bool        // Whether the queue is full or not. 队列满时是否丢弃请求
}

// Option optional parameter.
// Option 是可选参数类型，用于配置Options
type Option func(*Options)

// WithQuitHandler returns an Option which sets the Buffer goroutine exit handler.
// 设置goroutine退出时的清理函数
func WithQuitHandler(handler QuitHandler) Option {
	return func(o *Options) {
		o.handler = handler
	}
}

// WithBufferSize returns an Option which sets the length of each connection request queue.
// 设置每个连接请求队列的长度
func WithBufferSize(size int) Option {
	return func(opts *Options) {
		opts.bufferSize = size // 设置缓冲区大小
	}
}

// WithDropFull returns an Option which sets whether to drop the request when the queue is full.
// 设置队列满时是否丢弃请求
func WithDropFull(drop bool) Option {
	return func(opts *Options) {
		opts.dropFull = drop
	}
}
