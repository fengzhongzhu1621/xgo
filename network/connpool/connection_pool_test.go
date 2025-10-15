package connpool

import (
	"errors"
	"io"
	"net"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/buffer"
	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/dial"
	"github.com/stretchr/testify/require"
)

var (
	ErrFrameSet       = errors.New("framer not set")
	ErrReamFrame      = errors.New("ReadFrame failed")
	ErrRead           = errors.New("Read failed")
	ErrWrite          = errors.New("Write failed")
	ErrSyscallConn    = errors.New("SyscallConn Failed")
	ErrUnexpectedRead = errors.New("unexpected read from socket")

	mockChecker = func(*PoolConn, bool) bool { return true }
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 模拟连接实现 net.Conn 接口，用于读取和写入数据
type noopConn struct {
	closeFunc func()
	suc       bool
}

func (c *noopConn) Read(bs []byte) (int, error) {
	if !c.suc {
		return len(bs), ErrRead
	}
	return len(bs), nil
}
func (c *noopConn) Write(bs []byte) (int, error) {
	if !c.suc {
		return len(bs), ErrWrite
	}
	return len(bs), nil
}

func (c *noopConn) Close() error {
	c.closeFunc()
	return nil
}

func (c *noopConn) LocalAddr() net.Addr              { return nil }
func (c *noopConn) RemoteAddr() net.Addr             { return nil }
func (c *noopConn) SetDeadline(time.Time) error      { return nil }
func (c *noopConn) SetReadDeadline(time.Time) error  { return nil }
func (c *noopConn) SetWriteDeadline(time.Time) error { return nil }

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
var _ codec.IFramerBuilder = (*noopFramerBuilder)(nil)

type noopFramerBuilder struct {
	suc bool
}

func (fb *noopFramerBuilder) New(io.Reader) codec.IFramer {
	return &noopFramer{fb.suc}
}

var _ codec.IFramer = (*noopFramer)(nil)

type noopFramer struct {
	suc bool
}

func (fr *noopFramer) ReadFrame() ([]byte, error) {
	if !fr.suc {
		return make([]byte, 1), ErrReamFrame
	}
	return make([]byte, 1), nil
}

func (fr *noopFramer) IsSafe() bool {
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestInitialMinIdle(t *testing.T) {
	// 创建空的连接池管理器，并设置最小空闲连接数为 10
	// 每次创建连接时，增加 established 计数
	var established int32
	p := NewConnectionPool(
		WithMinIdle(10),
		WithDialFunc(func(*dial.DialOptions) (net.Conn, error) {
			// 每次创建连接时，增加 established 计数
			atomic.AddInt32(&established, 1)
			return &noopConn{closeFunc: func() {}}, nil
		}),
		WithHealthChecker(mockChecker))
	defer closePool(t, p)

	// 创建连接池，新建一个连接放到新创建的连接池中
	// 并设置连接的帧读取器为 noopFramerBuilder
	// 保持最小空闲连接数
	pc, err := p.Get(t.Name(), t.Name(), GetOptions{CustomReader: buffer.NewReader, DialTimeout: time.Second})
	require.Nil(t, err)
	require.Nil(t, pc.Close())

	// 等待 1 秒，检查 established 计数是否为 10
	start := time.Now()
	for time.Since(start) < time.Second {
		if established := atomic.LoadInt32(&established); established == 10 || established == 11 {
			// 判断是否保持最小空闲连接数
			return
		}
		runtime.Gosched()
	}
	require.FailNow(t, "expected 10/11 established connections for fresh pool")
}

// closePool 关闭连接池
func closePool(t *testing.T, p IPool) {
	v, ok := p.(*pool)
	if !ok {
		return
	}

	// 获取连接池的 key，并关闭连接
	key := getNodeKey(t.Name(), t.Name(), "")
	if pool, ok := v.connectionPools.Load(key); ok {
		pool.(*ConnectionPool).Close()
	}
}
