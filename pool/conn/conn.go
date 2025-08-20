package conn

import (
	"bufio"
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/fengzhongzhu1621/xgo/proto"
)

var noDeadline = time.Time{}

// Conn 连接接口.
type Conn struct {
	usedAt  int64    // atomic 原子操作，记录读或写的时间
	netConn net.Conn // 网络连接对象

	rd *proto.Reader
	bw *bufio.Writer
	wr *proto.Writer

	Inited    bool
	pooled    bool      // 标记连接是否放到连接池中
	createdAt time.Time // 连接创建时间
}

// NewConn 创建一个连接.
func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		netConn:   netConn,
		createdAt: time.Now(), // 连接创建时间
	}
	cn.rd = proto.NewReader(netConn)
	cn.bw = bufio.NewWriter(netConn)
	cn.wr = proto.NewWriter(cn.bw)
	cn.SetUsedAt(time.Now())
	return cn
}

// UsedAt 获得时间戳，转换为time.Time格式.
func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

// SetUsedAt 将时间转换为unix时间戳，并保存到变量中.
func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
	cn.rd.Reset(netConn) // 重置缓存，丢弃未被处理的数据
	cn.bw.Reset(netConn) // 重置缓存，丢弃未被处理的数据
}

// Write 通过连接写字节数组.
func (cn *Conn) Write(b []byte) (int, error) {
	return cn.netConn.Write(b)
}

// RemoteAddr 获得服务端地址.
func (cn *Conn) RemoteAddr() net.Addr {
	if cn.netConn != nil {
		return cn.netConn.RemoteAddr()
	}
	return nil
}

// WithReader 设置读超时，返回一个读操作.
func (cn *Conn) WithReader(
	ctx context.Context,
	timeout time.Duration,
	fn func(rd *proto.Reader) error,
) error {
	// 设置读超时时间
	if err := cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout)); err != nil {
		return err
	}
	return fn(cn.rd)
}

// WithWriter 设置写超时，返回一个写操作.
func (cn *Conn) WithWriter(
	ctx context.Context, timeout time.Duration, fn func(wr *proto.Writer) error,
) error {
	if err := cn.netConn.SetWriteDeadline(cn.deadline(ctx, timeout)); err != nil {
		return err
	}
	// 写缓存有数据，则重置连接丢弃缓存中的数据
	if cn.bw.Buffered() > 0 {
		cn.bw.Reset(cn.netConn)
	}

	if err := fn(cn.wr); err != nil {
		return err
	}
	// 写缓存中的数据发送出去
	return cn.bw.Flush()
}

// Close 关闭连接.
func (cn *Conn) Close() error {
	return cn.netConn.Close()
}

// deadline 返回连接的截止日期.
func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	// 记录当前时间
	tm := time.Now()
	cn.SetUsedAt(tm)

	// 计算超时结束时间
	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	// 上下文被取消时，返回上下文的取消时间
	if ctx != nil {
		// 返回 context.Context 被取消的时间，也就是完成工作的截止日期
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}
			// 如果没有超时
			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}

	return noDeadline
}
