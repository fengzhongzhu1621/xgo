package loadbalance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Hash 定义哈希函数类型
type Hash func(data []byte) uint32

// UInt32Slice 用于排序的uint32切片
type UInt32Slice []uint32

func (s UInt32Slice) Len() int           { return len(s) }
func (s UInt32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s UInt32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// ConsistentHashBalance 一致性哈希负载均衡器
type ConsistentHashBalance struct {
	mux      sync.RWMutex      // 读写锁，保护内部数据结构
	hash     Hash              // 哈希函数，默认使用crc32
	replicas int               // 每个节点的虚拟节点数（权重体现）
	keys     UInt32Slice       // 已排序的节点哈希切片
	hashMap  map[uint32]string // 节点哈希到节点地址的映射
}

// NewConsistentHashBalance 创建一致性哈希负载均衡器实例
// replicas: 每个节点的虚拟节点数（权重体现，值越大权重越高）
// fn: 自定义哈希函数，如果为nil则使用crc32.ChecksumIEEE
func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		// 最多32位，保证是一个2^32-1环
		m.hash = crc32.ChecksumIEEE // 默认使用crc32
	}
	return m
}

// IsEmpty 检查是否没有节点
func (c *ConsistentHashBalance) IsEmpty() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return len(c.keys) == 0
}

// Add 添加节点
// 参数：节点地址（如IP）
// 返回值：如果参数不足，返回错误
func (c *ConsistentHashBalance) Add(addr string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// 为每个真实节点生成replicas个虚拟节点
	for i := 0; i < c.replicas; i++ {
		// 虚拟节点名称：节点索引+节点地址（确保唯一性）
		virtualNodeName := strconv.Itoa(i) + addr
		hash := c.hash([]byte(virtualNodeName))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr // 记录哈希值到节点地址的映射
	}

	// 对所有虚拟节点的哈希值进行排序，方便后续二分查找
	sort.Sort(c.keys)
	return nil
}

// Remove 删除节点
// 参数：要删除的节点地址
// 返回值：如果节点不存在，返回错误
func (c *ConsistentHashBalance) Remove(addr string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	// 找到所有属于该节点的虚拟节点哈希值
	var hashesToRemove []uint32
	for hash, nodeAddr := range c.hashMap {
		if nodeAddr == addr {
			hashesToRemove = append(hashesToRemove, hash)
		}
	}

	// 如果没有找到该节点，返回错误
	if len(hashesToRemove) == 0 {
		return errors.New("node not found")
	}

	// 从keys中删除这些哈希值
	for _, hash := range hashesToRemove {
		index := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] == hash })
		if index < len(c.keys) && c.keys[index] == hash {
			c.keys = append(c.keys[:index], c.keys[index+1:]...)
		}
	}

	// 从hashMap中删除这些哈希值
	for _, hash := range hashesToRemove {
		delete(c.hashMap, hash)
	}

	return nil
}

// Get 根据数据找到对应的节点
// 参数：数据（如缓存键）
// 返回值：对应的节点地址和可能的错误
func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("no nodes available")
	}

	hash := c.hash([]byte(key)) // 计算数据的哈希值

	// 使用二分查找找到第一个 >= 数据哈希值的节点哈希值
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })

	// 如果数据哈希值比所有节点都大，则选择第一个节点（环状结构）
	if idx == len(c.keys) {
		idx = 0
	}

	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.hashMap[c.keys[idx]], nil
}

// SetHealthCheck 设置健康检查函数（预留接口）
// 参数：健康检查函数，返回true表示节点健康
// 返回值：无
func (c *ConsistentHashBalance) SetHealthCheck(healthCheck func(string) bool) {
	// 预留接口，实际实现可以在这里添加健康检查逻辑
	// 例如定期检查节点健康状态，自动移除不健康的节点
}

// GetNodes 获取所有节点信息（用于监控或调试）
// 返回值：节点地址列表
func (c *ConsistentHashBalance) GetNodes() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()

	nodes := make([]string, 0, len(c.hashMap))
	for _, addr := range c.hashMap {
		nodes = append(nodes, addr)
	}
	return nodes
}

// GetNodeHashes 获取所有节点的哈希值（用于调试）
// 返回值：节点哈希值列表
func (c *ConsistentHashBalance) GetNodeHashes() []uint32 {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.keys
}
