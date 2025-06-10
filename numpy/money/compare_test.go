package money

import (
	"testing"

	"github.com/Rhymond/go-money"
)

func TestCompareMoney(t *testing.T) {
	pound := money.New(100, money.CNY)     // 1 人名币
	twoPounds := money.New(200, money.CNY) // 2 人名币
	twoEuros := money.New(200, money.CNY)  // 2 人名币

	// 比较人名币数额
	pound.GreaterThan(twoPounds) // 返回 false, nil
	// 解释：1 人名币 < 2 人名币

	pound.LessThan(twoPounds) // 返回 true, nil
	// 解释：1 人名币 < 2 人名币

	twoPounds.Equals(twoEuros) // 返回 true, nil
	// 解释：2 人名币 == 2 人名币

	twoPounds.Compare(pound) // 返回 1, nil
	// 解释：2 人名币 > 1 人名币

	pound.Compare(twoPounds) // 返回 -1, nil
	// 解释：1 人名币 < 2 人名币

	pound.Compare(pound) // 返回 0, nil
	// 解释：1 人名币 == 1 人名币

	pound.Compare(twoEuros) // 返回 -1, nil
	// 解释：1 人名币 < 2 人名币
}
