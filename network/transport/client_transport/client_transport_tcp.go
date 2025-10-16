package client_transport

import (
	"context"
	"net"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/buffer"
	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/connpool"
	"github.com/fengzhongzhu1621/xgo/network/dial"
	"github.com/fengzhongzhu1621/xgo/network/multiplexed"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/xerror"
)

// tcpRoundTrip sends tcp request. It supports send, sendAndRcv, keepalive and multiplex.
func (c *clientTransport) tcpRoundTrip(ctx context.Context, reqData []byte,
	opts *options.RoundTripOptions) ([]byte, error) {
	if opts.Pool == nil {
		return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
			"tcp client transport: connection pool empty")
	}

	if opts.FramerBuilder == nil {
		return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
			"tcp client transport: framer builder empty")
	}

	// 从连接池获取一个空闲连接
	conn, err := c.dialTCP(ctx, opts)
	if err != nil {
		return nil, err
	}
	// TCP connection is exclusively multiplexed. Close determines whether connection should be put
	// back into the connection pool to be reused.
	defer conn.Close()

	// 上下文追加信息
	msg := codec.Message(ctx)
	msg.WithRemoteAddr(conn.RemoteAddr())
	msg.WithLocalAddr(conn.LocalAddr())

	if ctx.Err() == context.Canceled {
		return nil, xerror.NewFrameError(xerror.RetClientCanceled,
			"tcp client transport canceled before Write: "+ctx.Err().Error())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, xerror.NewFrameError(xerror.RetClientTimeout,
			"tcp client transport timeout before Write: "+ctx.Err().Error())
	}

	// 发送请求数据
	err = c.tcpWriteFrame(ctx, conn, reqData)
	if err != nil {
		return nil, err
	}

	// 读取响应数据
	rspData, err := c.tcpReadFrame(conn, opts)

	return rspData, err
}

// dialTCP establishes a TCP connection.
func (c *clientTransport) dialTCP(ctx context.Context, opts *options.RoundTripOptions) (net.Conn, error) {
	// If ctx has canceled or timeout, just return.
	if ctx.Err() == context.Canceled {
		return nil, xerror.NewFrameError(xerror.RetClientCanceled,
			"client canceled before tcp dial: "+ctx.Err().Error())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, xerror.NewFrameError(xerror.RetClientTimeout,
			"client timeout before tcp dial: "+ctx.Err().Error())
	}
	var timeout time.Duration
	d, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(d)
	}

	var conn net.Conn
	var err error
	// Short connection mode, directly dial a connection.
	if opts.DisableConnectionPool {
		// The connection is established using the minimum of ctx timeout and connecting timeout.
		if opts.DialTimeout > 0 && opts.DialTimeout < timeout {
			timeout = opts.DialTimeout
		}
		conn, err = dial.Dial(&dial.DialOptions{
			Network:       opts.Network,
			Address:       opts.Address,
			LocalAddr:     opts.LocalAddr,
			Timeout:       timeout,
			CACertFile:    opts.CACertFile,
			TLSCertFile:   opts.TLSCertFile,
			TLSKeyFile:    opts.TLSKeyFile,
			TLSServerName: opts.TLSServerName,
		})
		if err != nil {
			return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
				"tcp client transport dial: "+err.Error())
		}
		if ok {
			conn.SetDeadline(d)
		}
		return conn, nil
	}

	// Connection pool mode, get connection from pool.
	getOpts := connpool.NewGetOptions()
	getOpts.WithContext(ctx)
	getOpts.WithFramerBuilder(opts.FramerBuilder)
	getOpts.WithDialTLS(opts.TLSCertFile, opts.TLSKeyFile, opts.CACertFile, opts.TLSServerName)
	getOpts.WithLocalAddr(opts.LocalAddr)
	getOpts.WithDialTimeout(opts.DialTimeout)
	getOpts.WithProtocol(opts.Protocol)
	conn, err = opts.Pool.Get(opts.Network, opts.Address, getOpts)
	if err != nil {
		return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
			"tcp client transport connection pool: "+err.Error())
	}
	if ok {
		conn.SetDeadline(d)
	}
	return conn, nil
}

