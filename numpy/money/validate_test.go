package money

import (
	"testing"

	"github.com/Rhymond/go-money"
)

func TestValidate(t *testing.T) {
	pound := money.New(100, money.CNY)

	// 断定是否为零
	pound.IsZero() // 返回 false

	// 断定是否为正值
	pound.IsPositive() // 返回 true

	// 断定是否为负值
	pound.IsNegative() // 返回 false
}
