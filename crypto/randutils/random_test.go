package randutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/random"
)

// 生成随机int, 范围[min, max)
func TestRandomRandInt(t *testing.T) {
	rInt := random.RandInt(1, 10)
	fmt.Println(rInt)
}

func TestRandomRandString(t *testing.T) {
	randStr := random.RandString(6)
	fmt.Println(randStr)
}
