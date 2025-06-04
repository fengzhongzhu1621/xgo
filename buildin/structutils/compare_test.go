package structutils

import (
	"fmt"
	"reflect"
	"testing"
)

type Person struct {
	Name string
}

// 判断两个变量是否引用同一个对象（内存地址相同）
func TestCompareSameObject(t *testing.T) {
	p1 := &Person{Name: "Alice"}
	p2 := p1                     // p2 和 p1 指向同一个对象
	p3 := &Person{Name: "Alice"} // p3 是另一个新对象

	fmt.Println(p1 == p2) // true（指向同一个对象）
	fmt.Println(p1 == p3) // false（指向不同对象）
}

// 使用 reflect.DeepEqual 深度比较两个对象（值相同，可以不是同一个对象）
type BadPerson struct {
	Friends []string
}

func TestDeepEqual(t *testing.T) {
	p1 := BadPerson{Friends: []string{"Alice"}}
	p2 := BadPerson{Friends: []string{"Alice"}}
	p3 := BadPerson{Friends: []string{"Bob"}}

	fmt.Println(reflect.DeepEqual(p1, p2)) // true（内容相同）
	fmt.Println(reflect.DeepEqual(p1, p3)) // false（内容不同）
}

type Person2 struct {
	Name string
	Age  int
}

func (p Person2) Equal(other Person2) bool {
	return p.Name == other.Name && p.Age == other.Age
}

// 自定义比较逻辑（实现 Equal 方法）
// 如果对象比较逻辑复杂（如需要忽略某些字段），可以自定义 Equal 方法
func TestCustomEqual(t *testing.T) {
	p1 := Person2{Name: "Alice", Age: 30}
	p2 := Person2{Name: "Alice", Age: 30}
	p3 := Person2{Name: "Bob", Age: 30}

	fmt.Println(p1.Equal(p2)) // true
	fmt.Println(p1.Equal(p3)) // false
}
