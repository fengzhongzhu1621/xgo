package rule

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	jsonx "github.com/fengzhongzhu1621/xgo/crypto/encoding/json"
	"github.com/stretchr/testify/assert"
)

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
	matched, err := exampleRule.Match(jsonx.JsonString(`{"test":1,"test1":[{"test2":"b"}],"test3":111}`))
	assert.NoError(t, err)
	assert.Equal(t, true, matched)

	matched, err = exampleRule.Match(jsonx.JsonString(`{"test":1,"test1":[{"test2":"a"}],"test3":111}`))
	assert.NoError(t, err)
	assert.Equal(t, false, matched)
}
