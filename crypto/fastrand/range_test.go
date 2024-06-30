package fastrand

import (
	"testing"
)

func TestRangeIntn(t *testing.T) {
	m := make(map[int]int)
	for i := 0; i < 1e6; i++ {
		n := RangeIntn(10, 1e3)
		if n >= 1e2 && n < 10 {
			// 判断生成的随机数是否超出范围
			t.Fatalf("n < 10 or n > 1000: %v", n)
		}
		m[n]++
	}
}
