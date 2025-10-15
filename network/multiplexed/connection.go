package multiplexed

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	queue "github.com/fengzhongzhu1621/xgo/collections/queue/listqueue"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fengzhongzhu1621/xgo/network/dial"
	"github.com/fengzhongzhu1621/xgo/pool/packetbuffer"
)

const (
	maxBufferSize = 65535
)

// The following needs to be variables according to some test cases.
var (
	initialBackoff    = 5 * time.Millisecond
	maxBackoff        = 50 * time.Millisecond
	maxReconnectCount = 10
	// reconnectCountResetInterval is twice the expected total reconnect backoff time,
	// i.e. 2 * \sum_{i=1}^{maxReconnectCount}(i*initialBackoff).
	reconnectCountResetInterval = 5 * time.Millisecond * (1 + 10) * 10
)

// Connection 表示底层TCP/UDP连接，支持多路复用
type Connection struct {
	err                 error       // 连接错误状态
	address             string      // 目标地址
	network             string      // 网络协议类型
	enableIdleRemove    bool        // 是否启用空闲连接移除
	destroy             func()      // 连接销毁回调函数
	connsSubIdle        func()      // 空闲连接数减少回调
	connsAddIdle        func()      // 空闲连接数增加回调
	connsNeedIdleRemove func() bool // 检查是否需要移除空闲连接的回调

	// 重连相关字段
	reconnectCount    int       // 当前重连次数
	lastReconnectTime time.Time // 最后一次重连时间

	// mu 保护虚拟连接映射、空闲状态和连接关闭过程的并发安全
	mu       sync.RWMutex
	virConns map[uint32]*VirtualConnection // 虚拟连接映射表
	isIdle   bool                          // 连接是否空闲

	fp          IFrameParser  // 帧解析器
	done        chan struct{} // 底层连接关闭时关闭的通道
	writeBuffer chan []byte   // 写入缓冲区
	dropFull    bool          // 缓冲区满时是否丢弃数据
	maxVirConns int           // 最大虚拟连接数

	// UDP专用字段
	packetBuffer *packetbuffer.PacketBuffer // UDP数据包缓冲区
	addr         *net.UDPAddr               // UDP目标地址
	packetConn   net.PacketConn             // 底层UDP连接

	// TCP/Unix流专用字段
	conn       net.Conn          // 底层TCP连接
	connLocker sync.RWMutex      // 连接读写锁
	dialOpts   *dial.DialOptions // 拨号选项
	isStream   bool              // 是否为流式连接
	closed     bool              // 连接是否已关闭
}

// setRawConn 设置底层原始连接
// 使用写锁保护连接的设置
// 参数:
//
//	conn: 要设置的网络连接
func (c *Connection) setRawConn(conn net.Conn) {
	c.connLocker.Lock()
	defer c.connLocker.Unlock()

	c.conn = conn
}

// getRawConn 获取底层原始连接
// 使用读锁保护连接的读取
// 返回:
//
//	net.Conn: 底层网络连接
func (c *Connection) getRawConn() net.Conn {
	c.connLocker.RLock()
	defer c.connLocker.RUnlock()
	return c.conn
}

