package rule

import (
	"reflect"
	"testing"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/stretchr/testify/assert"
)

func TestArrayValidate(t *testing.T) {
	op := operator.GetOperator(operator.Array)

	opt := operator.NewDefaultExprOpt(map[string]criteria.FieldType{"element": criteria.Numeric})
	opt.MaxRulesDepth = 4

	// test filter array normal type
	err := op.ValidateValue(&CombinedRule{
		Condition: operator.And,
		Rules: []operator.IRuleFactory{
			&AtomRule{
				Field:    operator.ArrayElement,
				Operator: operator.Equal,
				Value:    1,
			},
			&CombinedRule{
				Condition: operator.Or,
				Rules: []operator.IRuleFactory{
					&AtomRule{
						Field:    operator.ArrayElement,
						Operator: operator.NotEqual,
						Value:    2,
					},
					&AtomRule{
						Field:    operator.ArrayElement,
						Operator: operator.NotEqual,
						Value:    4,
					},
				},
			},
		},
	}, opt)
	if err != nil {
		t.Errorf("validate failed, err: %v", err)
		return
	}

	// test invalid filter array type
	err = op.ValidateValue(1, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue("a", nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(false, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(map[string]interface{}{"test1": 1}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(struct{}{}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(nil, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]int64{1}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue([]interface{}{1, "a"}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}

	err = op.ValidateValue(CombinedRule{
		Condition: operator.And,
		Rules: []operator.IRuleFactory{
			&AtomRule{
				Field:    operator.ArrayElement,
				Operator: operator.Equal,
				Value:    1,
			},
			&CombinedRule{
				Condition: operator.Or,
				Rules: []operator.IRuleFactory{
					&AtomRule{
						Field:    "-1",
						Operator: operator.In,
						Value:    "a",
					},
				},
			},
		},
	}, nil)
	if err == nil {
		t.Errorf("validate should return error")
		return
	}
}

func TestArrayMongoCond(t *testing.T) {
	op := operator.GetOperator(operator.Array)

	// test filter array normal type
	cond, err := op.ToMgo("arr", &CombinedRule{
		Condition: operator.And,
		Rules: []operator.IRuleFactory{
			&AtomRule{
				Field:    operator.ArrayElement,
				Operator: operator.Equal,
				Value:    1,
			},
			&CombinedRule{
				Condition: operator.Or,
				Rules: []operator.IRuleFactory{
					&AtomRule{
						Field:    operator.ArrayElement,
						Operator: operator.NotEqual,
						Value:    "a",
					},
					&AtomRule{
						Field:    operator.ArrayElement,
						Operator: operator.In,
						Value:    []string{"b", "c"},
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("to mongo failed, err: %v", err)
		return
	}

	expectCond := map[string]interface{}{
		operator.DBAND: []map[string]interface{}{{
			"arr": map[string]interface{}{operator.DBEQ: 1},
		}, {
			operator.DBOR: []map[string]interface{}{{
				"arr": map[string]interface{}{operator.DBNE: "a"},
			}, {
				"arr": map[string]interface{}{operator.DBIN: []string{"b", "c"}},
			}},
		}},
	}

	if !reflect.DeepEqual(cond, expectCond) {
		t.Errorf("cond %+v is invalid", cond)
		return
	}
}

func TestArrayMatch(t *testing.T) {
	op := operator.GetOperator(operator.Array)

	// test filter array matched
	matched, err := op.Match([]string{"a", "cc"}, &CombinedRule{
		Condition: operator.And,
		Rules: []operator.IRuleFactory{
			&AtomRule{
				Field:    operator.ArrayElement,
				Operator: operator.Contains,
				Value:    "c",
			},
			&CombinedRule{
				Condition: operator.Or,
				Rules: []operator.IRuleFactory{
					&AtomRule{
						Field:    operator.ArrayElement,
						Operator: operator.NotEqual,
						Value:    "a",
					},
					&AtomRule{
						Field:    operator.ArrayElement,
						Operator: operator.In,
						Value:    []string{"bb", "cc"},
					},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, true, matched)
}
