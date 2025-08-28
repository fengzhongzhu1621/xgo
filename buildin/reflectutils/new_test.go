package reflectutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type S struct{}

func TestNew(t *testing.T) {
	// p 是一个 *Person 类型的指针。
	p := &Person{Name: "Alice", Age: 30}

	// New(p) 会：
	// 通过 reflect.ValueOf(p) 得到 p 的反射值。
	// 使用 reflect.Indirect 解引用，得到 Person 类型的反射值。
	// 提取 Person 的类型信息。
	// 使用 reflect.New 创建一个新的 *Person 指针，指向一个零值的 Person 结构体。
	// 最后将其转换为 interface{} 返回。
	newP := New(p)
	fmt.Printf("%T\n", newP) // 输出: *main.Person

	// 所以 newP 是一个新的 *Person 类型的指针，指向一个零值的 Person（即 Name 和 Age 都是零值）。
	fmt.Println(reflect.ValueOf(newP).IsNil()) // 输出: false

	i := New(3)
	require.Equal(t, reflect.Ptr, reflect.TypeOf(i).Kind())
	require.Equal(t, reflect.Int, reflect.TypeOf(i).Elem().Kind())

	s := New(S{})
	require.Equal(t, reflect.Ptr, reflect.TypeOf(s).Kind())
	require.Equal(t, reflect.TypeOf(S{}), reflect.TypeOf(s).Elem())

	s = New(&S{})
	require.Equal(t, reflect.Ptr, reflect.TypeOf(s).Kind())
	require.Equal(t, reflect.TypeOf(S{}), reflect.TypeOf(s).Elem())
}

type Service interface {
	DoSomething()
}

type MyService struct{}

func (s *MyService) DoSomething() {
	fmt.Println("Doing something")
}

// 允许框架在不明确知道具体类型的情况下，动态创建新的服务对象。
func CreateServiceFromExample(example interface{}) interface{} {
	// 假设 example 是某个实现了 Service 的实例
	return New(example)
}

// 序列化/反序列化辅助工具，根据已反序列化对象的类型，创建一个新的空对象，然后逐步填充字段值。
func DeserializeInto(data []byte, example interface{}) (interface{}, error) {
	// 假设我们已经解析了 data 的结构，知道了类型
	// 使用 example 的类型信息创建一个新的实例
	newobj := New(example)
	// 然后可以进一步将 data 的内容填充到 newobj 中（可能需要反射）
	return newobj, nil
}
