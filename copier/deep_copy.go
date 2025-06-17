package copier

import (
	"reflect"
)

// DeepCopy deeply copy v.
//
// The native Copy method will be used if v implements one.
//
// Pointer loops are preserved in new value.
// Slices sharing the same underlying array do not share in new value.
// Channels are directly copied.
// Unexported fields are not copied.
// 实现了一个功能完整的 深拷贝（Deep Copy） 工具函数 DeepCopy，
// 它使用 Go 的反射（reflect）机制递归地复制任意类型的值，
// 同时正确处理了指针循环引用、切片共享底层数组、通道直接复制、未导出字段不复制等问题。
func DeepCopy(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	if copier, ok := v.(Copier); ok {
		return copier.Copy()
	}

	return copy(reflect.ValueOf(v), make(map[uintptr]reflect.Value)).Interface(), nil
}

func copy(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Array:
		return copyArray(v, visited)
	case reflect.Slice:
		return copySlice(v, visited)
	case reflect.Map:
		return copyMap(v, visited)
	case reflect.Ptr:
		return copyPointer(v, visited)
	case reflect.Struct:
		return copyStruct(v, visited)
	case reflect.Interface:
		return copyInterface(v, visited)
	default:
		return copyDefault(v, visited)
	}
}

func copyArray(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	l, t := v.Len(), v.Type()

	// 创建一个新的数组 vv，类型和原数组一致。
	vv := reflect.New(t).Elem()
	for i := 0; i < l; i++ {
		// 遍历原数组的每个元素，递归调用 copy 深拷贝每个元素，并设置到新数组中。
		// 数组的每个元素都会被深拷贝，新数组和原数组完全独立。
		vv.Index(i).Set(copy(v.Index(i), visited))
	}

	return vv
}

func copySlice(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	l, t := v.Len(), v.Type()
	// 创建一个新的切片 vv，长度和容量与原切片一致。
	vv := reflect.MakeSlice(t, l, l)
	for i := 0; i < l; i++ {
		vv.Index(i).Set(copy(v.Index(i), visited))
	}

	return vv
}

func copyMap(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	// 如果已经拷贝过这个映射（指针相同），直接返回之前拷贝的结果（避免循环引用导致的重复拷贝）
	if vv, ok := visited[v.Pointer()]; ok {
		return vv
	}

	// 创建一个新的映射 vv，类型和原映射一致。
	vv := reflect.MakeMap(v.Type())
	visited[v.Pointer()] = vv

	iter := v.MapRange()
	for iter.Next() {
		vv.SetMapIndex(
			copy(iter.Key(), visited),
			copy(iter.Value(), visited))
	}

	return vv
}

func copyPointer(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	if vv, ok := visited[v.Pointer()]; ok {
		return vv
	}

	// 创建一个新的指针 vv，指向一个新分配的底层值（类型和原指针一致）。
	vv := reflect.New(v.Type().Elem())
	visited[v.Pointer()] = vv
	vv.Elem().Set(copy(v.Elem(), visited))

	return vv
}

func copyStruct(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	t := v.Type()
	vv := reflect.New(t)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			// 如果字段是未导出的（f.PkgPath != ""），跳过不拷贝。
			continue
		}
		vv.Elem().Field(i).Set(copy(v.Field(i), visited))
	}

	return vv.Elem()
}

func copyInterface(v reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if v.IsNil() {
		return v
	}

	vv := reflect.New(v.Type()).Elem()
	vv.Set(copy(v.Elem(), visited))

	return vv
}

func copyDefault(v reflect.Value, _ map[uintptr]reflect.Value) reflect.Value {
	// 直接返回原值（不拷贝）
	// 通道是引用类型，无法深拷贝其内部状态
	return v
}
