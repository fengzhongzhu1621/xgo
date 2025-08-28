package reflectutils

import "reflect"

// New create a new interface which has the same underlying type as v.
// 创建一个新的接口值，这个新接口值的底层类型（underlying type）与传入参数 v 的底层类型相同，但它的值是一个新的、零值的实例。
// https://github.com/trpc-ecosystem/go-utils/blob/main/copyutils/new.go
func New(v interface{}) interface{} {
	// 避免对 nil 值进行反射操作，因为反射操作在 nil 上可能会导致 panic。
	if v == nil {
		return nil
	}

	// reflect.New 接收一个 reflect.Type 类型的参数，并返回一个新的指针，这个指针指向一个该类型的零值实例。
	return reflect.New(
		// reflect.Indirect 的作用是获取一个值的“间接”值。
		// 如果传入的值是一个指针（*T），reflect.Indirect 会解引用它，返回指向的实际值的 reflect.Value。
		// 如果传入的值本身不是指针，则直接返回它自己。
		// 这一步的目的是确保我们拿到的是值的类型信息，而不是指针的类型信息。
		reflect.Indirect(
			reflect.ValueOf(v), // 使用反射包中的 reflect.ValueOf 获取 v 的反射值（reflect.Value 类型）
		).Type(), // 从 reflect.Value 中提取其底层类型（reflect.Type）。这个类型就是 v 最终指向的值的类型（去掉了指针层级）。
	).Interface() // 将 reflect.Value 转换回普通的 interface{} 类型，以便作为函数的返回值。
}

func CloneOrCreate(v interface{}) interface{} {
	// 使用 New 创建一个同类型的新实例
	newV := New(v)
	// 如果需要深拷贝，可以在这里进一步复制字段值
	return newV
}