// newVirConn 创建新的虚拟连接
// 参数:
//
//	ctx: 上下文
//	virConnID: 虚拟连接ID
//
// 返回:
//
//	*VirtualConnection: 新创建的虚拟连接
func (c *Connection) newVirConn(ctx context.Context, virConnID uint32) *VirtualConnection {
	ctx, cancel := context.WithCancel(ctx)
	vc := &VirtualConnection{
		id:         virConnID,
		conn:       c,
		ctx:        ctx,
		cancelFunc: cancel,
		recvQueue:  queue.New[[]byte](ctx.Done()),
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// 如果连接建立失败或正在重连，直接关闭虚拟连接
	if c.closed {
		vc.cancel(c.err)
	}
	// 考虑请求ID溢出或上层请求ID重复的情况，
	// 需要先读取并检查请求ID是否已存在，如果存在则向原虚拟连接返回错误
	if prevConn, ok := c.virConns[virConnID]; ok {
		prevConn.cancel(ErrDupRequestID)
	}
	c.virConns[virConnID] = vc
	if c.isIdle {
		c.isIdle = false
		c.connsSubIdle()
	}
	return vc
}

// send 发送数据到写入缓冲区
// 参数:
//
//	b: 要发送的数据字节切片
//
// 返回:
//
//	error: 发送过程中发生的错误
func (c *Connection) send(b []byte) error {
	// 如果设置了dropfull，队列满时丢弃数据
	if c.dropFull {
		select {
		case c.writeBuffer <- b:
			return nil
		default:
			return ErrSendQueueFull
		}
	}
	select {
	case c.writeBuffer <- b:
		return nil
	case <-c.done:
		return c.err
	}
}

// writeAll 写入所有数据到连接
// 根据连接类型调用相应的写入方法
// 参数:
//
//	b: 要写入的数据字节切片
//
// 返回:
//
//	error: 写入过程中发生的错误
func (c *Connection) writeAll(b []byte) error {
	if c.isStream {
		return c.writeTCP(b)
	}
	return c.writeUDP(b)
}

// writeUDP 写入数据到UDP连接
// 参数:
//
//	b: 要写入的数据字节切片
//
// 返回:
//
//	error: 写入过程中发生的错误
func (c *Connection) writeUDP(b []byte) error {
	num, err := c.packetConn.WriteTo(b, c.addr)
	if err != nil {
		return err
	}
	if num != len(b) {
		return ErrWriteNotFinished
	}
	return nil
}

// writeTCP 写入数据到TCP连接
// 确保所有数据都被写入，使用循环直到所有数据发送完成
// 参数:
//
//	b: 要写入的数据字节切片
//
// 返回:
//
//	error: 写入过程中发生的错误
func (c *Connection) writeTCP(b []byte) error {
	var sentNum, num int
	var err error
	conn := c.getRawConn()
	for sentNum < len(b) {
		num, err = conn.Write(b[sentNum:])
		if err != nil {
			return err
		}
		sentNum += num
	}
	return nil
}

// close 关闭连接
// 根据连接类型调用相应的关闭方法
// 参数:
//
//	lastErr: 关闭的原因错误
//	reconnect: 是否尝试重连
func (c *Connection) close(lastErr error, reconnect bool) {
	if c.isStream {
		c.closeTCP(lastErr, reconnect)
		return
	}
	c.closeUDP(lastErr)
}

// closeUDP 关闭UDP连接
// 参数:
//
//	lastErr: 关闭的原因错误
func (c *Connection) closeUDP(lastErr error) {
	c.destroy()
	c.err = lastErr
	close(c.done)

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, vc := range c.virConns {
		vc.cancel(lastErr)
	}
}

// closeTCP 关闭TCP连接
// 参数:
//
//	lastErr: 关闭的原因错误
//	reconnect: 是否尝试重连
func (c *Connection) closeTCP(lastErr error, reconnect bool) {
	if lastErr == nil {
		return
	}
	if needDestroy := c.doClose(lastErr, reconnect); needDestroy {
		c.destroy()
	}
}

// doClose 执行实际的连接关闭操作
// 参数:
//
//	lastErr: 关闭的原因错误
//	reconnect: 是否尝试重连
//
// 返回:
//
//	bool: 是否需要销毁连接
func (c *Connection) doClose(lastErr error, reconnect bool) (needDestroy bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 不要使用c.err != nil来判断，重连不会清除错误
	if c.closed {
		return false
	}
	c.closed = true
	c.err = lastErr

	// 当关闭`c.done`通道时，所有Read操作都会返回错误，
	// 因此应该清理所有现有连接，避免内存泄漏
	for _, vc := range c.virConns {
		vc.cancel(lastErr)
	}
	c.virConns = make(map[uint32]*VirtualConnection)
	close(c.done)
	if conn := c.getRawConn(); conn != nil {
		conn.Close()
	}
	if reconnect && c.doReconnectBackoff() {
		return !c.reconnect()
	}
	return true
}

// reconnect 尝试重连操作
// 返回:
//
//	bool: 重连是否成功
func (c *Connection) reconnect() (success bool) {
	for {
		conn, err := dial.TryConnect(c.dialOpts)
		if err != nil {
			logging.Tracef("reconnect fail: %+v", err)
			if !c.doReconnectBackoff() { // 如果当前重试次数大于最大重试次数，
				// doReconnectBackoff将返回false，因此移除相应的连接
				return false // 新的请求将触发重连
			}
			continue
		}
		c.setRawConn(conn)
		c.done = make(chan struct{})
		if !c.isIdle {
			c.isIdle = true
			c.connsAddIdle()
		}
		// 成功重连，移除关闭标志并重置c.err
		c.err = nil
		c.closed = false
		go c.reading()
		go c.writing()
		return true
	}
}

// doReconnectBackoff 执行重连退避策略
// 返回:
//
//	bool: 是否继续重连
func (c *Connection) doReconnectBackoff() bool {
	cur := time.Now()
	if !c.lastReconnectTime.IsZero() && c.lastReconnectTime.Add(reconnectCountResetInterval).Before(cur) {
		// 达到重置间隔时清除重连计数
		c.reconnectCount = 0
	}
	c.reconnectCount++
	c.lastReconnectTime = cur
	if c.reconnectCount > maxReconnectCount {
		logging.Tracef("reconnection reaches its limit: %d", maxReconnectCount)
		return false
	}
	currentBackoff := time.Duration(c.reconnectCount) * initialBackoff
	if currentBackoff > maxBackoff {
		currentBackoff = maxBackoff
	}
	time.Sleep(currentBackoff)
	return true
}

// remove 移除指定的虚拟连接
// 参数:
//
//	virConnID: 要移除的虚拟连接ID
func (c *Connection) remove(virConnID uint32) {
	if needDestroy := c.doRemove(virConnID); needDestroy {
		c.destroy()
	}
}

// doRemove 执行实际的虚拟连接移除操作
// 参数:
//
//	virConnID: 要移除的虚拟连接ID
//
// 返回:
//
//	bool: 是否需要销毁连接
func (c *Connection) doRemove(virConnID uint32) (needDestroy bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.virConns, virConnID)
	if c.enableIdleRemove {
		return c.idleRemove()
	}
	return false
}

