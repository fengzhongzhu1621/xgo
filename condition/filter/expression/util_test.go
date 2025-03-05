package expression

import (
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/fengzhongzhu1621/xgo/condition/filter/rule"
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
	exampleRule = &rule.CombinedRule{
		Condition: operator.And,
		Rules: []operator.IRuleFactory{
			&rule.AtomRule{
				Field:    "test",
				Operator: operator.Equal,
				Value:    1,
			},
			&rule.CombinedRule{
				Condition: operator.Or,
				Rules: []operator.IRuleFactory{
					&rule.AtomRule{
						Field:    "test1",
						Operator: operator.Array,
						Value: &rule.AtomRule{
							Field:    operator.ArrayElement,
							Operator: operator.Object,
							Value: &rule.CombinedRule{
								Condition: operator.And,
								Rules: []operator.IRuleFactory{
									&rule.AtomRule{
										Field:    "test2",
										Operator: operator.In,
										Value:    []string{"b", "c"},
									},
								},
							},
						},
					},
					&rule.AtomRule{
						Field:    "test3",
						Operator: operator.DatetimeLess,
						Value:    1,
					},
				},
			},
		},
	}
)
