package option

import (
	"fmt"
	"testing"
)

// Person 结构体定义
type Person struct {
	name    string
	age     int
	address string
}

// Option 是一个函数类型，接收一个Person指针并修改它
type Option func(*Person)

// NewPerson 是构造函数，接收可变数量的Option参数
func NewPerson(opts ...Option) *Person {
	p := &Person{} // 创建一个新的Person实例
	for _, opt := range opts {
		opt(p) // 依次应用每个Option
	}
	return p // 返回构建好的Person
}

// WithName 是一个Option，设置Person的name字段
func WithName(name string) Option {
	return func(person *Person) {
		person.name = name
	}
}

// WithAge 是一个Option，设置Person的age字段
func WithAge(age int) Option {
	return func(person *Person) {
		person.age = age
	}
}

// WithAddress 是一个Option，设置Person的address字段
func WithAddress(addr string) Option {
	return func(person *Person) {
		person.address = addr
	}
}

func TestOption(t *testing.T) {
	// 使用构建器模式创建Person实例
	newPerson := NewPerson(
		WithName("John Doe"),
		WithAge(30),
		WithAddress("Beijing"),
	)

	// 打印创建的Person
	fmt.Printf("Person: %+v\n", newPerson)
}
