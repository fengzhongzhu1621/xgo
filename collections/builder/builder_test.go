package builder

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

type PersonBuilder struct {
	name    string
	age     int
	address string
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{}
}

func (pb *PersonBuilder) WithName(name string) *PersonBuilder {
	pb.name = name
	return pb
}

func (pb *PersonBuilder) WithAge(age int) *PersonBuilder {
	pb.age = age
	return pb
}

func (pb *PersonBuilder) WithAddress(address string) *PersonBuilder {
	pb.address = address
	return pb
}

func (pb *PersonBuilder) Build() *Person {
	return &Person{
		name:    pb.name,
		age:     pb.age,
		address: pb.address,
	}
}

func TestBuilder(t *testing.T) {
	builder := NewPersonBuilder()
	builder.WithAge(30)
	builder.WithName("John Doe")
	builder.WithAddress("Beijing")

	newPerson := builder.Build()

	// 打印创建的Person
	fmt.Printf("Person: %+v\n", newPerson)
}
