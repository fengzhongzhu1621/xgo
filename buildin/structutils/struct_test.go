package structutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

type User2 struct {
	Name string
	Age  int
	_    struct{}
}

// 强制要求结构体在初始化时必须使用具名字段初始化
func TestDefine(t *testing.T) {
	user := User{"name", 10}
	user = User{Name: "name", Age: 20}
	assert.Equal(t, user.Age, 20)

	// _ 的作用，可以强制要求结构体在初始化时必须使用具名字段初始化（声明零值结构体变量的场景除外）
	// 优点
	// 1. 代码可读性：具名字段初始化使得代码可读性和可维护性更好
	// 2. 避免错误： 位置初始化需要严格遵循字段顺序，容易出错
	// user = User2{"name", 10} // 编译错误
	// user = User2{"name", 10, struct{}{}} // 编译错误
	user2 := User2{}
	user2 = User2{Name: "name", Age: 30}
	assert.Equal(t, user2.Age, 30)
}
