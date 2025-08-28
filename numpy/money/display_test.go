package money

import (
	"fmt"
	"testing"

	"github.com/Rhymond/go-money"
)

func TestDisplay(t *testing.T) {
	// 使用 Display() 格式化
	formatted := money.New(123456789, money.CNY).Display()
	fmt.Println(formatted) // 1,234,567.89 元

	// 使用 AsMajorUnits() 格式化为浮点数表示的金额值
	majorUnits := money.New(123456789, money.CNY).AsMajorUnits()
	fmt.Println(majorUnits) // 1.23456789e+06
}
