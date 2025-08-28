package expression

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
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

func TestExpressionValidateOption(t *testing.T) {
	expr := &Expression{
		IRuleFactory: &rule.CombinedRule{
			Condition: operator.And,
			Rules: []operator.IRuleFactory{
				&rule.AtomRule{
					Field:    "string",
					Operator: operator.Equal,
					Value:    "a",
				},
				&rule.CombinedRule{
					Condition: operator.Or,
					Rules: []operator.IRuleFactory{
						&rule.AtomRule{
							Field:    "int",
							Operator: operator.Greater,
							Value:    123,
						},
						&rule.AtomRule{
							Field:    "enum_array",
							Operator: operator.In,
							Value:    []string{"b", "c"},
						},
						&rule.AtomRule{
							Field:    "int_array",
							Operator: operator.NotIn,
							Value:    []int64{1, 3, 5},
						},
						&rule.AtomRule{
							Field:    "bool",
							Operator: operator.NotEqual,
							Value:    false,
						},
					},
				},
				&rule.AtomRule{
					Field:    "time",
					Operator: operator.DatetimeLessOrEqual,
					Value:    time.Now().Unix(),
				},
				&rule.AtomRule{
					Field:    "time",
					Operator: operator.DatetimeGreater,
					Value:    "2006-01-02 15:04:05",
				},
				// TODO confirm how to deal with filter object & array
			},
		},
	}

	opt := operator.NewDefaultExprOpt(map[string]criteria.FieldType{
		"string":     criteria.String,
		"int":        criteria.Numeric,
		"enum_array": criteria.Enum,
		"int_array":  criteria.Numeric,
		"bool":       criteria.Boolean,
		"time":       criteria.Time,
	})

	if err := expr.Validate(opt); err != nil {
		t.Errorf("validate expression failed, err: %v", err)
		return
	}

	// test invalidate scenario
	opt.RuleFields["string"] = criteria.Numeric
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a numeric") {
		t.Errorf("validate numeric type failed, err: %v", err)
		return
	}
	opt.RuleFields["string"] = criteria.String

	opt.RuleFields["int"] = criteria.String
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a string") {
		t.Errorf("validate string type failed, err: %v", err)
		return
	}
	opt.RuleFields["int"] = criteria.Numeric

	opt.RuleFields["enum_array"] = criteria.Boolean
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a boolean") {
		t.Errorf("validate bool type failed, err: %v", err)
		return
	}
	opt.RuleFields["enum_array"] = criteria.String

	opt.RuleFields["bool"] = criteria.Time
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "is not of time type") {
		t.Errorf("validate time type failed, err: %v", err)
		return
	}
	opt.RuleFields["bool"] = criteria.Boolean

	opt.RuleFields["time"] = criteria.Boolean
	if err := expr.Validate(opt); !strings.Contains(err.Error(), "value should be a boolean") {
		t.Errorf("validate boolean type failed, err: %v", err)
		return
	}
	opt.RuleFields["time"] = criteria.Time

	opt.MaxRulesDepth = 2
	if err := expr.Validate(opt); !strings.Contains(
		err.Error(),
		"expression rules depth exceeds maximum",
	) {
		t.Errorf("validate rule depth failed, err: %v", err)
		return
	}
	opt.MaxRulesDepth = 3

	opt.MaxRulesLimit = 3
	if err := expr.Validate(opt); !strings.Contains(
		err.Error(),
		"rules elements number exceeds limit: 3",
	) {
		t.Errorf("validate rule limit failed, err: %v", err)
		return
	}
	opt.MaxRulesLimit = 4

	opt.MaxInLimit = 1
	if err := expr.Validate(opt); !strings.Contains(
		err.Error(),
		"elements length 2 exceeds maximum 1",
	) {
		t.Errorf("validate rule in limit failed, err: %v", err)
		return
	}
	opt.MaxInLimit = 2

	opt.MaxNotInLimit = 2
	if err := expr.Validate(opt); !strings.Contains(
		err.Error(),
		"elements length 3 exceeds maximum 2",
	) {
		t.Errorf("validate rule in limit failed, err: %v", err)
		return
	}
	opt.MaxNotInLimit = 3
}
