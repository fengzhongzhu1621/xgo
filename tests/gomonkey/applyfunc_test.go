package gomonkey

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func add(a, b int) int {
	return a + b
}

func TestApplyFunc(t *testing.T) {
	// 模拟 add 函数，需要在测试结束时调用 Reset 方法，以恢复原始函数。否则，后续的测试可能会受到影响。
	patches := gomonkey.ApplyFunc(add, func(a, b int) int {
		return a * b
	})
	defer patches.Reset()

	result := add(2, 3)
	if !reflect.DeepEqual(result, 6) {
		t.Errorf("add(2, 3) = %d; want 6", result)
	}
}

type Host struct {
	IP   string
	Name string
}

func Convert2Json(h *Host) (string, error) {
	b, err := json.Marshal(h)
	return string(b), err
}

func TestConvert2Json(t *testing.T) {
	patches := gomonkey.ApplyFunc(json.Marshal, func(v interface{}) ([]byte, error) {
		return []byte(`{"IP":"1.1.1.1","Name":"Sky"}`), nil
	})

	defer patches.Reset()

	h := Host{Name: "xx", IP: "2.2.2.2"}
	s, err := Convert2Json(&h)

	expectedString := `{"IP":"1.1.1.1","Name":"Sky"}`

	if s != expectedString || err != nil {
		t.Errorf("expected %v, got %v", expectedString, s)
	}

}
