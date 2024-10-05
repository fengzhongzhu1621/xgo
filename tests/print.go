package tests

import (
	"fmt"

	"github.com/dablelv/cyan/encoding"
)

// PrintStruct 打印结构体
func PrintStruct(obj interface{}) {
	s, _ := encoding.ToIndentJSON(obj)
	fmt.Printf("%v\n", s)
}

// ToString 将结构体转换为字符串表示
func ToString(obj interface{}) string {
	s, _ := encoding.ToIndentJSON(obj)
	result := fmt.Sprintf("%v", s)
	return result
}
