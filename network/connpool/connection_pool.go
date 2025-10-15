package connpool

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/logging"
)

const (
	defaultDialTimeout     = 200 * time.Millisecond
	defaultIdleTimeout     = 50 * time.Second
	defaultMaxIdle         = 65536
	defaultCheckInterval   = 3 * time.Second
	defaultPoolIdleTimeout = 2 * defaultIdleTimeout
)

// ConnectionPool is the connection pool.
// 连接池需要具备以下功能：
// 提供可用连接，包括创建新连接和复用空闲连接；
// 回收上层使用过的连接作为空闲连接管理；
// 对连接池中空闲连接的管理能力，包括复用连接的选择策略，空闲连接的健康监测等；
// 根据用户配置调整连接池运行参数
type ConnectionPool struct {
	Dial func(context.Context) (net.Conn, error) // initialize the connection.
	// 业务的突发流量可能会导致大量新连接建立，创建连接是一个比较耗时的操作，可能导致请求超时。提前创建部分空闲连接可以起到预热效果。连接池在创建时创建 MinIdle 个连接备用。
	// 需要根据 MinIdle 预热空闲连接
	// MinIdle 是 ConnectionPool 维持的最小空闲连接，在初始化和周期检查中进行补充。 用户获取连接时，首先从空闲连接中获取，若没有空闲连接才会重新创建。当用户完成请求后，将连接归还给 ConnectionPool
	MinIdle     int           // Minimum number of idle connections.
	MaxIdle     int           // Maximum number of idle connections, 0 means no limit.
	MaxActive   int           // Maximum number of active connections, 0 means no limit.
	IdleTimeout time.Duration // idle connection timeout.
	// Whether to wait when the maximum number of active connections is reached.
	Wait               bool
	MaxConnLifetime    time.Duration  // Maximum lifetime of the connection.
	mu                 sync.Mutex     // Control concurrent locks. 并发控制
	checker            IHealthChecker // Idle connection health check function.
	closed             bool           // Whether the connection pool has been closed. 关闭标识
	tokenChan          chan struct{}  // control concurrency by applying token.
	idleSize           int            // idle connections size. 空闲连接数量
	idle               connList       // idle connection list. 空闲链接链表
	framerBuilder      codec.IFramerBuilder
	forceClosed        bool // Force close the connection, suitable for streaming scenarios.
	PushIdleConnToTail bool // connection to ip will be push tail when ConnectionPool.put method is called.
	// customReader creates a reader encapsulating the underlying connection.
	customReader    func(io.Reader) io.Reader
	onCloseFunc     func()        // execute when checker goroutine judge the connection_pool is useless.
	used            int32         // size of connections used by user, atomic.
	lastGetTime     int64         // last get connection millisecond timestamp, atomic. 上一次取出连接的时间
	poolIdleTimeout time.Duration // pool idle timeout.
}

// keepMinIdles 保持最小空闲连接数
// 当空闲连接数小于 MinIdle 时，创建新的空闲连接进行补充
func (p *ConnectionPool) keepMinIdles() {
	// 计算需要补充的连接数
	p.mu.Lock()
	count := p.MinIdle - p.idleSize
	if count > 0 {
		p.idleSize += count // 预增加空闲连接数
	}
	p.mu.Unlock()

	// 并发创建空闲连接
	for i := 0; i < count; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
			defer cancel()

			// 创建空闲连接
			if err := p.addIdleConn(ctx); err != nil {
				// 如果创建连接失败，减少空闲连接计数
				p.mu.Lock()
				p.idleSize--
				p.mu.Unlock()
			}
		}()
	}
}

// addIdleConn 添加空闲连接到连接池
// 参数:
//
//	ctx: 上下文对象
//
// 返回:
//
//	error: 添加过程中发生的错误
func (p *ConnectionPool) addIdleConn(ctx context.Context) error {
	// 判断连接池是否可用
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return ErrPoolClosed // 连接池已关闭，返回错误
	}
	p.mu.Unlock()

	// 建立新连接
	c, err := p.dial(ctx)
	if err != nil {
		return err // 拨号失败，返回错误
	}

	// 创建连接池连接对象
	pc := p.newPoolConn(c)
	p.mu.Lock()
	if p.closed {
		// 如果连接池在拨号过程中被关闭，关闭连接
		pc.closed = true
		pc.Conn.Close()
	} else {
		// 设置连接时间戳并添加到空闲列表
		// 连接池有 FIFO 和 LIFO 两种策略进行空闲连接的选择和淘汰，通过 PushIdleConnToTail 控制，应该根据业务的实际特点选择合适的管理策略。
		//
		// fifo，保证各个连接均匀使用，但是当调用方请求频率不高，但是恰巧每次能在连接空闲条件命中之前来一个请求，就会导致各个连接无法被释放，此时维持这么多的连接数是多余的。
		// lifo, 优先采用栈顶连接，栈底连接不频繁使用会优先淘汰。
		pc.t = time.Now()
		if !p.PushIdleConnToTail {
			p.idle.pushHead(pc) // 添加到列表头部
		} else {
			p.idle.pushTail(pc) // 添加到列表尾部
		}
	}
	p.mu.Unlock()

	return nil
}

