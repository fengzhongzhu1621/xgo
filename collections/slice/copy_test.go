package slice

import (
	"fmt"
	"testing"
)

func TestCopy(t *testing.T) {
	// 创建一个整型数组
	// [...]int 表示声明一个整型数组。这里的 ... 是一种特殊的用法，它让编译器根据后面初始化列表中的元素个数来确定数组的长度
	arr := [...]int{1, 2, 3, 4, 5}

	s1 := make([]int, 3)
	copy(s1, arr[2:5]) // s1 现在是 arr[2:5] 的副本

	s2 := make([]int, 3)
	copy(s2, arr[1:4])

	s1[0] = 99

	fmt.Println("arr:", arr) // [1 2 3 4 5]
	// 对 s1 的修改不会对 arr 和 s2 造成影响
	fmt.Println("s1:", s1) // [99 4 5]
	fmt.Println("s2:", s2) // [2 3 4]
}
