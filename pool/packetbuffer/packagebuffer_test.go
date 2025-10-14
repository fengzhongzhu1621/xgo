package packetbuffer

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

type udpServer struct {
	cancel context.CancelFunc
	conn   net.PacketConn
}

// 启动一个echo udp server
func (s *udpServer) start(ctx context.Context) error {
	var err error

	// 创建了一个 UDP 服务器
	s.conn, err = net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	ctx, s.cancel = context.WithCancel(ctx)
	go func() {
		buf := make([]byte, 65535)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			// 读取客户端信息，实现了一个 echo server
			n, addr, err := s.conn.ReadFrom(buf)
			if err != nil {
				log.Println("l.ReadFrom err: ", err)
				return
			}
			s.conn.WriteTo(buf[:n], addr)
		}
	}()
	return nil
}

func (s *udpServer) stop() {
	s.cancel()
	s.conn.Close()
}

func TestPacketReaderSucceed(t *testing.T) {
	// 启动一个 UDP 服务器
	s := &udpServer{}
	s.start(context.Background())
	t.Cleanup(s.stop)

	// 创建一个 UDP 客户端
	p, err := net.ListenPacket("udp", "127.0.0.1:0")
	require.Nil(t, err)
	// 向 UDP 服务器发送数据
	_, err = p.WriteTo([]byte("helloworldA"), s.conn.LocalAddr())
	require.Nil(t, err)

	buf := New(p, 65535)
	defer buf.Close()

	// 从 UDP 服务器读取数据
	result := make([]byte, 20)
	n, err := buf.Read(result)
	require.Nil(t, err)
	require.Equal(t, []byte("helloworldA"), result[:n])
	require.Equal(t, s.conn.LocalAddr(), buf.CurrentPacketAddr())

	// 读取下一个数据包（无数据包可读）
	_, err = buf.Read(result)
	require.Equal(t, io.EOF, err)

	// 重置缓冲区
	require.Nil(t, buf.Next())

	// 再次读取数据包
	_, err = p.WriteTo([]byte("helloworldB"), s.conn.LocalAddr())
	require.Nil(t, err)
	n, err = buf.Read(result)
	require.Nil(t, err)
	require.Equal(t, []byte("helloworldB"), result[:n])
}

func TestPacketReaderFailed(t *testing.T) {
	s := &udpServer{}
	s.start(context.Background())
	t.Cleanup(s.stop)

	p, err := net.ListenPacket("udp", "127.0.0.1:0")
	require.Nil(t, err)
	_, err = p.WriteTo([]byte("helloworld"), s.conn.LocalAddr())

	require.Nil(t, err)
	buf := New(p, 65535)
	defer buf.Close()

	n, err := buf.Read(nil)
	require.Nil(t, err)
	require.Equal(t, 0, n)
	result := make([]byte, 5)

	_, err = buf.Read(result)
	require.Nil(t, err)
	// There are some remaining data in the buf that have not been read.
	require.NotNil(t, buf.Next())
}
