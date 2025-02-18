package rule

import (
	"reflect"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/stretchr/testify/assert"
)

func TestObjectValidate(t *testing.T) {
	op := operator.GetOperator(operator.Object)

	opt := &operator.ExprOption{
		IgnoreRuleFields: true,
		MaxInLimit:       2,
		MaxRulesLimit:    2,
		MaxRulesDepth:    7,
	}

	// test filter object normal type
	err := op.ValidateValue(exampleRule, opt)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test filter mapstr type
	err = op.ValidateValue(&AtomRule{
		Field:    "test",
		Operator: operator.Object,
		Value: &AtomRule{
			Field:    "a",
			Operator: operator.NotEqual,
			Value:    4,
		},
	}, operator.NewDefaultExprOpt(map[string]criteria.FieldType{"test": criteria.MapString}))
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid filter object type
	err = op.ValidateValue(1, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue("a", opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(false, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(map[string]interface{}{"test1": 1}, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(struct{}{}, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(nil, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]int64{1}, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]interface{}{1, "a"}, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(CombinedRule{
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
						Operator: operator.In,
						Value:    "a",
					},
				},
			},
		},
	}, opt)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}
}

func TestObjectMongoCond(t *testing.T) {
	op := operator.GetOperator(operator.Object)

	// test filter object normal type
	cond, err := op.ToMgo("obj", exampleRule)
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	expectCond := map[string]interface{}{
		operator.DBAND: []map[string]interface{}{
			{
				"obj.test": map[string]interface{}{operator.DBEQ: 1},
			}, {
				operator.DBOR: []map[string]interface{}{
					{
						operator.DBAND: []map[string]interface{}{{
							"obj.test1.test2": map[string]interface{}{operator.DBIN: []string{"b", "c"}},
						}},
					}, {
						"obj.test3": map[string]interface{}{operator.DBLT: time.Unix(1, 0)},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(cond, expectCond) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestObjectMatch(t *testing.T) {
	op := operator.GetOperator(operator.Object)

	// test filter object normal type
	matched, err := op.Match(`{"test":1,"test1":[{"test2":"d"}],"test3":0}`, exampleRule)
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}
