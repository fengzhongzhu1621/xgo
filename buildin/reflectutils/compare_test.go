package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

type T struct {
	f func() int
}

func TestCompareFunc(t *testing.T) {
	t1 := T{f: func() int { return 1 }}
	t2 := T{f: func() int { return 1 }}

	// 尽管t1和t2的函数逻辑相同，但reflect.DeepEqual比较的是函数的指针，因此返回false。
	fmt.Println(reflect.DeepEqual(t1, t2)) // 输出: false
}