// Close 关闭连接池并释放所有连接
// 返回:
//
//	error: 关闭过程中发生的错误
func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil // 连接池已经关闭，直接返回
	}

	p.closed = true
	p.idle.count = 0
	p.idleSize = 0
	pc := p.idle.head
	p.idle.head, p.idle.tail = nil, nil // 清空链表
	p.mu.Unlock()

	// 遍历链表关闭所有空闲连接
	for ; pc != nil; pc = pc.next {
		pc.Conn.Close()
		pc.closed = true
	}

	return nil
}

// Get 从连接池获取连接
// 参数:
//
//	ctx: 上下文对象
//
// 返回:
//
//	*PoolConn: 连接池连接对象
//	error: 获取过程中发生的错误
func (p *ConnectionPool) Get(ctx context.Context) (*PoolConn, error) {
	var (
		pc  *PoolConn
		err error
	)
	if pc, err = p.get(ctx); err != nil {
		return nil, err
	}
	return pc, nil
}

// get 从连接池获取连接的核心逻辑
// 参数:
//
//	ctx: 上下文对象
//
// 返回:
//
//	*PoolConn: 连接池连接对象
//	error: 获取过程中发生的错误
func (p *ConnectionPool) get(ctx context.Context) (*PoolConn, error) {
	// 获取令牌控制并发
	if err := p.getToken(ctx); err != nil {
		return nil, err
	}

	// 更新最后获取连接时间和使用计数
	atomic.StoreInt64(&p.lastGetTime, time.Now().UnixMilli())
	atomic.AddInt32(&p.used, 1)

	// 优先尝试获取空闲连接
	if pc := p.getIdleConn(); pc != nil {
		return pc, nil
	}

	// 如果没有空闲连接，创建新连接
	pc, err := p.getNewConn(ctx)
	if err != nil {
		p.freeToken() // 创建失败，释放令牌
		return nil, err
	}
	return pc, nil
}

