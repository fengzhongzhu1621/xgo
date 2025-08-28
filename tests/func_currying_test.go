package tests

// 测试函数的基本功能

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// add 两个整数相加
func add(x int, y int) int {
	return x + y
}

func aggregate(a, b, c int, arithmetic func(int, int) int) int {
	return arithmetic(arithmetic(a, b), c)
}

// partial 装饰器函数，将一个函数装为另一个函数
func partial(mathFunc func(int, int) int) func(int) int {
	return func(i int) int {
		return mathFunc(i, i)
	}
}

// TestFunctionCurrying 测试函数作为参数传递
func TestFunctionCurrying(t *testing.T) {
	actual := aggregate(1, 2, 3, add)
	expect := 6

	assert.Equal(t, expect, actual)
}

// TestDecorator 测试装饰器模式
func TestDecorator(t *testing.T) {
	actual := partial(add)(1)
	expect := 2
	assert.Equal(t, expect, actual)
}

// concatter 返回一个基于装饰模式的闭包
func concatter() func(string) string {
	// 多次调用会复用
	value := ""

	// 返回一个闭包
	return func(s string) string {
		value += s + " "
		return value
	}
}

// TestClosure 测试闭包
func TestClosure(t *testing.T) {
	closure := concatter()
	assert.Equal(t, "hello ", closure("hello"))
	assert.Equal(t, "hello world ", closure("world"))
}

// 测试匿名函数
func TestAnonymous(t *testing.T) {
	actual := aggregate(1, 2, 3, func(x, y int) int { return x * y })
	expect := 6
	assert.Equal(t, expect, actual)
}
