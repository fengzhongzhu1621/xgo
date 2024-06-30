package pool

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// 连接池关闭错误
	// ErrClosed performs any operation on the closed client will return this error.
	ErrClosed = errors.New("redis: client is closed")

	// ErrPoolTimeout timed out waiting to get a connection from the connection pool.
	ErrPoolTimeout = errors.New("redis: connection pool timeout")
)

// 创建临时对象池，存放单次定时器.
var timers = sync.Pool{
	New: func() interface{} {
		// 创建一个小时定时器
		t := time.NewTimer(time.Hour)
		// 暂停定时器
		// Stop()在这里并不会停止定时器，停止后，Timer不会再被发送，但是Stop不会关闭通道
		t.Stop()
		return t
	},
}

// Stats contains pool state information and accumulated stats.
type Stats struct {
	Hits     uint32 // number of times free connection was found in the pool 空闲连接的命中次数
	Misses   uint32 // number of times free connection was NOT found in the pool	空闲连接没有命中的次数
	Timeouts uint32 // number of times a wait timeout occurred	获取空闲连接的超时的次数

	TotalConns uint32 // number of total connections in the pool	连接池中连接的数量
	IdleConns  uint32 // number of idle connections in the pool		连接池中空闲连接的数量
	StaleConns uint32 // number of stale connections removed from the pool 自动检测关闭的连接数量
}

// 连接池接口.
type Pooler interface {
	NewConn(context.Context) (*Conn, error) // 创建连接
	CloseConn(*Conn) error                  // 关闭连接

	Get(context.Context) (*Conn, error)   // 获得连接
	Put(context.Context, *Conn)           // 添加连接
	Remove(context.Context, *Conn, error) // 删除连接

	Len() int      // 连接的数量
	IdleLen() int  // 空闲连接的数量
	Stats() *Stats // 统计指标

	Close() error // 关闭连接池
}

type Options struct {
	Dialer  func(context.Context) (net.Conn, error) // 连接使用的调用函数
	OnClose func(*Conn) error                       // 连接关闭调用函数

	PoolFIFO           bool          // 默认从空闲连接最后一个获取，为true标识从头获取
	PoolSize           int           // 连接池大小
	MinIdleConns       int           // 最小空闲连接的数量
	MaxConnAge         time.Duration // 空闲连接存活时间
	PoolTimeout        time.Duration // 从连接池获得连接的超时时间
	IdleTimeout        time.Duration // 空闲连接超时时间
	IdleCheckFrequency time.Duration // 空闲连接检查定时器间隔
}

// 定义最近一次错误.
type lastDialErrorWrap struct {
	err error
}

// 连接池实现.
type ConnPool struct {
	opt *Options // 连接配置

	dialErrorsNum uint32 // atomic，连接错误数

	lastDialError atomic.Value // 最近一次错误

	queue chan struct{} // 存放空闲连接的标识，如果存在一个空闲连接则向此队列添加一个空对象

	connsMu      sync.Mutex // 互斥锁
	conns        []*Conn    // 存放所有的连接
	idleConns    []*Conn    // 存放空闲连接
	poolSize     int        // 连接池大小
	idleConnsLen int        // 空闲连接数量

	stats Stats

	_closed  uint32        // atomic，连接池是否关闭标识；1表示连接池已关闭
	closedCh chan struct{} // 关闭连接池时，停止空闲检测
}

var _ Pooler = (*ConnPool)(nil)

// 创建连接池.
func NewConnPool(opt *Options) *ConnPool {
	p := &ConnPool{
		opt: opt, // 连接池配置

		queue:     make(chan struct{}, opt.PoolSize), // 定义管道，为连接池大小
		conns:     make([]*Conn, 0, opt.PoolSize),    // 定义切片，存放所有连接
		idleConns: make([]*Conn, 0, opt.PoolSize),    // 定义切片，存放空闲连接
		closedCh:  make(chan struct{}),               // 定义关闭管道
	}
	// 初始化连接池的大小和空闲连接，预分配连接
	p.connsMu.Lock()
	p.checkMinIdleConns()
	p.connsMu.Unlock()

	// 空闲检查
	if opt.IdleTimeout > 0 && opt.IdleCheckFrequency > 0 {
		go p.reaper(opt.IdleCheckFrequency)
	}

	return p
}

