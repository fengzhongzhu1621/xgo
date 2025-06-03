package uuid

// shortid 是一种用于生成简短、唯一 ID 的工具，适合需要短链或验证码的场景。生成的 ID 是 7-14 个字符，包含字母、数字、特殊字符。
import (
	"fmt"
	"testing"

	"github.com/teris-io/shortid"
)

// 唯一性：基于随机数生成，具有较高的唯一性。
// 无序性：生成的 ID 无时间排序特性。
// 可读性：短 ID 的长度适中，便于展示
//
// 不具备排序特性，不能用于需要时间顺序的场景。
//
// 短链接、验证码、临时文件名等场景。
func TestShortId(t *testing.T) {
	id, _ := shortid.Generate()
	fmt.Println(id) // 输出类似 "dppUrjK3"
}
