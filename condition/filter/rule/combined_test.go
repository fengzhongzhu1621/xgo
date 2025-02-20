package rule

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
)

func TestCombinedRuleRuleWithType(t *testing.T) {
	var rule operator.IRuleFactory

	rule = new(CombinedRule)
	if rule.WithType() != CombinedType {
		t.Errorf("rule type %s is invalid", rule.WithType())
		return
	}
}