// 初始化连接池的大小和空闲连接，预分配连接.
func (p *ConnPool) checkMinIdleConns() {
	if p.opt.MinIdleConns == 0 {
		return
	}
	for p.poolSize < p.opt.PoolSize && p.idleConnsLen < p.opt.MinIdleConns {
		p.poolSize++
		p.idleConnsLen++

		go func() {
			// 创建一个空闲连接
			err := p.addIdleConn()
			if err != nil && err != ErrClosed {
				// 初始化连接失败
				p.connsMu.Lock()
				p.poolSize--
				p.idleConnsLen--
				p.connsMu.Unlock()
			}
		}()
	}
}

// 创建并添加空闲连接.
func (p *ConnPool) addIdleConn() error {
	// 创建一个新连接
	cn, err := p.dialConn(context.TODO(), true)
	if err != nil {
		return err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// 新的连接创建后，连接池可能被关闭了，所以需要重新判断
	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() {
		_ = cn.Close()
		return ErrClosed
	}
	// 将新的连接加入到连接池
	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	return nil
}

// 创建一个新连接，此连接不再连接池中，但是需要放到p.conns中统一管理.
func (p *ConnPool) NewConn(ctx context.Context) (*Conn, error) {
	return p.newConn(ctx, false)
}

func (p *ConnPool) newConn(ctx context.Context, pooled bool) (*Conn, error) {
	cn, err := p.dialConn(ctx, pooled)
	if err != nil {
		return nil, err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() {
		_ = cn.Close()
		return nil, ErrClosed
	}

	p.conns = append(p.conns, cn)
	if pooled {
		// If pool is full remove the cn on next Put.
		if p.poolSize >= p.opt.PoolSize {
			cn.pooled = false
		} else {
			p.poolSize++
		}
	}

	return cn, nil
}

// 创建连接.
func (p *ConnPool) dialConn(ctx context.Context, pooled bool) (*Conn, error) {
	// 判断连接池是否已经关闭
	if p.closed() {
		return nil, ErrClosed
	}
	// 连接错误数过多
	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.opt.PoolSize) {
		return nil, p.getLastDialError()
	}
	// 执行连接
	netConn, err := p.opt.Dialer(ctx)
	if err != nil {
		// 记录最近一次连接失败错误
		p.setLastDialError(err)
		// 失败则不停的尝试重连
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.opt.PoolSize) {
			go p.tryDial()
		}
		return nil, err
	}
	// 创建一个连接
	cn := NewConn(netConn)
	// 标记连接在连接池中
	cn.pooled = pooled
	return cn, nil
}

// 失败则不停的尝试重连.
func (p *ConnPool) tryDial() {
	for {
		// 判断连接池是否已经关闭
		if p.closed() {
			return
		}
		// 尝试连接
		conn, err := p.opt.Dialer(context.Background())
		if err != nil {
			// 失败重试
			p.setLastDialError(err)
			time.Sleep(time.Second)
			continue
		}
		// 连接成功则恢复连接错误数量，表示可以重新创建连接了
		atomic.StoreUint32(&p.dialErrorsNum, 0)
		// 关闭连接
		_ = conn.Close()
		return
	}
}

// 存储最近一次错误.
func (p *ConnPool) setLastDialError(err error) {
	p.lastDialError.Store(&lastDialErrorWrap{err: err})
}

// 获得最近一次错误.
func (p *ConnPool) getLastDialError() error {
	err, _ := p.lastDialError.Load().(*lastDialErrorWrap)
	if err != nil {
		return err.err
	}
	return nil
}

