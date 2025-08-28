//go:build !purego
// +build !purego

package reflectutils

import (
	"reflect"
	"unsafe"
)

// Pointer is an opaque typed pointer and is guaranteed to be comparable.
type Pointer struct {
	// 这是一个不安全的指针，它可以指向任何类型的值。unsafe.Pointer 是 Go 语言中一个特殊的类型，允许绕过类型系统的限制，进行底层内存操作。
	p unsafe.Pointer
	// 用于存储指针所指向的具体类型信息。通过保存类型信息，可以在后续操作中了解指针指向的数据类型，从而进行类型安全的操作。
	t reflect.Type
}

// PointerOf returns a Pointer from v, which must be a
// reflect.Ptr, reflect.Slice, or reflect.Map.
// 将 reflect.Value 转换为一个包含底层指针和类型信息的 Pointer。
// v reflect.Value: 必须是一个指针、切片或映射类型的 reflect.Value。这是因为只有这些类型才具有底层的指针表示。
func PointerOf(v reflect.Value) Pointer {
	// The proper representation of a pointer is unsafe.Pointer,
	// which is necessary if the GC ever uses a moving collector.

	return Pointer{unsafe.Pointer(v.Pointer()), v.Type()}
}

// IsNil reports whether the pointer is nil.
// 检查 Pointer 结构体中的指针是否为 nil
func (p Pointer) IsNil() bool {
	return p.p == nil
}

// Uintptr returns the pointer as an uintptr.
// 将 Pointer 结构体中的 unsafe.Pointer 转换为 uintptr 类型
// 当需要进行低级内存操作或将指针与其他需要 uintptr 类型的 API 交互时使用。
func (p Pointer) Uintptr() uintptr {
	return uintptr(p.p)
}
