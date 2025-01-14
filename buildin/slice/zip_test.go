package slice

import (
	"fmt"
	"testing"

	"github.com/araujo88/lambda-go/pkg/tuple"
)

// tuple 包提供了一种通用的 tuple 数据结构，可用于处理成对数据。
// 当你需要从函数中返回多个值或将相关数据分组时，这尤其方便。
// 使用 Zip 函数将两个片段合并为一个元组片段的示例
func TestZip(t *testing.T) {
	names := []string{"Alice", "Bob", "Charlie"}
	ages := []int{25, 30, 35}

	pairs := tuple.Zip(names, ages)
	for _, pair := range pairs {
		fmt.Printf("%s is %d years old\n", pair.First, pair.Second)
	}
}
