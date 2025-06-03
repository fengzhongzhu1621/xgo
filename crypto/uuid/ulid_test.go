package uuid

// ULID（Universally Unique Lexicographically Sortable Identifier）是一种可排序的唯一 ID 生成方案。它的目标是提供一种比 UUID 更具可读性和排序性的唯一 ID。
import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
)

// 唯一性：基于时间戳和随机数的组合，保证全局唯一性。
// 排序性：ID 基于时间戳生成，具备排序特性。
// 可读性：使用 26 个字符的 Base32 编码，较 UUID 更紧凑。
//
// 具备时间有序性，方便插入有序数据库。
// 紧凑且易读，适合展示和日志记录。
//
// 在极高并发场景下，存在理论上的碰撞可能。
func TestUlid(t *testing.T) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	fmt.Println(id.String()) // 输出类似 "01FHZB1E8B9B8YZW6CVDBKMT6T"
}
