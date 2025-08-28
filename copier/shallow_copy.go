package copier

import (
	"errors"
	"fmt"
	"reflect"
)

// ShallowCopy copies src to dst with only one copy depth.
//
// The native CopyTo method will be used if src implements one.
//
// ShallowCopy does copy unexported fields.
// dst must be a pointer type.
// 将 src 的值浅拷贝到 dst 中，要求 dst 必须是一个指针类型，并且 dst 和 src 的底层类型必须一致。
func ShallowCopy(dst, src interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("shallow copy should never panic: %v", r)
		}
	}()

	if src == nil {
		return nil
	}

	// 使用了 Go 的接口类型断言，检查 src 是否实现了 CopierTo 接口。
	// 如果实现了，就调用 copierTo.CopyTo(dst) 方法，将拷贝的控制权交给 src 自己。
	// 这是一种扩展机制，允许某些特殊类型的对象自定义拷贝行为（例如包含不可导出字段或需要特殊处理的类型）。
	if copierTo, ok := src.(CopierTo); ok {
		return copierTo.CopyTo(dst)
	}

	// 检查 dst 是否是指针类型
	dstV := reflect.ValueOf(dst)
	if dstV.Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer")
	}

	// 解引用 dst 和 src
	// src 可能是指针也可能是非指针，统一解引用后可以简化后续的类型比较和赋值操作。
	dstV = reflect.Indirect(dstV)
	srcV := reflect.Indirect(reflect.ValueOf(src))

	// 检查 dst 和 src 的类型是否匹配
	if dstV.Type() != srcV.Type() {
		return fmt.Errorf("dst %s and src %s type miss match", dstV.Type(), srcV.Type())
	}

	// 使用 reflect.Value 的 Set 方法，将 srcV 的值直接复制到 dstV 中。
	// 这是一种浅拷贝操作：
	// 如果 srcV 是基本类型（如 int、string），会直接复制值。
	// 如果 srcV 是引用类型（如 *int、[]int、map[string]int），只会复制指针/引用本身，而不会递归地复制内部数据。
	dstV.Set(srcV)

	return nil
}
