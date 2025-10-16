package transport

import (
	"context"
	"encoding/binary"
	"errors"
	"io"

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
	// 模拟服务端报错，将连接关闭消息交给业务处理层处理，并会关闭 tcp 连接
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

// ReadFrame 从数据流中读取一个完整的帧
// 帧格式: [4字节长度头] + [实际数据内容]
// 返回值:
// - []byte: 包含长度头和数据的完整帧
// - error: 读取过程中发生的错误
func (f *framer) ReadFrame() ([]byte, error) {
	// 检查帧构建器是否设置了预定义错误
	// 用于测试场景或错误注入
	if f.fb.errSet {
		return nil, f.fb.err
	}

	// 	+----------------+---------------------+
	// | 4字节长度头    | 实际数据内容        |
	// | (大端序uint32) | (长度由头部指定)    |
	// +----------------+---------------------+

	var lenData [4]byte // 4字节缓冲区，用于存储帧长度信息

	// 读取4个字节的帧长度数据
	// io.ReadFull确保读取完整4个字节，否则返回错误
	_, err := io.ReadFull(f.r, lenData[:])
	if err != nil {
		return nil, err
	}
	// 将4字节的大端序数据转换为uint32类型的长度值
	// 大端序(网络字节序)确保跨平台兼容性
	length := binary.BigEndian.Uint32(lenData[:])

	// 分配完整帧的缓冲区: 长度头(4字节) + 实际数据(length字节)
	msg := make([]byte, len(lenData)+int(length))

	// 将长度头数据复制到消息缓冲区的前4个字节
	copy(msg, lenData[:])

	// 从数据流中读取实际的数据内容到消息缓冲区的剩余部分
	// 从第5个字节开始(msg[len(lenData):])存放实际数据
	_, err = io.ReadFull(f.r, msg[len(lenData):])
	if err != nil {
		return nil, err
	}

	// 返回完整的帧数据(包含长度头+实际数据)
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
