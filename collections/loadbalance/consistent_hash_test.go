package loadbalance

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/tests"
)

func TestConsistentHashBalance(t *testing.T) {
	// 创建一致性哈希负载均衡器，每个节点生成100个虚拟节点（权重高）
	ch := NewConsistentHashBalance(100, nil)

	// 添加节点
	ch.Add("192.168.1.1")
	ch.Add("192.168.1.2")
	ch.Add("192.168.1.3")

	// 获取数据对应的节点
	node, err := ch.Get("some_data_key")
	if err != nil {
		panic(err)
	}
	println("Data should go to node:", node)

	// 获取所有节点
	nodes := ch.GetNodes()
	println("All nodes:", tests.ToString(nodes))

	// 删除节点
	ch.Remove("192.168.1.2")

	// 再次获取数据对应的节点
	node, err = ch.Get("some_data_key")
	if err != nil {
		panic(err)
	}
	println("After removal, data should go to node:", node)
}
