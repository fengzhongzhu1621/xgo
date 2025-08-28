package money

import (
	"fmt"
	"log"
	"testing"

	"github.com/Rhymond/go-money"
)

func TestOperator(t *testing.T) {
	// 初始化金额
	yuan := money.New(100, money.CNY)          // 1 人民币
	twoYuan := money.New(200, money.CNY)       // 2 人民币
	negativeYuan := money.New(-100, money.CNY) // -1 人民币

	// 加法
	result, err := yuan.Add(twoYuan)
	if err != nil {
		log.Fatalf("加法错误: %v", err)
	}
	fmt.Printf("1 人民币 + 2 人民币 = %s\n", result.Display()) // 输出: ￥3.00

	// 减法
	result, err = yuan.Subtract(twoYuan)
	if err != nil {
		log.Fatalf("减法错误: %v", err)
	}
	fmt.Printf("1 人民币 - 2 人民币 = %s\n", result.Display()) // 输出: -￥1.00

	// 乘法
	result = yuan.Multiply(2)
	fmt.Printf("1 人民币 * 2 = %s\n", result.Display()) // 输出: ￥2.00

	// 绝对值
	result = negativeYuan.Absolute()
	fmt.Printf("-1 人民币 的绝对值 = %s\n", result.Display()) // 输出: ￥1.00

	// 负值
	result = negativeYuan.Negative()
	fmt.Printf("-1 人民币 的负值 = %s\n", result.Display()) // 输出: ￥1.00

	// 混合货币类型示例（将引发错误）
	usd := money.New(150, money.USD)
	_, err = yuan.Add(usd)
	if err != nil {
		fmt.Printf("尝试将 CNY 和 USD 相加时出错: %v\n", err)
	}
}
