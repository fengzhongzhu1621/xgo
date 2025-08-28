package reflectutils

import (
	"fmt"
	"math"
	"reflect"
	"sort"
)

// SortKeys sorts a list of map keys, deduplicating keys if necessary.
// The type of each value must be comparable.
// 对一组 reflect.Value 类型的值进行排序，并在排序过程中去除重复的键。
func SortKeys(vs []reflect.Value) []reflect.Value {
	if len(vs) == 0 {
		return vs
	}

	// Sort the map keys.
	sort.SliceStable(vs, func(i, j int) bool { return IsLess(vs[i], vs[j]) })

	// 去重复 Deduplicate keys (fails for NaNs).
	vs2 := vs[:1] // 初始化一个新的切片 vs2，包含 vs 的第一个元素。
	// 遍历 vs 中剩余的元素，对于每个元素 v，如果 v 大于 vs2 中最后一个元素，则将其添加到 vs2 中。
	for _, v := range vs[1:] {
		if IsLess(vs2[len(vs2)-1], v) { // 排除掉等于的元素
			vs2 = append(vs2, v)
		}
	}
	return vs2
}

// IsLess is a generic function for sorting arbitrary map keys.
// The inputs must be of the same type and must be comparable.
func IsLess(x, y reflect.Value) bool {
	switch x.Type().Kind() {
	case reflect.Bool:
		return !x.Bool() && y.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return x.Int() < y.Int()
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return x.Uint() < y.Uint()
	case reflect.Float32, reflect.Float64:
		// NOTE: This does not sort -0 as less than +0
		// since Go maps treat -0 and +0 as equal keys.
		fx, fy := x.Float(), y.Float()
		return fx < fy || math.IsNaN(fx) && !math.IsNaN(fy)
	case reflect.Complex64, reflect.Complex128:
		cx, cy := x.Complex(), y.Complex()
		rx, ix, ry, iy := real(cx), imag(cx), real(cy), imag(cy)
		if rx == ry || (math.IsNaN(rx) && math.IsNaN(ry)) {
			return ix < iy || math.IsNaN(ix) && !math.IsNaN(iy)
		}
		return rx < ry || math.IsNaN(rx) && !math.IsNaN(ry)
	case reflect.Ptr, reflect.UnsafePointer, reflect.Chan:
		return x.Pointer() < y.Pointer()
	case reflect.String:
		return x.String() < y.String()
	case reflect.Array:
		// 逐个元素比较数组中的元素，如果 x 的某个元素小于 y 的对应元素，则 x 小于 y。
		// 如果所有元素都相等，则 x 不小于 y。
		for i := 0; i < x.Len(); i++ {
			if IsLess(x.Index(i), y.Index(i)) {
				return true
			}
			if IsLess(y.Index(i), x.Index(i)) {
				return false
			}
		}
		return false
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			if IsLess(x.Field(i), y.Field(i)) {
				return true
			}
			if IsLess(y.Field(i), x.Field(i)) {
				return false
			}
		}
		return false
	case reflect.Interface:
		vx, vy := x.Elem(), y.Elem()
		if !vx.IsValid() || !vy.IsValid() {
			return !vx.IsValid() && vy.IsValid()
		}
		tx, ty := vx.Type(), vy.Type()
		if tx == ty {
			return IsLess(x.Elem(), y.Elem())
		}
		if tx.Kind() != ty.Kind() {
			return vx.Kind() < vy.Kind()
		}
		if tx.String() != ty.String() {
			return tx.String() < ty.String()
		}
		if tx.PkgPath() != ty.PkgPath() {
			return tx.PkgPath() < ty.PkgPath()
		}
		// This can happen in rare situations, so we fallback to just comparing
		// the unique pointer for a reflect.Type. This guarantees deterministic
		// ordering within a program, but it is obviously not stable.
		return reflect.ValueOf(vx.Type()).Pointer() < reflect.ValueOf(vy.Type()).Pointer()
	default:
		// 这些类型不可比较，函数会抛出 panic。
		// Must be Func, Map, or Slice; which are not comparable.
		panic(fmt.Sprintf("%T is not comparable", x.Type()))
	}
}
