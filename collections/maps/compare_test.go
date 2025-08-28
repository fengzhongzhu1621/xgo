package maps

import (
	"fmt"
	"maps"
	"testing"
)

// func Equal[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool
// 判断两个 map 是否包含相同的键值对，内部使用 == 进行比较
// Note: map 类型的 key 和 value 必须是 comparable 类型
func TestEqual(t *testing.T) {
	m1 := map[int]int{0: 0, 1: 1, 2: 2}
	m2 := map[int]int{0: 0, 1: 1}
	m3 := map[int]int{0: 0, 1: 1, 2: 2}

	fmt.Println(maps.Equal(m1, m2)) // false
	fmt.Println(maps.Equal(m1, m3)) // true
}

// func EqualFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](m1 M1, m2 M2, eq func(V1, V2) bool) bool
// 类似 Equal 函数，只不过是通过 eq 函数进行比较值，键仍使用 == 进行比较。注意：value 可以为任意类型（any）。
func TestEqualFunc(t *testing.T) {
	type User struct {
		Nickname string
		IdCard   string
	}

	m1 := map[int]User{0: {Nickname: "李四", IdCard: "111"}, 1: {Nickname: "张三", IdCard: "222"}}
	m2 := map[int]User{0: {Nickname: "李四", IdCard: "111"}}
	m3 := map[int]User{0: {Nickname: "王五", IdCard: "111"}, 1: {Nickname: "张三", IdCard: "222"}}

	fmt.Println(maps.EqualFunc(m1, m2, func(user User, user2 User) bool {
		return user.IdCard == user2.IdCard
	})) // false
	fmt.Println(maps.EqualFunc(m1, m3, func(user User, user2 User) bool {
		return user.IdCard == user2.IdCard
	})) // true
}
