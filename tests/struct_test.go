package tests

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

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
