package intern

import (
	"fmt"
	"testing"

	"github.com/josharian/intern"
)

func TestIntern(t *testing.T) {
	a := "hello"
	b := []byte("hello")

	// 使用 String() 驻留字符串
	s1 := intern.String(a)
	s2 := intern.String(a)
	fmt.Println(s1 == s2)

	// 使用 Bytes() 驻留从 []byte 转换的字符串
	s3 := intern.Bytes(b)
	s4 := intern.Bytes(b)
	fmt.Println(s3 == s4)
}
