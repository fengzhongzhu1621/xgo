// package packetbuffer for manipulating byte slices
// 这是一个数据包缓冲区实现，用于高效处理网络数据包。主要功能：

// 核心特性：
//
// * 封装了 net.PacketConn，实现 io.Reader 接口
// * 使用内存分配器管理缓冲区，避免频繁内存分配
// * 支持连续读取数据包内容
package packetbuffer

import (
	"fmt"
	"io"
	"net"

	"github.com/fengzhongzhu1621/xgo/pool/allocator"
)

var _ io.Reader = (*PacketBuffer)(nil)

// New creates a packet buffer with specific packet connection and size.
// 创建指定大小和数据包连接的数据包缓冲区
func New(conn net.PacketConn, size int) *PacketBuffer {
	buf, i := allocator.Malloc(size) // 从内存分配器获取缓冲区
	return &PacketBuffer{
		buf:      buf,  // 数据存储缓冲区
		conn:     conn, // 底层数据包连接
		toBeFree: i,    // 需要释放的内存标识
	}
}

// PacketBuffer encapsulates a packet connection and implements the io.Reader interface.
// PacketBuffer 封装数据包连接并实现 io.Reader 接口
type PacketBuffer struct {
	buf      []byte         // 数据存储缓冲区
	toBeFree interface{}    // 需要释放的内存标识
	conn     net.PacketConn // 底层数据包连接

	raddr net.Addr // 当前数据包的远程地址
	r, w  int      // 读指针和写指针位置
}

// Read reads data from the packet. Continuous reads cannot cross between multiple packet only if Close is called.
// 从数据包读取数据，连续读取不能跨数据包（除非调用Close）
func (pb *PacketBuffer) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil // 目标缓冲区为空，直接返回
	}
	if pb.w == 0 {
		// 缓冲区为空，连接读取新数据包到缓冲区
		n, raddr, err := pb.conn.ReadFrom(pb.buf)
		if err != nil {
			return 0, err // 读取失败
		}
		pb.w = n         // 设置写指针到数据末尾
		pb.raddr = raddr // 保存远程地址
	}
	n = copy(p, pb.buf[pb.r:pb.w]) // 复制数据到目标缓冲区
	if n == 0 {
		return 0, io.EOF // 没有数据可读
	}
	pb.r += n // 移动读指针
	return n, nil
}

// Next is used to distinguish continuous logic reads. It indicates that the reading on current packet has finished.
// If there remains data unconsumed, Next returns an error and discards the remaining data.
// Next 用于区分连续逻辑读取，标记当前数据包读取完成
// 如果还有未消费的数据，Next会返回错误并丢弃剩余数据
func (pb *PacketBuffer) Next() error {
	if pb.w == 0 {
		return nil // 没有数据可处理
	}
	var err error
	if remain := pb.w - pb.r; remain != 0 {
		err = fmt.Errorf("packet data is not drained, the remaining %d will be dropped", remain) // 数据未完全消费
	}
	pb.r, pb.w = 0, 0 // 重置读写指针
	pb.raddr = nil    // 清空远程地址
	return err
}

// CurrentPacketAddr returns current packet's remote address.
// 返回当前数据包的远程地址
func (pb *PacketBuffer) CurrentPacketAddr() net.Addr {
	return pb.raddr // 返回保存的远程地址
}

// Close closes this buffer and releases resource.
// 关闭缓冲区并释放资源
func (pb *PacketBuffer) Close() {
	allocator.Free(pb.toBeFree) // 通过分配器释放内存
}