// if p.Wait is True, return err when timeout.
// if p.Wait is False, return err when token empty immediately.
// token 是一个用于并发控制的 ch, 其缓冲长度根据 MaxActive 设置，代表用户可以同时使用 MaxActive 个连接，
// 当活跃连接被归还连接池或关闭时，归还 token.
// 如果设置 Wait=True, 会在获取不到 token 时等待直到超时返回，
// 如果设置 Wait=False, 会在获取不到 token 时直接返回 ErrPoolLimit。
// 成功获取 token 后，优先从 idle list 中获取空闲连接，如果失败则新创建连接返回。
func (p *ConnectionPool) getToken(ctx context.Context) error {
	if p.MaxActive <= 0 {
		return nil
	}

	if p.Wait {
		select {
		case p.tokenChan <- struct{}{}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	} else {
		select {
		case p.tokenChan <- struct{}{}:
			return nil
		default:
			// 令牌桶满（连接已满）
			return ErrPoolLimit
		}
	}
}

func (p *ConnectionPool) freeToken() {
	if p.MaxActive <= 0 {
		return
	}
	<-p.tokenChan
}

// getIdleConn 从空闲列表中获取空闲连接
// 返回:
//
//	*PoolConn: 连接池连接对象
//	error: 获取过程中发生的错误
func (p *ConnectionPool) getIdleConn() *PoolConn {
	p.mu.Lock()
	for p.idle.head != nil {
		pc := p.idle.head
		p.idle.popHead() // 出栈
		p.idleSize--
		p.mu.Unlock()

		// 空闲连接检查
		if p.checker(pc, true) {
			return pc
		}

		// 如果检查失败，关闭连接并继续检查下一个空闲连接
		pc.Conn.Close()
		pc.closed = true
		p.mu.Lock()
	}
	p.mu.Unlock()

	return nil
}

// getNewConn 创建新连接
func (p *ConnectionPool) getNewConn(ctx context.Context) (*PoolConn, error) {
	// If the connection pool has been closed, return an error directly.
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	p.mu.Unlock()

	c, err := p.dial(ctx)
	if err != nil {
		return nil, err
	}

	return p.newPoolConn(c), nil
}

// newPoolConn 创建新的连接池连接对象
// 参数:
//
//	c: 原始的网络连接
//
// 返回:
//
//	*PoolConn: 新创建的连接池连接对象
func (p *ConnectionPool) newPoolConn(c net.Conn) *PoolConn {
	// 初始化连接池连接对象
	pc := &PoolConn{
		Conn:       c,             // 原始网络连接
		created:    time.Now(),    // 连接创建时间
		pool:       p,             // 所属连接池
		forceClose: p.forceClosed, // 是否强制关闭（流式场景适用）
		inPool:     false,         // 初始状态不在池中
	}

	// 如果设置了帧构建器，创建帧读取器
	if p.framerBuilder != nil {
		// 使用自定义读取器包装原始连接，然后创建帧读取器
		pc.fr = p.framerBuilder.New(p.customReader(pc))
		// 检查帧读取器是否安全，决定是否需要复制帧数据
		pc.copyFrame = !codec.IsSafeFramer(pc.fr)
	}

	return pc
}

func (p *ConnectionPool) checkHealthOnce() {
	p.mu.Lock()
	n := p.idle.count
	for i := 0; i < n && p.idle.head != nil; i++ {
		pc := p.idle.head
		p.idle.popHead()
		p.idleSize--
		p.mu.Unlock()

		// 检查连接是否健康
		if p.checker(pc, false) {
			p.mu.Lock()
			p.idleSize++
			p.idle.pushTail(pc)
		} else {
			pc.Conn.Close()
			pc.closed = true
			p.mu.Lock()
		}
	}
	p.mu.Unlock()
}

// checkRoutine 检查连接池健康状态的定时任务，巡检空闲连接，检查连接是否健康，
// 如果连接不健康，则关闭连接并从连接池中移除，如果连接池空闲连接数小于最小空闲连接数，则创建新的连接并放到连接池中。
func (p *ConnectionPool) checkRoutine(interval time.Duration) {
	for {
		time.Sleep(interval)
		p.mu.Lock()
		closed := p.closed
		p.mu.Unlock()
		if closed {
			return
		}

		p.checkHealthOnce()

		if p.checkPoolIdleTimeout() {
			return // 连接池被关闭，退出检查
		}

		// Check if the minimum number of idle connections is met.
		p.checkMinIdle()
	}
}

func (p *ConnectionPool) checkMinIdle() {
	if p.MinIdle <= 0 {
		return
	}
	// 保持最小空闲连接数
	p.keepMinIdles()
}

// checkPoolIdleTimeout 检查连接池是否处于无用状态
// 当连接池长时间未被使用且没有活跃连接时，自动关闭连接池以释放资源
// 返回:
//
//	bool: 如果连接池被关闭则返回 true，否则返回 false
func (p *ConnectionPool) checkPoolIdleTimeout() bool {
	p.mu.Lock()
	// 获取最后使用连接池的时间戳
	lastGetTime := atomic.LoadInt64(&p.lastGetTime)
	// 如果从未使用过连接池或未设置超时时间，直接返回
	if lastGetTime == 0 || p.poolIdleTimeout == 0 {
		p.mu.Unlock()
		return false
	}

	// 检查是否满足关闭条件：
	// 1. 当前时间与最后使用时间差超过池空闲超时时间
	// 2. 设置了关闭回调函数
	// 3. 当前没有活跃连接在使用
	if time.Now().UnixMilli()-lastGetTime > p.poolIdleTimeout.Milliseconds() &&
		p.onCloseFunc != nil && atomic.LoadInt32(&p.used) == 0 {
		p.mu.Unlock()
		// 执行关闭前的回调函数
		p.onCloseFunc()
		// 关闭连接池
		if err := p.Close(); err != nil {
			logging.Errorf("failed to close ConnectionPool, error: %v", err)
		}
		return true
	}
	p.mu.Unlock()
	return false
}

// RegisterChecker registers the idle connection check method.
func (p *ConnectionPool) RegisterChecker(interval time.Duration, checker IHealthChecker) {
	if interval <= 0 || checker == nil {
		return
	}
	p.mu.Lock()
	p.checker = checker
	p.mu.Unlock()
	go p.checkRoutine(interval)
}

// defaultChecker is the default idle connection check method,
// returning true means the connection is available normally.
// 检查协程：健康检查扫描 idle 链表，如果未通过安全检查则将连接直接关闭，首先检查连接是否正常，
// 然后检查是否到达 IdleTimeout 和 MaxConnLifetime. 可以使用 WithHealthChecker 自定义健康检查策略。
// 除周期性的检查空闲连接，在每次从 idle list 获取空闲连接是都会检查，此时将 isFast 设为 true, 只进行连接存活确认
//
// 连接池检测连接空闲的时间，通常也要做成可配置化的，目的是为了与 server 端配合（尤其要考虑不同框架的场景），如果配合的不好，也会出问题。
// 比如 pool 空闲连接检测时间是 1min，server 也是 1min，
// 可能会存在这样的情景，就是 server 端密集关闭空闲连接的时候，client 端还没检测到，发送数据的时候发现大量失败，而不得不通过上层重试解决。
// 比较好的做法是，server 空闲连接检测时长设置为 pool 空闲连接检测时长大一些，尽量让 client 端主动关闭连接，避免取出的连接被 server 关闭而不自知。
//
// 空闲连接数量检查 同 KeepMinIdles, 周期性的将空闲连接数补充到 MinIdle 个。
// ConnectionPool 空闲检查 transport 不会主动关闭 ConnectionPool, 会导致后台检查协程空转。通过设置 poolIdleTimeout, 周期性检查在此时间内用户使用连接数为 0, 来保证长时间未使用的 ConnectionPool 自动关闭。
func (p *ConnectionPool) defaultChecker(pc *PoolConn, isFast bool) bool {
	// Check whether the connection status is abnormal:
	// closed, network exception or sticky packet processing exception.
	if pc.isRemoteError(isFast) {
		return false
	}
	// Based on performance considerations, the quick check only does the RemoteErr check.
	if isFast {
		return true
	}
	// Check if the connection has exceeded the maximum idle time, if so close the connection.
	if p.IdleTimeout > 0 && pc.t.Add(p.IdleTimeout).Before(time.Now()) {
		return false
	}
	// Check if the connection is still alive.
	if p.MaxConnLifetime > 0 && pc.created.Add(p.MaxConnLifetime).Before(time.Now()) {
		return false
	}
	return true
}

// dial establishes a connection.
func (p *ConnectionPool) dial(ctx context.Context) (net.Conn, error) {
	if p.Dial != nil {
		return p.Dial(ctx)
	}
	return nil, errors.New("must pass Dial to pool")
}

// put 尝试将连接释放回连接池
// 根据连接状态决定是放回空闲列表还是直接关闭
// 参数:
//
//	pc: 要释放的连接池连接对象
//	forceClose: 是否强制关闭连接（当连接读写失败时为 true）
//
// 返回:
//
//	error: 释放过程中发生的错误
func (p *ConnectionPool) put(pc *PoolConn, forceClose bool) error {
	// 如果连接已经关闭，直接返回
	if pc.closed {
		return nil
	}

	p.mu.Lock()
	// 检查连接池是否未关闭且不需要强制关闭连接
	if !p.closed && !forceClose {
		// 更新连接的时间戳
		pc.t = time.Now()

		// 根据配置决定将连接添加到空闲列表的头部还是尾部
		if !p.PushIdleConnToTail {
			p.idle.pushHead(pc) // 添加到头部（栈顶）
		} else {
			p.idle.pushTail(pc) // 添加到尾部（队列尾）
		}

		// 如果空闲连接数已达到最大值，移除最旧的连接
		if p.idleSize >= p.MaxIdle {
			pc = p.idle.tail // 获取最旧的连接
			p.idle.popTail() // 从尾部移除
		} else {
			p.idleSize++ // 增加空闲连接计数
			pc = nil     // 连接已成功放回池中，不需要关闭
		}
	}
	p.mu.Unlock()

	// 如果连接需要被关闭（连接池已关闭、强制关闭或超过最大空闲数）
	if pc != nil {
		pc.closed = true
		pc.Conn.Close() // 关闭底层连接
	}

	// 释放令牌，允许其他请求获取连接
	p.freeToken()
	// 减少活跃连接计数
	atomic.AddInt32(&p.used, -1)

	return nil
}
