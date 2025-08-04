package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// 结构体中使用 _，强制要求该结构体在初始化时必须使用具名字段初始化
type User struct {
	Name string
	Age  int
	// 在结构体中定义一个名为 _ 的字段，可以强制要求该结构体在初始化时必须使用具名字段初始化（声明零值结构体变量的场景除外）
	_ struct{}
}

func TestStructSuccess(t *testing.T) {
	user := User{}
	user = User{Name: "foo", Age: 18}
	user = User{"bar", 18, struct{}{}}
	assert.Equal(t, user.Age, 18)
}

func TestStructFailure(t *testing.T) {
	// 编译错误 too few values in struct literal of type User
	// _ = User{"陈明勇", 18}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// 禁止结构体比较
// A Value can represent any Go value, but unlike type any,
// it can represent most small values without an allocation.
// The zero Value corresponds to nil.
//
// 可以发现，这里有一个匿名字段 _ [0]func()，并且注释写着 // disallow ==。
// _ [0]func() 的目的显然是为了禁止比较。
// 既然使用 map[string]int 和 _ [0]func() 都能实现禁止结构体相等性比较，那么我为什么说 _ [0]func() 是更优雅的做法呢？
// _ [0]func() 有着比其他实现方式更优的特点：
// 它不占内存空间！
// 使用匿名字段 _ 语义也更强。
// 不过值得注意的是：当使用 _ [0]func() 时，不要把它放在结构体最后一个字段，推荐放在第一个字段。这与结构体内存对齐有关
type Value struct {
	_ [0]func() // disallow ==
	// num holds the value for Kinds Int64, Uint64, Float64, Bool and Duration,
	// the string length for KindString, and nanoseconds since the epoch for KindTime.
	num uint64
	// If any is of type Kind, then the value is in num as described above.
	// If any is of type *time.Location, then the Kind is Time and time.Time value
	// can be constructed from the Unix nanos in num and the location (monotonic time
	// is not preserved).
	// If any is of type stringptr, then the Kind is String and the string value
	// consists of the length in num and the pointer in any.
	// Otherwise, the Kind is Any and any is the value.
	// (This implies that Attrs cannot store values of type Kind, *time.Location
	// or stringptr.)
	any any
}

// 结构体是否可以比较，不取决于字段是否可导出，而是取决于其是否包含不可比较字段。
// 如果全部字段都是可比较的，那么这个结构体就是可比较的。
// 如果其中有一个字段不可比较，那么这个结构体就是不可比较的。
// 不过虽然我们不可以使用 == 对 n1、n2 进行比较，但我们可以使用 reflect.DeepEqual() 对二者进行比较
// fmt.Println(reflect.DeepEqual(n1, n2))
