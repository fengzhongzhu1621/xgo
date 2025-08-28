package conn

import "context"

// 简单连接池，只包含一个连接.
type SingleConnPool struct {
	pool      Pooler
	cn        *Conn
	stickyErr error
}

// 判断SingleConnPool是否实现了Pooler中的所有接口.
var _ Pooler = (*SingleConnPool)(nil)

// 创建一个简单连接池.
func NewSingleConnPool(pool Pooler, cn *Conn) *SingleConnPool {
	return &SingleConnPool{
		pool: pool,
		cn:   cn,
	}
}

// 创建一个连接并放到连接池.
func (p *SingleConnPool) NewConn(ctx context.Context) (*Conn, error) {
	return p.pool.NewConn(ctx)
}

// 关闭一个连接.
func (p *SingleConnPool) CloseConn(cn *Conn) error {
	return p.pool.CloseConn(cn)
}

func (p *SingleConnPool) Get(ctx context.Context) (*Conn, error) {
	if p.stickyErr != nil {
		return nil, p.stickyErr
	}
	return p.cn, nil
}

func (p *SingleConnPool) Put(ctx context.Context, cn *Conn) {}

func (p *SingleConnPool) Remove(ctx context.Context, cn *Conn, reason error) {
	p.cn = nil
	p.stickyErr = reason
}

func (p *SingleConnPool) Close() error {
	p.cn = nil
	p.stickyErr = ErrClosed
	return nil
}

func (p *SingleConnPool) Len() int {
	return 0
}

func (p *SingleConnPool) IdleLen() int {
	return 0
}

func (p *SingleConnPool) Stats() *Stats {
	return &Stats{}
}
