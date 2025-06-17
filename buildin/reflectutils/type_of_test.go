package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypeOf(t *testing.T) {
	// reflect.TypeOf(map[string]interface{}(nil))返回一个表示map[string]interface{}类型的reflect.Type值。
	tp := reflect.TypeOf(map[string]interface{}(nil))
	fmt.Println(tp) // map[string]interface {}
}

func TestTypeAndKind(t *testing.T) {
	var a int
	var s string
	var p Person
	var sl []int
	var m map[string]int

	// 获取 reflect.Type
	ta := reflect.TypeOf(a)
	ts := reflect.TypeOf(s)
	tp := reflect.TypeOf(p)
	tsl := reflect.TypeOf(sl)
	tm := reflect.TypeOf(m)

	// 获取 reflect.Kind
	ka := ta.Kind()
	ks := ts.Kind()
	kp := tp.Kind()
	ksl := tsl.Kind()
	km := tm.Kind()

	fmt.Println("Type & Kind 示例:")
	fmt.Printf("a: Type = %v, Kind = %v\n", ta, ka)    // a: Type = int, Kind = int
	fmt.Printf("s: Type = %v, Kind = %v\n", ts, ks)    // s: Type = string, Kind = string
	fmt.Printf("p: Type = %v, Kind = %v\n", tp, kp)    // p: Type = main.Person, Kind = struct
	fmt.Printf("sl: Type = %v, Kind = %v\n", tsl, ksl) // sl: Type = []int, Kind = slice
	fmt.Printf("m: Type = %v, Kind = %v\n", tm, km)    // m: Type = map[string]int, Kind = map
}
