package operator

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
)

// OpType defines the operators supported by cc.
type OpType string

// Validate test the operator is valid or not.
func (op OpType) Validate() error {
	switch op {
	case Equal,
		NotEqual,
		In,
		NotIn,
		Less,
		LessOrEqual,
		Greater,
		GreaterOrEqual,
		DatetimeLess,
		DatetimeLessOrEqual,
		DatetimeGreater,
		DatetimeGreaterOrEqual,
		BeginsWith,
		BeginsWithInsensitive,
		NotBeginsWith,
		NotBeginsWithInsensitive,
		Contains,
		ContainsSensitive,
		NotContains,
		NotContainsInsensitive,
		EndsWith,
		EndsWithInsensitive,
		NotEndsWith,
		NotEndsWithInsensitive,
		IsEmpty,
		IsNotEmpty,
		Size,
		IsNull,
		IsNotNull,
		Exist,
		NotExist,
		Object,
		Array:
	default:
		return fmt.Errorf("unsupported operator: %s", op)
	}

	return nil
}

const (
	// Unknown is an unsupported operator
	Unknown OpType = "unknown"

	// generic operator

	// Equal operator
	Equal OpType = "equal"
	// NotEqual operator
	NotEqual OpType = "not_equal"

	// set operator that is used to filter element using the value array

	// In operator
	In OpType = "in"
	// NotIn operator
	NotIn OpType = "not_in"

	// numeric compare operator

	// Less operator
	Less OpType = "less"
	// LessOrEqual operator
	LessOrEqual OpType = "less_or_equal"
	// Greater operator
	Greater OpType = "greater"
	// GreaterOrEqual operator
	GreaterOrEqual OpType = "greater_or_equal"

	// datetime operator, ** need to be parsed to mongo in coreservice to avoid json marshaling **

	// DatetimeLess operator
	DatetimeLess OpType = "datetime_less"
	// DatetimeLessOrEqual operator
	DatetimeLessOrEqual OpType = "datetime_less_or_equal"
	// DatetimeGreater operator
	DatetimeGreater OpType = "datetime_greater"
	// DatetimeGreaterOrEqual operator
	DatetimeGreaterOrEqual OpType = "datetime_greater_or_equal"

	// string operator

	// BeginsWith operator with case-sensitive
	BeginsWith OpType = "begins_with"
	// BeginsWithInsensitive operator with case-insensitive
	BeginsWithInsensitive OpType = "begins_with_i"
	// NotBeginsWith operator with case-sensitive
	NotBeginsWith OpType = "not_begins_with"
	// NotBeginsWithInsensitive operator with case-insensitive
	NotBeginsWithInsensitive OpType = "not_begins_with_i"
	// Contains operator with case-insensitive, compatible for the query builder's same operator that's case-insensitive
	Contains OpType = "contains"
	// ContainsSensitive operator with case-sensitive
	ContainsSensitive OpType = "contains_s"
	// NotContains operator with case-sensitive
	NotContains OpType = "not_contains"
	// NotContainsInsensitive operator with case-insensitive
	NotContainsInsensitive OpType = "not_contains_i"
	// EndsWith operator with case-sensitive
	EndsWith OpType = "ends_with"
	// EndsWithInsensitive operator with case-insensitive
	EndsWithInsensitive OpType = "ends_with_i"
	// NotEndsWith operator with case-sensitive
	NotEndsWith OpType = "not_ends_with"
	// NotEndsWithInsensitive operator with case-insensitive
	NotEndsWithInsensitive OpType = "not_ends_with_i"

	// array operator

	// IsEmpty operator
	IsEmpty OpType = "is_empty"
	// IsNotEmpty operator
	IsNotEmpty OpType = "is_not_empty"
	// Size operator
	Size OpType = "size"

	// null check operator

	// IsNull operator
	IsNull OpType = "is_null"
	// IsNotNull operator
	IsNotNull OpType = "is_not_null"

	// existence check operator

	// Exist operator
	Exist OpType = "exist"
	// NotExist operator
	NotExist OpType = "not_exist"

	// filter embedded elements operator

	// Object filter object fields operator
	Object OpType = "filter_object"
	// Array filter array elements operator
	Array OpType = "filter_array"
)

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
	// WithType get a rule's type 规则类型
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
