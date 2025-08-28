package validator

import "reflect"

// IsNil 判断对象是否为 nil
// * nil 是 Golang 中的一个预定义标识符，用于表示指针、接口、切片、映射、通道和函数类型的零值。当你想要检查一个变量是否为其类型的零值时，可以使用 == 或 != 运算符与 nil 进行比较。
// * IsNil 是一个函数，通常用于判断接口类型的变量是否为 nil。由于接口类型的变量可以存储任意类型的值，因此在使用 == 或 != 运算符与 nil 进行比较时可能会出现误判。
// Note: 注意的是，IsNil 函数只能用于接口类型的变量，对于其他类型的变量，仍然需要使用 == 或 != 运算符与 nil 进行比较。
func IsNil(value any) bool {
	defer func() { recover() }()
	// nil 用于直接判断变量是否为其类型的零值，而 IsNil 函数用于判断接口类型的变量是否为 nil
	return value == nil || reflect.ValueOf(value).IsNil()
}
