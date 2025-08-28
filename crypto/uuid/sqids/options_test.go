package sqids

import (
	"fmt"
	"testing"

	"github.com/sqids/sqids-go"
)

func TestSqidsOption(t *testing.T) {
	// 1. 创建 Sqids 实例（自定义 Alphabet）
	s, err := sqids.New(sqids.Options{
		Alphabet: "FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE",
	})
	if err != nil {
		fmt.Printf("创建 Sqids 实例失败: %v\n", err)
		return
	}

	// 2. 编码数字数组
	id, err := s.Encode([]uint64{1, 2, 3})
	if err != nil {
		fmt.Printf("编码失败: %v\n", err)
		return
	}
	fmt.Println("编码结果:", id) // B4aajs

	// 3. 解码短 ID
	numbers := s.Decode(id)
	fmt.Println("解码结果:", numbers) // [1 2 3]
}