// tcpWriteFrame 写入TCP帧数据到连接
// 参数:
//   - ctx: 上下文，用于取消和超时控制
//   - conn: 网络连接对象
//   - reqData: 要发送的请求数据字节切片
//
// 返回值: 错误信息，成功返回nil
func (c *clientTransport) tcpWriteFrame(ctx context.Context, conn net.Conn, reqData []byte) error {
	sentNum := 0  // 已发送字节数
	num := 0      // 单次写入的字节数
	var err error // 错误变量

	// 循环发送直到所有数据发送完成
	for sentNum < len(reqData) {
		// 从已发送位置开始写入剩余数据
		num, err = conn.Write(reqData[sentNum:])
		if err != nil {
			// 检查是否为超时错误
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// 返回客户端超时错误
				return xerror.NewFrameError(xerror.RetClientTimeout,
					"tcp client transport Write: "+err.Error())
			}
			// 返回网络错误
			return xerror.NewFrameError(xerror.RetClientNetErr,
				"tcp client transport Write: "+err.Error())
		}

		// 更新已发送字节数
		sentNum += num
	}

	// 所有数据发送成功，返回nil
	return nil
}

// tcpReadFrame 从TCP连接读取帧数据
// 参数:
// - conn: 网络连接对象
// - opts: 往返选项配置，包含连接池、请求类型等设置
// 返回值:
// - []byte: 读取到的帧数据
// - error: 读取过程中发生的错误
func (c *clientTransport) tcpReadFrame(conn net.Conn, opts *options.RoundTripOptions) ([]byte, error) {
	// 检查请求类型是否为"仅发送"模式
	// 如果是仅发送模式，不需要等待响应，直接返回无响应错误
	if opts.ReqType == codec.SendOnly {
		return nil, xerror.ErrClientNoResponse
	}

	var fr codec.IFramer // 帧读取器接口变量

	// 根据连接池禁用标志选择不同的帧读取器创建策略
	if opts.DisableConnectionPool {
		// 短连接模式：为每个连接创建新的帧读取器
		// 使用缓冲读取器包装连接，提高读取性能
		// 这种模式下帧读取器与连接生命周期绑定
		fr = opts.FramerBuilder.New(buffer.NewReader(conn))
	} else {
		// 连接池模式：尝试从连接对象本身获取帧读取器
		// 在连接池中，帧读取器通常已经与连接绑定，可以复用
		var ok bool
		fr, ok = conn.(codec.IFramer) // 类型断言，检查连接是否实现了IFramer接口

		// 如果连接没有实现帧读取器接口，返回连接失败错误
		if !ok {
			return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
				"tcp client transport: framer not implemented")
		}
	}

	// 使用帧读取器读取完整的帧数据
	rspData, err := fr.ReadFrame()
	if err != nil {
		// 错误处理：区分超时错误和其他网络错误
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// 超时错误：客户端读取超时
			return nil, xerror.NewFrameError(xerror.RetClientTimeout,
				"tcp client transport ReadFrame: "+err.Error())
		}
		// 其他读取错误：帧读取失败
		return nil, xerror.NewFrameError(xerror.RetClientReadFrameErr,
			"tcp client transport ReadFrame: "+err.Error())
	}

	// 成功读取帧数据，返回给调用方
	return rspData, nil
}

// multiplexed handle multiplexed request.
func (c *clientTransport) multiplexed(ctx context.Context, req []byte, opts *options.RoundTripOptions) ([]byte, error) {
	if opts.FramerBuilder == nil {
		return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
			"tcp client transport: framer builder empty")
	}
	getOpts := multiplexed.NewGetOptions()
	getOpts.WithVID(opts.Msg.RequestID())
	fp, ok := opts.FramerBuilder.(multiplexed.IFrameParser)
	if !ok {
		return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
			"frame builder does not implement multiplexed.FrameParser")
	}
	getOpts.WithFrameParser(fp)
	getOpts.WithDialTLS(opts.TLSCertFile, opts.TLSKeyFile, opts.CACertFile, opts.TLSServerName)
	getOpts.WithLocalAddr(opts.LocalAddr)
	conn, err := opts.Multiplexed.GetMuxConn(ctx, opts.Network, opts.Address, getOpts)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 上下文追加信息
	msg := codec.Message(ctx)
	msg.WithRemoteAddr(conn.RemoteAddr())

	if err := conn.Write(req); err != nil {
		return nil, xerror.NewFrameError(xerror.RetClientNetErr,
			"tcp client multiplexed transport Write: "+err.Error())
	}

	// SendOnly does not need to read response.
	if opts.ReqType == codec.SendOnly {
		return nil, xerror.ErrClientNoResponse
	}

	buf, err := conn.Read()
	if err != nil {
		if err == context.Canceled {
			return nil, xerror.NewFrameError(xerror.RetClientCanceled,
				"tcp client multiplexed transport ReadFrame: "+err.Error())
		}
		if err == context.DeadlineExceeded {
			return nil, xerror.NewFrameError(xerror.RetClientTimeout,
				"tcp client multiplexed transport ReadFrame: "+err.Error())
		}
		return nil, xerror.NewFrameError(xerror.RetClientNetErr,
			"tcp client multiplexed transport ReadFrame: "+err.Error())
	}
	return buf, nil
}
