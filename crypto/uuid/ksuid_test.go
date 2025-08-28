package uuid

// ksuid 是一种可排序的唯一 ID，27 字节的 Base62 编码字符串，包含时间戳信息，方便时间排序。
import (
	"fmt"
	"testing"

	"github.com/segmentio/ksuid"
)

// 唯一性：基于时间戳、随机数生成，全局唯一。
// 排序性：具备时间排序特性，适合检索和排序。
// 可读性：27 字符的 Base62 字符串，适合展示和日志记录。
func TestKsuid(t *testing.T) {
	id := ksuid.New()
	fmt.Println(id.String()) // 输出类似 "1CHzCEl82ZZe2r2KS1qEz3XF8Ve"
}
