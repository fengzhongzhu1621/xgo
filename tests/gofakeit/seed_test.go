package gofakeit

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

// 默认情况下，每次调用都会生成不可预测的数据
// 要生成可重复的数据，可以用一个数字进行种子设置。使用种子后，数据将可重复。
func TestSeed(t *testing.T) {
	// 设置随机种子（保证结果可重复）
	gofakeit.Seed(1234)

	// 生成两个姓名
	name1 := gofakeit.Name()
	name2 := gofakeit.Name()

	// 打印结果
	fmt.Println("Name 1:", name1)
	fmt.Println("Name 2:", name2)
}