// 从连接池中获取一个连接
// Get returns existed connection from the pool or creates a new one.
func (p *ConnPool) Get(ctx context.Context) (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}
	// 判断连接池是否有空闲的连接
	if err := p.waitTurn(ctx); err != nil {
		return nil, err
	}

	for {
		// 获得一个空闲连接
		p.connsMu.Lock()
		cn, err := p.popIdle()
		p.connsMu.Unlock()

		if err != nil {
			return nil, err
		}

		if cn == nil {
			break
		}
		// 判断连接是否无效，无效连接关闭
		if p.isStaleConn(cn) {
			_ = p.CloseConn(cn)
			continue
		}
		// 获得可用空闲连接，记录命中次数
		atomic.AddUint32(&p.stats.Hits, 1)
		return cn, nil
	}

	atomic.AddUint32(&p.stats.Misses, 1)

	// 找不到空闲连接就创建一个新的，并加入到连接池
	newcn, err := p.newConn(ctx, true)
	if err != nil {
		p.freeTurn()
		return nil, err
	}

	return newcn, nil
}

// 预先占用一个空闲连接，分为3个来源
// 	1. Get()时从空闲连接中获取
//	2. Get()时如果没有空闲连接，则创建一个新连接
//	3. 空闲检查时定时调用
func (p *ConnPool) getTurn() {
	p.queue <- struct{}{}
}

// 释放预先占用的空闲连接
//	1. 空闲检查定时调用结束时操作
//	2. Put()中调用
func (p *ConnPool) freeTurn() {
	<-p.queue
}

// 判断连接池是否有空闲的连接.
func (p *ConnPool) waitTurn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case p.queue <- struct{}{}:
		return nil
	default:
	}

	// 下面是队列满的处理逻辑，表示连接池中没有空闲连接了

	// 重置单次定时器，设置过期时间
	timer := timers.Get().(*time.Timer)
	timer.Reset(p.opt.PoolTimeout)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return ctx.Err()
	case p.queue <- struct{}{}:
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return nil
	case <-timer.C:
		// 超时返回错误
		timers.Put(timer)
		atomic.AddUint32(&p.stats.Timeouts, 1)
		return ErrPoolTimeout
	}
}

// 获得一个空闲连接.
func (p *ConnPool) popIdle() (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}
	n := len(p.idleConns)
	if n == 0 {
		return nil, nil
	}

	var cn *Conn
	if p.opt.PoolFIFO {
		// 获得第一个空闲连接
		cn = p.idleConns[0]
		copy(p.idleConns, p.idleConns[1:])
		p.idleConns = p.idleConns[:n-1]
	} else {
		// 获得最后一个空闲连接
		idx := n - 1
		cn = p.idleConns[idx]
		p.idleConns = p.idleConns[:idx]
	}
	p.idleConnsLen--
	// 创建空闲连接
	p.checkMinIdleConns()
	return cn, nil
}

// 将一个连接加入到空闲连接列表.
func (p *ConnPool) Put(ctx context.Context, cn *Conn) {
	// 判断连接是否还有未传输完成的数据
	if cn.rd.Buffered() > 0 {
		p.Remove(ctx, cn, BadConnError{})
		return
	}
	// 如果连接是孤立的，则直接删除
	if !cn.pooled {
		p.Remove(ctx, cn, nil)
		return
	}
	// 添加到空闲连接列表
	p.connsMu.Lock()
	p.idleConns = append(p.idleConns, cn)
	p.idleConnsLen++
	p.connsMu.Unlock()
	p.freeTurn()
}

func (p *ConnPool) Remove(_ context.Context, cn *Conn, _ error) {
	p.removeConnWithLock(cn)
	p.freeTurn()
	_ = p.closeConn(cn)
}

// 关闭指定连接.
func (p *ConnPool) CloseConn(cn *Conn) error {
	p.removeConnWithLock(cn)
	return p.closeConn(cn)
}

func (p *ConnPool) removeConnWithLock(cn *Conn) {
	p.connsMu.Lock()
	p.removeConn(cn)
	p.connsMu.Unlock()
}

// 从连接池移除连接.
func (p *ConnPool) removeConn(cn *Conn) {
	for i, c := range p.conns {
		if c == cn {
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			if cn.pooled {
				p.poolSize--
				// 添加新的空闲连接
				p.checkMinIdleConns()
			}
			return
		}
	}
}

