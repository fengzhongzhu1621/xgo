package loadbalance

import (
	"errors"
	"fmt"
	"sync"
)

// 想象你在餐厅点餐：
// * 有3个厨师：A（厉害）、B（一般）、C（新手）
// * 你想让厉害的厨师多做菜：
//   - 每轮先给每个厨师发"代金券"（currentWeight），厉害的发3张，一般发2张，新手发1张
//   - 然后选代金券最多的厨师做菜（A）
//   - 做完菜后，把这个厨师的代金券扣掉总代金券数（A的3张扣掉6张，变成-3）
//   - 下一轮再重新发代金券（A又得到3张，变成0）
//
// 这样长期来看，厉害的厨师做的菜最多（符合权重比例），但每一轮的选择是动态的。

// WeightNode 加权节点
type WeightNode struct {
	addr            string // 节点地址
	Weight          int    // 初始权重
	currentWeight   int    // 当前临时权重
	effectiveWeight int    // 有效权重，默认等于Weight，故障时减1
}

// WeightRoundRobinBalance 加权轮询负载均衡器
type WeightRoundRobinBalance struct {
	mu     sync.Mutex // 互斥锁，保证并发安全
	curIdx int        // 当前索引（虽然加权轮询不需要，但保留以兼容接口）
	rss    []*WeightNode
	rsw    []int // 权重总和缓存（虽然当前实现不需要，但保留以兼容接口）
}

// NewWeightRoundRobinBalance 创建加权轮询负载均衡器实例
func NewWeightRoundRobinBalance() *WeightRoundRobinBalance {
	return &WeightRoundRobinBalance{
		rss: make([]*WeightNode, 0),
	}
}

// Add 添加节点
// 参数：节点地址和权重（字符串形式）
// 返回值：如果参数不足或权重解析失败，返回错误
func (r *WeightRoundRobinBalance) Add(addr string, weight int) error {
	if weight <= 0 {
		return errors.New("weight must be a positive integer")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	node := &WeightNode{
		addr:            addr,        // 节点地址
		Weight:          int(weight), // 节点权重，初始值
		currentWeight:   0,           // 当前临时权重
		effectiveWeight: int(weight), // 有效权重，默认等于Weight，故障时减1
	}
	r.rss = append(r.rss, node)
	return nil
}

// Next 获取下一个节点（加权轮询算法）
// 返回值：选中的节点地址
func (r *WeightRoundRobinBalance) Next() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.rss) == 0 {
		return ""
	}

	var best *WeightNode
	totalWeight := 0

	// 1. 计算总权重并更新每个节点的currentWeight
	for _, node := range r.rss {
		// 更新当前节点的currentWeight
		node.currentWeight += node.effectiveWeight
		fmt.Println(node.currentWeight)
		// 累加总权重
		totalWeight += node.effectiveWeight
		// 动态调整有效权重（故障恢复逻辑）
		// 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if node.effectiveWeight < node.Weight {
			node.effectiveWeight++
		}
		// 选择currentWeight最大的节点，选中最大临时权重节点
		if best == nil || node.currentWeight > best.currentWeight {
			best = node
		}
	}

	// 2. 将选中的节点的currentWeight减去总权重，变更临时权重为 临时权重-有效权重之和
	if best != nil {
		best.currentWeight -= totalWeight
		fmt.Println(best.currentWeight)
	}

	return best.addr
}

// Get 获取节点（包装Next方法）
// 参数：未使用（保留接口一致性）
// 返回值：节点地址和错误（这里总是返回nil错误）
func (r *WeightRoundRobinBalance) Get() (string, error) {
	return r.Next(), nil
}

// Update 更新节点状态（预留接口，用于故障检测等）
func (r *WeightRoundRobinBalance) Update() {
	// 这里可以实现节点健康检查逻辑
	// 例如：检测节点是否可用，更新effectiveWeight
	// 当前实现为空，留给用户自定义
}

// GetNodes 获取所有节点（用于监控或调试）
func (r *WeightRoundRobinBalance) GetNodes() []*WeightNode {
	r.mu.Lock()
	defer r.mu.Unlock()

	nodes := make([]*WeightNode, len(r.rss))
	copy(nodes, r.rss)
	return nodes
}

// SetEffectiveWeight 设置节点的有效权重（用于手动调整）
func (r *WeightRoundRobinBalance) SetEffectiveWeight(addr string, weight int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, node := range r.rss {
		if node.addr == addr {
			if weight < 0 {
				return errors.New("weight cannot be negative")
			}
			node.effectiveWeight = weight
			return nil
		}
	}
	return errors.New("node not found")
}
