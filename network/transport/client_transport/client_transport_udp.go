package client_transport

import (
	"context"
	"fmt"
	"net"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/pool/packetbuffer"
	"github.com/fengzhongzhu1621/xgo/xerror"
)

const defaultUDPRecvBufSize = 64 * 1024

// udpRoundTrip sends UDP requests.
func (c *clientTransport) udpRoundTrip(ctx context.Context, reqData []byte,
	opts *options.RoundTripOptions) ([]byte, error) {
	if opts.FramerBuilder == nil {
		return nil, xerror.NewFrameError(xerror.RetClientConnectFail,
			"udp client transport: framer builder empty")
	}

	conn, addr, err := c.dialUDP(ctx, opts)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	msg := codec.Message(ctx)
	msg.WithRemoteAddr(addr)
	msg.WithLocalAddr(conn.LocalAddr())

	if ctx.Err() == context.Canceled {
		return nil, xerror.NewFrameError(xerror.RetClientCanceled,
			"udp client transport canceled before Write: "+ctx.Err().Error())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, xerror.NewFrameError(xerror.RetClientTimeout,
			"udp client transport timeout before Write: "+ctx.Err().Error())
	}

	if err := c.udpWriteFrame(conn, reqData, addr, opts); err != nil {
		return nil, err
	}
	return c.udpReadFrame(ctx, conn, opts)
}

// udpReadFrame reads UDP frame.
func (c *clientTransport) udpReadFrame(
	ctx context.Context, conn net.PacketConn, opts *options.RoundTripOptions) ([]byte, error) {
	// If it is SendOnly, returns directly without waiting for the server's response.
	if opts.ReqType == codec.SendOnly {
		return nil, xerror.ErrClientNoResponse
	}

	select {
	case <-ctx.Done():
		return nil, xerror.NewFrameError(xerror.RetClientTimeout, "udp client transport select after Write: "+ctx.Err().Error())
	default:
	}

	buf := packetbuffer.New(conn, defaultUDPRecvBufSize)
	defer buf.Close()
	fr := opts.FramerBuilder.New(buf)
	req, err := fr.ReadFrame()
	if err != nil {
		if e, ok := err.(net.Error); ok {
			if e.Timeout() {
				return nil, xerror.NewFrameError(xerror.RetClientTimeout,
					"udp client transport ReadFrame: "+err.Error())
			}
			return nil, xerror.NewFrameError(xerror.RetClientNetErr,
				"udp client transport ReadFrom: "+err.Error())
		}
		return nil, xerror.NewFrameError(xerror.RetClientReadFrameErr,
			"udp client transport ReadFrame: "+err.Error())
	}
	// One packet of udp corresponds to one trpc packet,
	// and after parsing, there should not be any remaining data
	if err := buf.Next(); err != nil {
		return nil, xerror.NewFrameError(xerror.RetClientReadFrameErr,
			fmt.Sprintf("udp client transport ReadFrame: %s", err))
	}

	// Framer is used for every request so there is no need to copy memory.
	return req, nil
}

// udpWriteReqData write UDP frame.
func (c *clientTransport) udpWriteFrame(conn net.PacketConn,
	reqData []byte, addr *net.UDPAddr, opts *options.RoundTripOptions) error {
	// Sending udp request packets
	var num int
	var err error
	if opts.ConnectionMode == options.Connected {
		udpconn := conn.(*net.UDPConn)
		num, err = udpconn.Write(reqData)
	} else {
		num, err = conn.WriteTo(reqData, addr)
	}
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return xerror.NewFrameError(xerror.RetClientTimeout, "udp client transport WriteTo: "+err.Error())
		}
		return xerror.NewFrameError(xerror.RetClientNetErr, "udp client transport WriteTo: "+err.Error())
	}
	if num != len(reqData) {
		return xerror.NewFrameError(xerror.RetClientNetErr, "udp client transport WriteTo: num mismatch")
	}
	return nil
}

// dialUDP establishes an UDP connection.
func (c *clientTransport) dialUDP(ctx context.Context, opts *options.RoundTripOptions) (net.PacketConn, *net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr(opts.Network, opts.Address)
	if err != nil {
		return nil, nil, xerror.NewFrameError(xerror.RetClientNetErr,
			"udp client transport ResolveUDPAddr: "+err.Error())
	}

	var conn net.PacketConn
	if opts.ConnectionMode == options.Connected {
		var localAddr net.Addr
		if opts.LocalAddr != "" {
			localAddr, err = net.ResolveUDPAddr(opts.Network, opts.LocalAddr)
			if err != nil {
				return nil, nil, xerror.NewFrameError(xerror.RetClientNetErr,
					"udp client transport LocalAddr ResolveUDPAddr: "+err.Error())
			}
		}
		dialer := net.Dialer{
			LocalAddr: localAddr,
		}
		var udpConn net.Conn
		udpConn, err = dialer.Dial(opts.Network, opts.Address)
		if err != nil {
			return nil, nil, xerror.NewFrameError(xerror.RetClientConnectFail,
				fmt.Sprintf("dial udp fail: %s", err.Error()))
		}

		var ok bool
		conn, ok = udpConn.(net.PacketConn)
		if !ok {
			return nil, nil, xerror.NewFrameError(xerror.RetClientConnectFail,
				"udp conn not implement net.PacketConn")
		}
	} else {
		// Listen on all available IP addresses of the local system by default,
		// and a port number is automatically chosen.
		const defaultLocalAddr = ":"
		localAddr := defaultLocalAddr
		if opts.LocalAddr != "" {
			localAddr = opts.LocalAddr
		}
		conn, err = net.ListenPacket(opts.Network, localAddr)
	}
	if err != nil {
		return nil, nil, xerror.NewFrameError(xerror.RetClientNetErr, "udp client transport Dial: "+err.Error())
	}
	d, ok := ctx.Deadline()
	if ok {
		conn.SetDeadline(d)
	}
	return conn, addr, nil
}
