package rule

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCombinedRuleRuleWithType(t *testing.T) {
	var rule operator.IRuleFactory = new(CombinedRule)
	if rule.WithType() != CombinedType {
		t.Errorf("rule type %s is invalid", rule.WithType())
		return
	}
}

func TestJsonUnmarshalRule(t *testing.T) {
	ruleJson, err := json.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	rule := new(CombinedRule)
	err = json.Unmarshal(ruleJson, rule)
	if err != nil {
		t.Error(err)
		return
	}

	testExampleRule(t, rule)
}

func TestBsonUnmarshalRule(t *testing.T) {
	ruleBson, err := bson.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	rule := new(CombinedRule)
	err = bson.Unmarshal(ruleBson, rule)
	if err != nil {
		t.Error(err)
		return
	}

	testExampleRule(t, rule)
}

func TestRuleValidate(t *testing.T) {
	// test combined rule validation
	rule := exampleRule

	// TODO confirm how to deal with object & array
	opt := &operator.ExprOption{
		RuleFields: map[string]criteria.FieldType{
			"test":                criteria.Numeric,
			"test1":               criteria.Array,
			"test1.element":       criteria.Object,
			"test1.element.test2": criteria.String,
			"test3":               criteria.Time,
		},
		MaxInLimit:    2,
		MaxRulesLimit: 2,
		MaxRulesDepth: 6,
	}

	if err := rule.Validate(opt); err != nil {
		t.Errorf("rule validate failed, err: %v", err)
		return
	}

	// test invalidate scenario
	opt.RuleFields["test"] = criteria.String
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.RuleFields["test"] = criteria.Numeric

	delete(opt.RuleFields, "test")
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.RuleFields["test"] = criteria.Numeric

	opt.MaxInLimit = 1
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.MaxInLimit = 0

	opt.MaxRulesLimit = 1
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.MaxRulesLimit = 0

	opt.MaxRulesDepth = 5
	if err := rule.Validate(opt); err == nil {
		t.Error("rule validate failed")
		return
	}
	opt.MaxRulesDepth = 0
}

func TestRuleFields(t *testing.T) {
	// test combined rule to mongo
	rule := exampleRule
	mgo, err := rule.ToMgo(nil)
	if err != nil {
		t.Errorf("covert rule to mongo failed, err: %v", err)
		return
	}

	expectMgo := map[string]interface{}{
		operator.DBAND: []map[string]interface{}{{
			"test": map[string]interface{}{operator.DBEQ: 1},
		}, {
			operator.DBOR: []map[string]interface{}{{
				operator.DBAND: []map[string]interface{}{{
					"test1.test2": map[string]interface{}{operator.DBIN: []string{"b", "c"}},
				}},
			}, {
				"test3": map[string]interface{}{operator.DBLT: time.Unix(1, 0)},
			}},
		}},
	}

	if !reflect.DeepEqual(mgo, expectMgo) {
		t.Errorf("rule mongo condition %+v is invalid", mgo)
		return
	}

	// test invalid combined rule to mongo scenario
	rule = &CombinedRule{
		Condition: "test",
		Rules:     []operator.IRuleFactory{exampleRule},
	}

	if _, err = rule.ToMgo(nil); err == nil {
		t.Errorf("covert rule to mongo should fail")
		return
	}

	rule = &CombinedRule{
		Condition: "",
		Rules:     []operator.IRuleFactory{exampleRule},
	}

	if _, err = rule.ToMgo(nil); err == nil {
		t.Errorf("covert rule to mongo should fail")
		return
	}

	rule = &CombinedRule{
		Condition: "test",
		Rules: []operator.IRuleFactory{
			&AtomRule{
				Field:    "test1",
				Operator: operator.In,
				Value:    []interface{}{"a", 1, "c"},
			},
		},
	}

	if _, err = rule.ToMgo(nil); err == nil {
		t.Errorf("covert rule to mongo should fail")
		return
	}
}
