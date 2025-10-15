package dial

import (
	"net"
)

// TryConnect 尝试建立连接并设置TCP keep-alive选项
// 参数:
//
//	opts: 拨号选项配置
//
// 返回:
//
//	net.Conn: 成功建立的网络连接
//	error: 连接过程中发生的错误
func TryConnect(opts *DialOptions) (net.Conn, error) {
	conn, err := Dial(opts)
	if err != nil {
		return nil, err
	}
	if c, ok := conn.(*net.TCPConn); ok {
		c.SetKeepAlive(true)
	}
	return conn, nil
}
