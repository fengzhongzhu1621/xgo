package gomonkey

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

type MyStruct struct{}

func (s *MyStruct) MyMethod() int {
	return 42
}

func (m *MyStruct) MyMethod2(a int) (int, error) {
	return a * 2, nil
}

func TestApplyMethodReturn(t *testing.T) {
	myStruct := &MyStruct{}
	patches := gomonkey.ApplyMethodReturn(myStruct, "MyMethod", 0)
	defer patches.Reset()

	result := myStruct.MyMethod()
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestApplyMethodReturn2(t *testing.T) {
	var c *MyStruct
	patches := gomonkey.ApplyMethodReturn(c, "MyMethod", 0)
	defer patches.Reset()

	myStruct := &MyStruct{}
	result := myStruct.MyMethod()
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestApplyMethodReturn3(t *testing.T) {
	myStruct := &MyStruct{}

	// 使用 gomonkey.ApplyMethodReturn 模拟 MyMethod 方法的返回值
	patches := gomonkey.ApplyMethodReturn(myStruct, "MyMethod2", 42, nil)
	defer patches.Reset()

	// 调用 MyMethod 方法，它现在将返回模拟的值 42
	result, _ := myStruct.MyMethod2(21)
	assert.Equal(t, 42, result)
}
