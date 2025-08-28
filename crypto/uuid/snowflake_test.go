package uuid

import (
	"fmt"
	"testing"

	"github.com/bwmarrin/snowflake"
)

// 唯一性：由时间戳、机器 ID 和序列号组成，保证了全局唯一性。
// 排序性：ID 基于时间戳生成，具备时间排序特性。
// 性能：生成速度非常快，适用于高并发环境。
//
// ID 紧凑，适合用于数据库的主键或分布式系统的唯一 ID。
// 生成的 ID 有序，方便插入有序数据库。
//
// 需要配置机器节点 ID，部署时稍微复杂。
func TestSnowflake(t *testing.T) {
	node, _ := snowflake.NewNode(1)
	id := node.Generate()
	fmt.Println(id)         // 输出类似 "1234567890123456789"
	fmt.Println(id.Int64()) // 获取 int64 类型的 ID
}
