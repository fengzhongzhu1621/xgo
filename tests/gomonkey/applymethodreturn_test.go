package gomonkey

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

type MyStruct struct{}

func (s *MyStruct) MyMethod() int {
	return 42
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
