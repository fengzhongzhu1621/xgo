package multiplexed

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	// 是 HashiCorp 提供的一个 Go 语言库，它能将多个错误合并成一个标准 error，非常适合需要收集和统一管理多个错误的场景
	"github.com/hashicorp/go-multierror"
)

// Connections 表示特定目标地址的连接集合
type Connections struct {
	nodeKey    string       // 节点键，格式为"network_address"
	maxIdle    int          // 最大空闲连接数
	opts       *PoolOptions // 连接池配置选项
	destructor func()       // 销毁回调函数，当连接集合被驱逐时调用

	// mu 保护以下字段的并发安全
	mu              sync.Mutex
	conns           []*Connection // 连接列表
	currentIdle     int32         // 当前空闲连接数（原子操作）
	roundRobinIndex int           // 轮询索引，用于负载均衡
	expelled        bool          // 是否已被驱逐
	err             error         // 连接集合的错误状态
}

// initialize 初始化连接集合
// 根据配置的连接数量创建初始连接
// 参数:
//
//	opts: 连接获取配置选项
func (cs *Connections) initialize(opts *GetOptions) {
	for i := 0; i < cs.opts.connectNumberPerHost; i++ {
		cs.newConn(opts)
	}
}

// addIdle 增加空闲连接计数
// 如果启用了最大空闲连接限制，则原子增加计数
func (cs *Connections) addIdle() {
	if cs.maxIdle > 0 {
		atomic.AddInt32(&cs.currentIdle, 1)
	}
}

// subIdle 减少空闲连接计数
// 如果启用了最大空闲连接限制，则原子减少计数
func (cs *Connections) subIdle() {
	if cs.maxIdle > 0 {
		atomic.AddInt32(&cs.currentIdle, -1)
	}
}

// expel 驱逐指定的连接
// 参数:
//
//	c: 要驱逐的连接
func (cs *Connections) expel(c *Connection) {
	cs.mu.Lock()
	cs.subIdle()
	cs.conns = filterOutConnection(cs.conns, c)
	cs.err = multierror.Append(cs.err, c.err).ErrorOrNil()
	if cs.expelled || len(cs.conns) > 0 {
		cs.mu.Unlock()
		return
	}
	cs.expelled = true
	cs.mu.Unlock()
	cs.destructor()
}

// newConn 创建新的物理连接
// 参数:
//
//	opts: 连接获取配置选项
//
// 返回:
//
//	*Connection: 新创建的连接
func (cs *Connections) newConn(opts *GetOptions) *Connection {
	c := &Connection{
		network:          opts.network,
		address:          opts.address,
		virConns:         make(map[uint32]*VirtualConnection),
		done:             make(chan struct{}),
		dropFull:         cs.opts.dropFull,
		maxVirConns:      cs.opts.maxVirConnsPerConn,
		writeBuffer:      make(chan []byte, cs.opts.sendQueueSize),
		isStream:         opts.isStream,
		isIdle:           true,
		enableIdleRemove: cs.maxIdle > 0 && cs.opts.maxVirConnsPerConn > 0,
		connsAddIdle:     func() { cs.addIdle() },
		connsSubIdle:     func() { cs.subIdle() },
		connsNeedIdleRemove: func() bool {
			return int(atomic.LoadInt32(&cs.currentIdle)) > cs.maxIdle
		},
	}
	c.destroy = func() { cs.expel(c) }
	cs.conns = append(cs.conns, c)

	// 增加空闲连接计数
	cs.addIdle()

	go c.startConnect(opts, cs.opts.dialTimeout)

	return c
}

// pickSingleConcrete 从连接集合中选取一个具体的物理连接
// 参数:
//
//	ctx: 上下文
//	opts: 连接获取配置选项
//
// 返回:
//
//	*Connection: 选中的物理连接
//	error: 选取过程中发生的错误
func (cs *Connections) pickSingleConcrete(ctx context.Context, opts *GetOptions) (*Connection, error) {
	// 始终需要加锁，因为cs.conns的长度可能在另一个goroutine中改变
	// 示例情况:
	//  1. 在空闲连接移除期间，cs.conns的长度会减少
	//  2. 如果达到最大重试次数，cs.conns的长度会减少
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.expelled {
		return nil, fmt.Errorf("node key: %s, err: %w, caused by sub errors on conns: %+v",
			cs.nodeKey, ErrConnectionsHaveBeenExpelled, cs.err)
	}
	if cs.opts.maxVirConnsPerConn == 0 {
		// 每个具体连接的虚拟连接数无限制，执行轮询
		cs.roundRobinIndex = (cs.roundRobinIndex + 1) % cs.opts.connectNumberPerHost
		if cs.roundRobinIndex >= len(cs.conns) {
			// 当前具体连接数已减少到预期数量以下
			// 创建一个新的具体连接来补充
			cs.roundRobinIndex = len(cs.conns)
			return cs.newConn(opts), nil
		}
		return cs.conns[cs.roundRobinIndex], nil
	}
	for _, c := range cs.conns {
		if c.canGetVirConn() {
			return c, nil
		}
	}
	return cs.newConn(opts), nil
}
