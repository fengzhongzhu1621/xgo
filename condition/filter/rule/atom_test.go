package rule

import (
	"reflect"
	"testing"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
)

func TestAtomRuleRuleWithType(t *testing.T) {
	var rule operator.IRuleFactory = new(AtomRule)
	if rule.WithType() != AtomType {
		t.Errorf("rule type %s is invalid", rule.WithType())
		return
	}
}

func TestAtomRuleValidate(t *testing.T) {
	var (
		rule operator.IRuleFactory
		opt  *operator.ExprOption
	)

	// 验证 NotIn 操作符原子规则
	opt = operator.NewDefaultExprOpt(map[string]criteria.FieldType{"test1": criteria.String})
	rule = &AtomRule{
		Field:    "test1",
		Operator: operator.NotIn,
		Value:    []string{"a", "b", "c"},
	}
	if err := rule.Validate(opt); err != nil {
		t.Errorf("rule validate failed, err: %v", err)
		return
	}

	// 限制 NotIn 操作符匹配的值的数量
	opt = &operator.ExprOption{
		RuleFields: map[string]criteria.FieldType{
			"test1": criteria.String,
		},
		MaxNotInLimit: 3,
	}
	if err := rule.Validate(opt); err != nil {
		t.Errorf("rule validate failed, err: %v", err)
		return
	}

	// 验证 Field 的类型配置错误
	opt = &operator.ExprOption{
		RuleFields: map[string]criteria.FieldType{
			"test2": criteria.String,
		},
	}
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}

	// 验证没有配置 Field 的类型
	opt = &operator.ExprOption{
		MaxRulesLimit: 10,
		MaxNotInLimit: 2,
	}
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
}

func TestAtomRuleFields(t *testing.T) {
	var rule operator.IRuleFactory = &AtomRule{
		Field:    "test1",
		Operator: operator.Equal,
		Value:    1,
	}

	fields := rule.RuleFields()
	if !reflect.DeepEqual(fields, []string{"test1"}) {
		t.Errorf("rule fields %+v is invalid", fields)
		return
	}
}

func TestAtomRuleToMgo(t *testing.T) {
	var rule operator.IRuleFactory = &AtomRule{
		Field:    "test1",
		Operator: operator.NotIn,
		Value:    []string{"a", "b", "c"},
	}

	mgo, err := rule.ToMgo(nil)
	if err != nil {
		t.Errorf("covert rule to mongo failed, err: %v", err)
		return
	}

	expectMgo := map[string]interface{}{
		"test1": map[string]interface{}{
			operator.DBNIN: []string{"a", "b", "c"},
		},
	}

	if !reflect.DeepEqual(mgo, expectMgo) {
		t.Errorf("rule mongo condition %+v is invalid", mgo)
		return
	}
}
