package transport

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/xerror"
)

var savedListenerPort int

var defaultStreamID uint32 = 100

type helloRequest struct {
	Name string
	Msg  string
}

type helloResponse struct {
	Name string
	Msg  string
	Code int
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// echoHandler 模拟服务端业务处理逻辑（返回成功）
type echoHandler struct{}

func (h *echoHandler) Handle(ctx context.Context, req []byte) ([]byte, error) {
	rsp := make([]byte, len(req))
	copy(rsp, req)
	return rsp, nil
}

// echoHandler 模拟服务端流式业务处理逻辑（返回失败）
type echoStreamHandler struct{}

func (h *echoStreamHandler) Handle(ctx context.Context, req []byte) ([]byte, error) {
	rsp := make([]byte, len(req))
	copy(rsp, req)
	return rsp, xerror.ErrServerNoResponse
}

type errorHandler struct{}

func (h *errorHandler) Handle(ctx context.Context, req []byte) ([]byte, error) {
	return nil, errors.New("handle error")
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type framerBuilder struct {
	errSet bool
	err    error
	safe   bool
}

// SetError sets frameBuilder error.
func (fb *framerBuilder) SetError(err error) {
	fb.errSet = true
	fb.err = err
}

func (fb *framerBuilder) ClearError() {
	fb.errSet = false
	fb.err = nil
}

func (fb *framerBuilder) New(r io.Reader) codec.IFramer {
	return &framer{r: r, fb: fb}
}

type framer struct {
	fb *framerBuilder
	r  io.Reader
}

func (f *framer) ReadFrame() ([]byte, error) {
	if f.fb.errSet {
		return nil, f.fb.err
	}
	var lenData [4]byte

	_, err := io.ReadFull(f.r, lenData[:])
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lenData[:])

	msg := make([]byte, len(lenData)+int(length))
	copy(msg, lenData[:])

	_, err = io.ReadFull(f.r, msg[len(lenData):])
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (f *framer) IsSafe() bool {
	return f.fb.safe
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/* --------------------------------------------------- mock multiplexed framer ---------------------------------------------- */

var _ codec.IFramerBuilder = (*multiplexedFramerBuilder)(nil)

type multiplexedFramerBuilder struct {
	errSet bool
	err    error
	safe   bool
}

func (fb *multiplexedFramerBuilder) New(r io.Reader) codec.IFramer {
	return &multiplexedFramer{r: r, fb: fb}
}

func (fb *multiplexedFramerBuilder) SetError(err error) {
	fb.errSet = true
	fb.err = err
}

func (fb *multiplexedFramerBuilder) ClearError() {
	fb.errSet = false
	fb.err = nil
}

func (fb *multiplexedFramerBuilder) Parse(rc io.Reader) (vid uint32, buf []byte, err error) {
	buf, err = fb.New(rc).ReadFrame()
	if err != nil {
		return 0, nil, err
	}
	return binary.BigEndian.Uint32(buf[:4]), buf, nil
}

var _ codec.IFramer = (*multiplexedFramer)(nil)
var _ codec.ISafeFramer = (*multiplexedFramer)(nil)

type multiplexedFramer struct {
	fb *multiplexedFramerBuilder
	r  io.Reader
}

func (f *multiplexedFramer) ReadFrame() ([]byte, error) {
	if f.fb.errSet {
		return nil, f.fb.err
	}
	var headData [8]byte

	_, err := io.ReadFull(f.r, headData[:])
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(headData[4:])

	msg := make([]byte, len(headData)+int(length))
	copy(msg, headData[:])

	_, err = io.ReadFull(f.r, msg[len(headData):])
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (f *multiplexedFramer) IsSafe() bool {
	return f.fb.safe
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getFreeAddr(network string) string {
	p, err := getFreePort(network)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(":%d", p)
}

func getFreePort(network string) (int, error) {
	if network == "tcp" || network == "tcp4" || network == "tcp6" {
		addr, err := net.ResolveTCPAddr(network, "localhost:0")
		if err != nil {
			return -1, err
		}

		l, err := net.ListenTCP(network, addr)
		if err != nil {
			return -1, err
		}
		defer l.Close()

		return l.Addr().(*net.TCPAddr).Port, nil
	}

	if network == "udp" || network == "udp4" || network == "udp6" {
		addr, err := net.ResolveUDPAddr(network, "localhost:0")
		if err != nil {
			return -1, err
		}

		l, err := net.ListenUDP(network, addr)
		if err != nil {
			return -1, err
		}
		defer l.Close()

		return l.LocalAddr().(*net.UDPAddr).Port, nil
	}

	return -1, errors.New("invalid network")
}
