package operator

import "github.com/fengzhongzhu1621/xgo/condition/filter/criteria"

// RuleType is the expression rule's rule type.
type RuleType string

// RuleOption defines the options of a rule.
type RuleOption struct {
	// Parent field name, used when filtering object/array elements
	Parent string
	// ParentType parent type, used when filtering object/array elements
	ParentType criteria.FieldType
}

type IRuleFactory interface {
	// WithType get a rule's type
	WithType() RuleType
	// Validate this rule is valid or not
	Validate(opt *ExprOption) error
	// RuleFields get this rule's fields
	RuleFields() []string
	// ToMgo convert this rule to a mongo condition
	ToMgo(opt ...*RuleOption) (map[string]interface{}, error)
	// Match checks if the input data matches this rule
	Match(data MatchedData, opt ...*RuleOption) (bool, error)
}
