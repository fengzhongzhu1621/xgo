package maps

import (
	"fmt"
	"maps"
	"testing"
)

// func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2)
// 复制 src 中的所有键值对到 dst 中，如果 dst 中包含 src 中的任意 key，则该 key 对应的 value 将会被覆盖
func TestCopy(t *testing.T) {
	m1 := map[string]string{"Name": "张三", "City": "深圳"}
	m2 := map[string]string{"City": "广州", "Phone": "xxxxxx"}

	maps.Copy(m1, m2)

	fmt.Println(m1) // map[City:广州 Name:张三 Phone:xxxxxx]
}
