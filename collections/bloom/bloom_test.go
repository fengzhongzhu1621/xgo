package bloom

import (
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
)

func TestCreate(t *testing.T) {
	// bloom.New(size, hashCount) size是位数组大小，hashCount是哈希函数个数。
	filter := bloom.New(1000, 3)

	// 向布隆过滤器中添加数据
	filter.Add([]byte("a"))
	filter.Add([]byte("b"))

	// 测试数据是否可能在过滤器中
	fmt.Println(filter.Test([]byte("a"))) // true (可能存在误判)
	fmt.Println(filter.Test([]byte("b"))) // true (可能存在误判)
	fmt.Println(filter.Test([]byte("c"))) // false（一定正确）
}
