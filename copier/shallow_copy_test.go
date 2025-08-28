package copier

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Config struct {
	Timeout int
}

func TestShallowCopy(t *testing.T) {
	src := Config{Timeout: 10}
	var dst Config

	err := ShallowCopy(&dst, &src)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("dst =", dst) // 输出: dst = {10}
	}
}

func TestShallowCopy2(t *testing.T) {
	type Data struct {
		Values []int
	}

	src := Data{Values: []int{1, 2, 3}}
	var dst Data

	err := ShallowCopy(&dst, &src)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("dst =", dst) // 输出: dst = {Values:[1 2 3]}
		dst.Values[0] = 100
		// dst 修改后，src也会同步变更
		fmt.Println("src =", src) // 输出: src = {Values:[100 2 3]} （共享底层数组）
	}
}

type shallowCopyTestData struct {
	I             int
	unexportedInt int
	IntPtr        *int
}

func TestShallowCopy3(t *testing.T) {
	require.NotNil(t, ShallowCopy(1, 2), "dst must be a pointer")

	var i int
	require.NotNil(t, ShallowCopy(&i, 1.0), "type miss match")
	require.Nil(t, ShallowCopy(&i, 2))
	require.Equal(t, 2, i)

	s := shallowCopyTestData{I: 1, unexportedInt: 2, IntPtr: &i}
	var ss shallowCopyTestData
	require.Nil(t, ShallowCopy(&ss, &s), "src may or may not be a pointer")
	require.Equal(t, s, ss)
	require.Equal(t, s.unexportedInt, ss.unexportedInt)
	require.Equal(t,
		reflect.ValueOf(s.IntPtr).Pointer(),
		reflect.ValueOf(ss.IntPtr).Pointer())
}
