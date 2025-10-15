package server_transport

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/ip"
	"github.com/fengzhongzhu1621/xgo/network/transport/options"
	"github.com/fengzhongzhu1621/xgo/xerror"
)

var _ IServerTransport = (*serverStreamTransport)(nil)

// serverStreamTransport 实现ServerStreamTransport接口并保持与原始serverTransport的向后兼容性
// 支持流式服务器传输，可以发送流式消息
// 流式服务器传输的实现方式是，在serverTransport的基础上，添加一个Send方法，用于发送流式消息
type serverStreamTransport struct {
	// 保持与原始serverTransport的向后兼容性
	serverTransport
}

// NewServerStreamTransport 创建新的ServerTransport，包装在serverStreamTransport中作为返回的ServerStreamTransport接口
func NewServerStreamTransport(opt ...options.ServerTransportOption) IServerStreamTransport {
	s := newServerTransport(opt...)
	return &serverStreamTransport{s}
}

// DefaultServerStreamTransport 默认的ServerStreamTransport实例
var DefaultServerStreamTransport = NewServerStreamTransport()

// ListenAndServe 实现ServerTransport接口
// 为了兼容普通RPC和流式RPC，我们使用serverTransport.ListenAndServe函数
func (st *serverStreamTransport) ListenAndServe(ctx context.Context, opts ...options.ListenServeOption) error {
	return st.serverTransport.ListenAndServe(ctx, opts...)
}

// Send 发送流消息的方法
func (st *serverStreamTransport) Send(ctx context.Context, req []byte) error {
	msg := codec.Message(ctx)
	raddr := msg.RemoteAddr() // 远程地址
	laddr := msg.LocalAddr()  // 本地地址
	if raddr == nil || laddr == nil {
		return xerror.NewFrameError(xerror.RetServerSystemErr,
			fmt.Sprintf("Address is invalid, local: %s, remote: %s", laddr, raddr))
	}
	key := ip.AddrToKey(laddr, raddr) // 生成连接键
	st.serverTransport.m.RLock()
	tc, ok := st.serverTransport.addrToConn[key] // 从连接映射中查找连接
	st.serverTransport.m.RUnlock()
	if ok && tc != nil {
		if _, err := tc.rwc.Write(req); err != nil { // 向连接写入数据
			tc.close()    // 关闭连接
			st.Close(ctx) // 关闭流传输
			return err
		}
		return nil
	}
	return xerror.NewFrameError(xerror.RetServerSystemErr, "Can't find conn by addr") // 找不到连接错误
}

// Close 关闭ServerStreamTransport，同时清理缓存的连接
func (st *serverStreamTransport) Close(ctx context.Context) {
	msg := codec.Message(ctx)
	key := ip.AddrToKey(msg.LocalAddr(), msg.RemoteAddr()) // 生成连接键
	st.m.Lock()
	delete(st.serverTransport.addrToConn, key) // 从连接映射中删除连接
	st.m.Unlock()
}
