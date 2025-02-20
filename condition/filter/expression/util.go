package expression

import (
	"fmt"
	"reflect"

	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/fengzhongzhu1621/xgo/condition/filter/rule"
)

// GenAtomFilter generate atom rule filter
func GenAtomFilter(field string, op operator.OpType, value interface{}) *Expression {
	return &Expression{
		IRuleFactory: &rule.AtomRule{
			Field:    field,
			Operator: op,
			Value:    value,
		},
	}
}

// And merge expressions using 'and' operation.
func And(rules ...operator.IRuleFactory) (*Expression, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("rules are not set")
	}

	andRules := make([]operator.IRuleFactory, 0)
	for _, rule_obj := range rules {
		if rule_obj == nil || reflect.ValueOf(rule_obj).IsNil() {
			continue
		}

		for expr, ok := rule_obj.(*Expression); ok; expr, ok = rule_obj.(*Expression) {
			rule_obj = expr.IRuleFactory
		}

		switch rule_obj.WithType() {
		case rule.AtomType:
			andRules = append(andRules, rule_obj)
		case rule.CombinedType:
			combinedRule, ok := rule_obj.(*rule.CombinedRule)
			if !ok {
				return nil, fmt.Errorf("combined rule type is invalid")
			}
			if combinedRule.Condition == operator.And {
				andRules = append(andRules, combinedRule.Rules...)
				continue
			}
			andRules = append(andRules, combinedRule)
		default:
			return nil, fmt.Errorf("rule type %s is invalid", rule_obj.WithType())
		}
	}

	if len(andRules) == 0 {
		return nil, fmt.Errorf("rules are all nil")
	}

	if len(andRules) == 1 {
		return &Expression{
			IRuleFactory: andRules[0],
		}, nil
	}

	return &Expression{
		IRuleFactory: &rule.CombinedRule{
			Condition: operator.And,
			Rules:     andRules,
		},
	}, nil
}
