package structutils

import (
	"fmt"
	"testing"
	"unsafe"
)

type Inefficient struct {
	a bool  // 1 byte，7 个字节的填充，用于对齐 b
	b int64 // 8 bytes
	c bool  // 1 byte，7 个字节的填充，以保持对齐
}

type Efficient struct {
	b int64 // 8 bytes
	a bool  // 1 byte
	c bool  // 1 byte
	// 6 个字节的填充
}

// 计算结构体大小
func TestAlignment(t *testing.T) {
	inefficient := Inefficient{}
	efficient := Efficient{}

	// Inefficient: 24 bytes
	fmt.Printf("Inefficient: %d bytes\n", unsafe.Sizeof(inefficient))
	// Efficient: 16 bytes
	fmt.Printf("Efficient: %d bytes\n", unsafe.Sizeof(efficient))
}
