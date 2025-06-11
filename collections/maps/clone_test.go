package maps

import (
	"fmt"
	"maps"
	"testing"
)

// func Clone[M ~map[K]V, K comparable, V any](m M) M
// Clone 函数接收一个 m 参数，该函数的功能是返回 m 的副本，底层基于浅层克隆去实现，使用普通赋值的方式去设置新的键值对。
// 修改克隆后的 map 任意 key 的 value 将有可能影响原 map 的 value。
func TestClone(t *testing.T) {
	type Programmer struct {
		Name string
		City string
	}

	m1 := map[string]Programmer{
		"programmer-01": {Name: "李四", City: "深圳"},
		"programmer-02": {Name: "张三", City: "广州"},
	}
	m2 := maps.Clone(m1)
	fmt.Printf("m1: %v\n", m1)
	fmt.Printf("m2: %v\n", m2)
	// m1: map[programmer-01:{李四 深圳} programmer-02:{张三 广州}]
	// m2: map[programmer-01:{李四 深圳} programmer-02:{张三 广州}]

	// value 是指针类型。
	// 如果 m1 的 value 是指针类型，那么在对克隆后的 m2 中的任意 key 对应的 value 进行修改操作后，都会直接影响到 m1。
	// 这是因为 m1 和 m2 共享了同一组指向相同 Programmer 结构体的指针，因此对一个指针的修改会在两个 map 中都可见。
	type Programmer2 struct {
		Name string
		City string
	}

	m3 := map[string]*Programmer2{
		"programmer-01": {Name: "李四", City: "深圳"},
		"programmer-02": {Name: "张三", City: "广州"},
	}
	fmt.Printf("m3: %v, %v\n", *m3["programmer-01"], *m3["programmer-02"])
	m4 := maps.Clone(m3)
	fmt.Printf("m4: %v, %v\n", *m4["programmer-01"], *m4["programmer-02"])
	m4["programmer-02"].City = "海口"
	fmt.Printf("m4 修改后，m3: %v, %v\n", *m3["programmer-01"], *m3["programmer-02"])
	fmt.Printf("m4 修改后，m4: %v, %v\n", *m4["programmer-01"], *m4["programmer-02"])
	// m4 修改后，m3: {李四 深圳}, {张三 海口}
	// m4 修改后，m4: {李四 深圳}, {张三 海口}
}
