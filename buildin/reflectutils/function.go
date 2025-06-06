package reflectutils

import (
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type FuncType int

const (
	_ FuncType = iota

	tbFunc  // func(T) bool
	ttbFunc // func(T, T) bool
	trbFunc // func(T, R) bool
	tibFunc // func(T, I) bool
	trFunc  // func(T) R

	Equal             = ttbFunc // func(T, T) bool
	EqualAssignable   = tibFunc // func(T, I) bool; encapsulates func(T, T) bool
	Transformer       = trFunc  // func(T) R
	ValueFilter       = ttbFunc // func(T, T) bool
	Less              = ttbFunc // func(T, T) bool
	ValuePredicate    = tbFunc  // func(T) bool
	KeyValuePredicate = trbFunc // func(T, R) bool
)

var BoolType = reflect.TypeOf(true)

// IsType reports whether the reflect.Type is of the specified function type.
func IsType(t reflect.Type, ft FuncType) bool {
	// IsVariadic: 判断函数是否接受可变参数
	if t == nil || t.Kind() != reflect.Func || t.IsVariadic() {
		return false
	}
	// 获得函数输入和输出参数的数量
	ni, no := t.NumIn(), t.NumOut()

	switch ft {
	case tbFunc: // func(T) bool
		if ni == 1 && no == 1 && t.Out(0) == BoolType {
			return true
		}
	case ttbFunc: // func(T, T) bool
		if ni == 2 && no == 1 && t.In(0) == t.In(1) && t.Out(0) == BoolType {
			return true
		}
	case trbFunc: // func(T, R) bool
		if ni == 2 && no == 1 && t.Out(0) == BoolType {
			return true
		}
	case tibFunc: // func(T, I) bool
		// 判断 t 类型的值可否赋值给 u 类型。 func (t *rtype) AssignableTo(u reflect.Type) bool
		if ni == 2 && no == 1 && t.In(0).AssignableTo(t.In(1)) && t.Out(0) == BoolType {
			return true
		}
	case trFunc: // func(T) R
		if ni == 1 && no == 1 {
			return true
		}
	}
	return false
}

// 在函数名中提取最后一个标识符部分，例如从 "mypkg.myMethod" 中提取 "myMethod"。
var lastIdentRx = regexp.MustCompile(`[_\p{L}][_\p{L}\p{N}]*$`)

// NameOf returns the name of the function value.
// 从 reflect.Value 提取函数或方法的名称
func NameOf(v reflect.Value) string {
	// 通过程序计数器（PC）获取对应的 *runtime.Func 对象，该对象包含了函数的名称和其他调试信息。
	fnc := runtime.FuncForPC(v.Pointer())
	if fnc == nil {
		return "<unknown>"
	}
	fullName := fnc.Name() // e.g., "long/path/name/mypkg.(*MyType).(long/path/name/mypkg.myMethod)-fm"

	// Method closures have a "-fm" suffix.
	fullName = strings.TrimSuffix(fullName, "-fm")

	var name string
	for len(fullName) > 0 {
		inParen := strings.HasSuffix(fullName, ")")
		fullName = strings.TrimSuffix(fullName, ")")

		s := lastIdentRx.FindString(fullName)
		if s == "" {
			break
		}
		name = s + "." + name
		fullName = strings.TrimSuffix(fullName, s)

		if i := strings.LastIndexByte(fullName, '('); inParen && i >= 0 {
			fullName = fullName[:i]
		}
		fullName = strings.TrimSuffix(fullName, ".")
	}
	return strings.TrimSuffix(name, ".")
}
