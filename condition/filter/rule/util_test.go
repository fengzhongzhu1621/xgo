package rule

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	jsonx "github.com/fengzhongzhu1621/xgo/crypto/encoding/json"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/stretchr/testify/assert"
)

//	{
//		"condition": "AND",
//		"rules": [
//				{
//						"field": "test",
//						"operator": "equal",
//						"value": 1
//				},
//				{
//						"condition": "OR",
//						"rules": [
//								{
//										"field": "test1",
//										"operator": "filter_array",
//										"value": {
//												"field": "element",
//												"operator": "filter_object",
//												"value": {
//														"condition": "AND",
//														"rules": [
//																{
//																		"field": "test2",
//																		"operator": "in",
//																		"value": [
//																				"b",
//																				"c"
//																		]
//																}
//														]
//												}
//										}
//								},
//								{
//										"field": "test3",
//										"operator": "datetime_less",
//										"value": 1
//								}
//						]
//				}
//		]
//	}
var (
	exampleRule = &CombinedRule{
		Condition: operator.And,
		Rules: []operator.IRuleFactory{
			&AtomRule{
				Field:    "test",
				Operator: operator.Equal,
				Value:    1,
			},
			&CombinedRule{
				Condition: operator.Or,
				Rules: []operator.IRuleFactory{
					&AtomRule{
						Field:    "test1",
						Operator: operator.Array,
						Value: &AtomRule{
							Field:    operator.ArrayElement,
							Operator: operator.Object,
							Value: &CombinedRule{
								Condition: operator.And,
								Rules: []operator.IRuleFactory{
									&AtomRule{
										Field:    "test2",
										Operator: operator.In,
										Value:    []string{"b", "c"},
									},
								},
							},
						},
					},
					&AtomRule{
						Field:    "test3",
						Operator: operator.DatetimeLess,
						Value:    1,
					},
				},
			},
		},
	}
)

func TestRuleMatch(t *testing.T) {
	matched, err := exampleRule.Match(
		jsonx.JsonString(`{"test":1,"test1":[{"test2":"b"}],"test3":111}`),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = exampleRule.Match(
		jsonx.JsonString(`{"test":1,"test1":[{"test2":"a"}],"test3":111}`),
	)
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}

func testExampleRule(t *testing.T, r operator.IRuleFactory) {
	if r == nil {
		t.Errorf("rule is nil")
		return
	}

	rule, ok := r.(*CombinedRule)
	if !ok {
		t.Errorf("rule %+v is not combined type", r)
		return
	}

	if rule.Condition != operator.And {
		t.Errorf("rule condition %s is not and", rule.Condition)
		return
	}

	if len(rule.Rules) != 2 {
		t.Errorf("rules length %d is not 2", len(rule.Rules))
		return
	}

	subAtomRule, ok := rule.Rules[0].(*AtomRule)
	if !ok {
		t.Errorf("first sub rule %+v is not atom type", rule.Rules[0])
		return
	}

	if subAtomRule.Field != "test" {
		t.Errorf("first sub rule field %s is not test", subAtomRule.Field)
		return
	}

	if subAtomRule.Operator != operator.Equal {
		t.Errorf("first sub rule op %s is not equal", subAtomRule.Operator)
		return
	}

	intVal, err := cast.ToInt64E(subAtomRule.Value)
	if err != nil {
		t.Errorf("first sub rule value %v is invalid", subAtomRule.Value)
		return
	}

	if intVal != 1 {
		t.Errorf("first sub rule value %v is not 1", subAtomRule.Value)
		return
	}

	subCombinedRule, ok := rule.Rules[1].(*CombinedRule)
	if !ok {
		t.Errorf("second sub rule %+v is not combined type", rule.Rules[1])
		return
	}

	if subCombinedRule.Condition != operator.Or {
		t.Errorf("second sub rule condition %s is not combined type", rule.Condition)
		return
	}

	if len(subCombinedRule.Rules) != 2 {
		t.Errorf("second sub rules length %d is not 2", len(rule.Rules))
		return
	}

	subAtomRule1, ok := subCombinedRule.Rules[0].(*AtomRule)
	if !ok {
		t.Errorf("first sub sub rule %+v is not atom type", subCombinedRule.Rules[0])
		return
	}

	if subAtomRule1.Field != "test1" {
		t.Errorf("first sub sub rule field %s is not test1", subAtomRule1.Field)
		return
	}

	if subAtomRule1.Operator != operator.Array {
		t.Errorf("first sub sub rule op %s is not ne", subAtomRule1.Operator)
		return
	}

	filterArrVal, ok := subAtomRule1.Value.(*AtomRule)
	if !ok {
		t.Errorf("first sub sub rule value %v is invalid", subAtomRule1.Value)
		return
	}

	if filterArrVal.Field != operator.ArrayElement {
		t.Errorf("filter array rule field %s is not %s", subAtomRule.Field, operator.ArrayElement)
		return
	}

	if filterArrVal.Operator != operator.Object {
		t.Errorf("filter array rule op %s is not filter object", subAtomRule.Operator)
		return
	}

	filterArrValRule, ok := filterArrVal.Value.(*CombinedRule)
	if !ok {
		t.Errorf("filter array rule value %v is invalid", filterArrVal.Value)
		return
	}

	if filterArrValRule.Condition != operator.And {
		t.Errorf("filter array sub condition %s is not and", filterArrValRule.Condition)
		return
	}

	if len(filterArrValRule.Rules) != 1 {
		t.Errorf("filter array sub rules length %d is not 1", len(rule.Rules))
		return
	}

	filterObjRule, ok := filterArrValRule.Rules[0].(*AtomRule)
	if !ok {
		t.Errorf("filter object rule %+v is not atom type", filterArrValRule.Rules[0])
		return
	}

	if filterObjRule.Field != "test2" {
		t.Errorf("filter object rule field %s is not test2", filterObjRule.Field)
		return
	}

	if filterObjRule.Operator != operator.In {
		t.Errorf("filter object rule op %s is not in", filterObjRule.Operator)
		return
	}

	arrVal, ok := filterObjRule.Value.([]interface{})
	if !ok {
		t.Errorf("filter object rule value %v is invalid", filterObjRule.Value)
		return
	}

	if len(arrVal) != 2 {
		t.Errorf("array value length %d is not 2", len(arrVal))
		return
	}

	strVal1, ok := arrVal[0].(string)
	if !ok {
		t.Errorf("first array value %v is invalid", arrVal[0])
		return
	}

	if strVal1 != "b" {
		t.Errorf("first array value %v is not b", arrVal[0])
		return
	}

	strVal2, ok := arrVal[1].(string)
	if !ok {
		t.Errorf("second array value %v is invalid", arrVal[1])
		return
	}

	if strVal2 != "c" {
		t.Errorf("second array value %v is not c", arrVal[1])
		return
	}

	subAtomRule2, ok := subCombinedRule.Rules[1].(*AtomRule)
	if !ok {
		t.Errorf("second sub sub rule %+v is not atom type", subCombinedRule.Rules[1])
		return
	}

	if subAtomRule2.Field != "test3" {
		t.Errorf("second sub sub rule field %s is not test3", subAtomRule2.Field)
		return
	}

	if subAtomRule2.Operator != operator.DatetimeLess {
		t.Errorf("second sub sub rule op %s is not datetime less", subAtomRule2.Operator)
		return
	}

	if !validator.IsNumeric(subAtomRule2.Value) {
		t.Errorf("second sub rule value %v is invalid", subAtomRule2.Value)
		return
	}
}