func (p *ConnPool) closeConn(cn *Conn) error {
	if p.opt.OnClose != nil {
		_ = p.opt.OnClose(cn)
	}
	return cn.Close()
}

// Len returns total number of connections.
func (p *ConnPool) Len() int {
	p.connsMu.Lock()
	n := len(p.conns)
	p.connsMu.Unlock()
	return n
}

// IdleLen returns number of idle connections.
func (p *ConnPool) IdleLen() int {
	p.connsMu.Lock()
	n := p.idleConnsLen
	p.connsMu.Unlock()
	return n
}

func (p *ConnPool) Stats() *Stats {
	idleLen := p.IdleLen()
	return &Stats{
		Hits:     atomic.LoadUint32(&p.stats.Hits),
		Misses:   atomic.LoadUint32(&p.stats.Misses),
		Timeouts: atomic.LoadUint32(&p.stats.Timeouts),

		TotalConns: uint32(p.Len()),
		IdleConns:  uint32(idleLen),
		StaleConns: atomic.LoadUint32(&p.stats.StaleConns),
	}
}

// 判断连接池是否已经关闭.
func (p *ConnPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

func (p *ConnPool) Filter(fn func(*Conn) bool) error {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	var firstErr error
	for _, cn := range p.conns {
		if fn(cn) {
			if err := p.closeConn(cn); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// 关闭连接池.
func (p *ConnPool) Close() error {
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return ErrClosed
	}
	close(p.closedCh)

	var firstErr error
	p.connsMu.Lock()
	for _, cn := range p.conns {
		if err := p.closeConn(cn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.poolSize = 0
	p.idleConns = nil
	p.idleConnsLen = 0
	p.connsMu.Unlock()

	return firstErr
}

// 定时检查空闲的连接，如果过期则关闭.
func (p *ConnPool) reaper(frequency time.Duration) {
	// 创建定时器
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// It is possible that ticker and closedCh arrive together,
			// and select pseudo-randomly pick ticker case, we double
			// check here to prevent being executed after closed.
			if p.closed() {
				return
			}
			// 关闭过期的所有空闲连接
			_, err := p.ReapStaleConns()
			if err != nil {
				continue
			}
		case <-p.closedCh:
			// 关闭连接池时，停止空闲检测
			return
		}
	}
}

// 关闭过期的所有空闲连接.
func (p *ConnPool) ReapStaleConns() (int, error) {
	var n int
	for {
		// 空闲检查和Get()操作共用并发，即如果没有空闲连接，queue队列满，则空闲检查等待，减少性能消耗
		p.getTurn()

		p.connsMu.Lock()
		// 从空闲连接列表中判断第一个连接是否是无效连接，如果无效则删除
		cn := p.reapStaleConn()
		p.connsMu.Unlock()

		p.freeTurn()
		// 关闭无效连接
		if cn != nil {
			_ = p.closeConn(cn)
			n++
		} else {
			break
		}
	}
	// 记录空闲连接的数量
	atomic.AddUint32(&p.stats.StaleConns, uint32(n))
	return n, nil
}

// 从空闲连接列表中判断第一个连接是否是无效连接，如果无效则删除.
func (p *ConnPool) reapStaleConn() *Conn {
	if len(p.idleConns) == 0 {
		return nil
	}
	// 获得一个空闲连接
	cn := p.idleConns[0]
	// 判断连接是否是无效的
	if !p.isStaleConn(cn) {
		return nil
	}
	// 删除第一个空闲连接
	p.idleConns = append(p.idleConns[:0], p.idleConns[1:]...)
	p.idleConnsLen--
	p.removeConn(cn)

	return cn
}

// 判断连接是否是无效的.
func (p *ConnPool) isStaleConn(cn *Conn) bool {
	if p.opt.IdleTimeout == 0 && p.opt.MaxConnAge == 0 {
		return false
	}

	now := time.Now()
	if p.opt.IdleTimeout > 0 && now.Sub(cn.UsedAt()) >= p.opt.IdleTimeout {
		return true
	}
	if p.opt.MaxConnAge > 0 && now.Sub(cn.createdAt) >= p.opt.MaxConnAge {
		return true
	}

	return false
}
