package slice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice2(t *testing.T) {
	a := make([]int, 10, 10)
	fmt.Println("a = ", a)

	b := a
	b[0] = 10

	// slice的底层结构其中一个实际上是有一个指针，指向了一个数组。那么，
	// 在把a赋值给b的时候，只是把slice的结构也就是Array、Len和Cap复制给了b，但Array指向的数组还是同一个。
	// 所以，这就是为什么更改了b[0]，a[0]的值也更改了的原因。
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)
}

func TestCreateBigSlice(t *testing.T) {
	size := 1 << 30 // 1G
	s := make([]byte, size)
	assert.Equal(t, 1024*1024*1024, len(s))
	assert.Equal(t, 1024*1024*1024, cap(s))
}