// idleRemove 空闲连接移除逻辑
// 返回:
//
//	bool: 是否需要销毁连接
func (c *Connection) idleRemove() (needDestroy bool) {
	// 判断当前连接是否空闲
	if len(c.virConns) != 0 {
		return false
	}
	// 检查连接是否已关闭
	if c.closed {
		return false
	}
	if !c.isIdle {
		c.isIdle = true
		c.connsAddIdle()
	}
	// 判断当前节点空闲连接是否超过最大值
	if !c.connsNeedIdleRemove() {
		return false
	}
	// 关闭当前连接
	c.closed = true
	close(c.done)
	if conn := c.getRawConn(); conn != nil {
		conn.Close()
	}
	// 从连接集中移除当前连接
	return true
}

// canGetVirConn 检查是否可以获取新的虚拟连接
// 返回:
//
//	bool: 是否可以获取新连接
func (c *Connection) canGetVirConn() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.maxVirConns == 0 || // 0表示无限制
		len(c.virConns) < c.maxVirConns
}

// startConnect 开始实际执行连接逻辑
// 参数:
//
//	opts: 连接配置选项
//	dialTimeout: 拨号超时时间
func (c *Connection) startConnect(opts *GetOptions, dialTimeout time.Duration) {
	c.fp = opts.FP
	if err := c.dial(dialTimeout, opts); err != nil {
		// 第一次连接建立失败直接失败，
		// 让上层触发下一次重新建立连接
		c.close(err, false)
		return
	}

	go c.reading()
	go c.writing()
}

// dial 执行拨号操作
// 参数:
//
//	timeout: 拨号超时时间
//	opts: 连接配置选项
//
// 返回:
//
//	error: 拨号过程中发生的错误
func (c *Connection) dial(timeout time.Duration, opts *GetOptions) error {
	if c.isStream {
		conn, dialOpts, err := dialTCP(timeout, opts)
		c.dialOpts = dialOpts
		if err != nil {
			return err
		}
		c.setRawConn(conn)
	} else {
		conn, addr, err := dialUDP(opts)
		if err != nil {
			return err
		}
		c.addr = addr
		c.packetConn = conn
		c.packetBuffer = packetbuffer.New(conn, maxBufferSize)
	}
	return nil
}

// reading 读取数据的协程
// 持续从连接读取数据并分发给对应的虚拟连接
func (c *Connection) reading() {
	var lastErr error
	for {
		select {
		case <-c.done:
			return
		default:
		}

		vid, buf, err := c.parse()
		if err != nil {
			// 如果TCP解包出现错误，可能会导致所有后续解析出现问题，
			// 因此需要关闭并重连
			if c.isStream {
				lastErr = err
				logging.Tracef("reconnect on read err: %+v", err)
				break
			}
			// UDP按单个数据包处理，接收非法数据包不影响后续数据包处理逻辑，可以继续接收数据包
			logging.Tracef("decode packet err: %s", err)
			continue
		}

		c.mu.RLock()
		vc, ok := c.virConns[vid]
		c.mu.RUnlock()
		if !ok {
			continue
		}
		vc.recvQueue.Put(buf)
	}
	c.close(lastErr, true)
}

// writing 写入数据的协程
// 持续从写入缓冲区读取数据并写入到连接
func (c *Connection) writing() {
	var lastErr error
L:
	for {
		select {
		case <-c.done:
			return
		case it := <-c.writeBuffer:
			if err := c.writeAll(it); err != nil {
				if c.isStream { // 如果TCP写入数据失败，会导致对端关闭连接
					lastErr = err
					logging.Tracef("reconnect on write err: %+v", err)
					break L
				}
				// UDP发送数据包失败，可以继续发送数据包
				logging.Tracef("multiplexed send UDP packet failed: %v", err)
				continue
			}
		}
	}
	c.close(lastErr, true)
}

// parse 解析接收到的数据
// 根据连接类型调用相应的解析方法
// 返回:
//
//	vid: 虚拟连接ID
//	buf: 解析出的数据
//	err: 解析过程中发生的错误
func (c *Connection) parse() (vid uint32, buf []byte, err error) {
	if c.isStream {
		return c.fp.Parse(c.getRawConn())
	}

	defer func() {
		closeErr := c.packetBuffer.Next()
		if closeErr == nil {
			return
		}
		if err == nil {
			err = closeErr
			return
		}
		err = fmt.Errorf("parse error %w, close packet error %s", err, closeErr)
	}()

	return c.fp.Parse(c.packetBuffer)
}
