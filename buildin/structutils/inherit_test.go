package structutils

import (
	"fmt"
	"testing"
)

// ////////////////////////////////////////////////////////////////
type Parent struct {
	// 父类字段和方法
	name string
}

func (p *Parent) test() {}

// ////////////////////////////////////////////////////////////////
// 内嵌结构体：模拟继承，直接访问 Parent 字段和方法
// 内嵌不继承私有字段和方法
// 内嵌采取就近访问,非线性查找
type Child struct {
	Parent // 继承Parent
	// 子类字段和方法
}

// 覆盖Parent的test方法
// 子类可以定义与父类同名的方法,这相当于方法 Override
func (c *Child) test() {
	// 先调用Parent的实现
	c.Parent.test()

	// Child的处理逻辑
}

// ////////////////////////////////////////////////////////////////
// 组合结构体：组合需要通过实例名访问
type Child2 struct {
	parent Parent
}

// ////////////////////////////////////////////////////////////////
// 多重继承
type (
	A struct{}
	B struct{}
	C struct{ A B }
)

// ////////////////////////////////////////////////////////////////
// 测试构造函数
func TestConstruct(t *testing.T) {
	c := Child{
		Parent{name: "Tom"},
	}
	fmt.Println(c)
}

// ////////////////////////////////////////////////////////////////
// 装饰器模式
type Shape interface {
	Draw()
}

type Circle struct {
	radius float32
}

func (c *Circle) Draw() {
	fmt.Printf("Draw circle, radius=%f\n", c.radius)
}

type ColoredCircle struct {
	Circle
	color string
}

func (c *ColoredCircle) Draw() {
	fmt.Printf("Draw colored circle, color=%s\n", c.color)
	c.Circle.Draw()
}

func TestDecorator(t *testing.T) {
	cc := &ColoredCircle{
		Circle: Circle{10},
		color:  "Red",
	}
	cc.Draw()
}

// ////////////////////////////////////////////////////////////////
// 多态
type Sayer interface {
	Say()
}

type Dog struct{}

func (d *Dog) Say() {
	println("Wang!")
}

type Cat struct{}

func (c *Cat) Say() {
	println("Miao!")
}

func TestPolymorphism(t *testing.T) {
	animals := []Sayer{
		&Dog{},
		&Cat{},
	}

	for _, animal := range animals {
		animal.Say()
	}
}
