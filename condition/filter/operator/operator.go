package operator

import "github.com/fengzhongzhu1621/xgo/condition/filter/criteria"

const (
	// MaxRulesDepth defines the maximum number of rules depth
	MaxRulesDepth = uint(3)
)

// FieldType define the table's field data type.
type FieldType string

// ExprOption defines how to validate an expression.
type ExprOption struct {
	// RuleFields:
	// 1. used to test if all the expression rule's field
	//    is in the RuleFields' key restricts.
	// 2. all the expression's rule field should be a sub-set
	//    of the RuleFields' key.
	RuleFields map[string]criteria.FieldType
	// IgnoreRuleFields defines if expression rule field validation needs to be ignored.
	IgnoreRuleFields bool
	// MaxInLimit defines the maximum element of the in operator
	MaxInLimit uint
	// MaxNotInLimit defines the maximum element of the nin operator
	MaxNotInLimit uint
	// MaxRulesLimit defines the maximum number of rules an expression allows.
	MaxRulesLimit uint
	// MaxRulesDepth defines the maximum depth of rules an expression allows.
	MaxRulesDepth uint
}

// NewDefaultExprOpt init an expression option with default limit option.
func NewDefaultExprOpt(ruleFields map[string]criteria.FieldType) *ExprOption {
	return &ExprOption{
		RuleFields:    ruleFields,
		MaxInLimit:    500,
		MaxNotInLimit: 500,
		MaxRulesLimit: 50,
		MaxRulesDepth: MaxRulesDepth,
	}
}

func CloneExprOption(opt *ExprOption) *ExprOption {
	return &ExprOption{
		RuleFields:       opt.RuleFields,
		IgnoreRuleFields: opt.IgnoreRuleFields,
		MaxInLimit:       opt.MaxInLimit,
		MaxNotInLimit:    opt.MaxNotInLimit,
		MaxRulesLimit:    opt.MaxRulesLimit,
		MaxRulesDepth:    opt.MaxRulesDepth,
	}
}

// Operator is a collection of supported query operators.
type IOperator interface {
	// Name is the operator's name
	Name() OpType
	// ValidateValue validate the operator's value is valid or not
	ValidateValue(v interface{}, opt *ExprOption) error
	// ToMgo generate an operator's mongo condition with its field and value.
	ToMgo(field string, value interface{}) (map[string]interface{}, error)
	// Match checks if the first data matches the second data by this operator
	Match(value1, value2 interface{}) (bool, error)
}
