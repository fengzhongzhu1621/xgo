package uuid

// sonyflake 是一种高效、分布式的唯一 ID 生成算法，生成的 ID 是 64 位整数，由时间戳、机器 ID 和序列号组成。
import (
	"fmt"
	"testing"

	"github.com/sony/sonyflake"
)

// 唯一性：基于时间戳、机器 ID 和序列号，确保唯一性。
// 排序性：具备时间排序特性。
// 可扩展性：支持自定义机器 ID 位数和时间戳精度。
func TestSonyflake(t *testing.T) {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, _ := sf.NextID()
	fmt.Println(id) // 输出类似 "173301874793540608"
}
