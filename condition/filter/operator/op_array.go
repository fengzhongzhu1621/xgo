package operator

import (
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
)

const (
	ArrayElement = "element"
)

// ArrayOp is filter array operator
type ArrayOp OpType

// Name is filter array operator name
// 数组中的每个元素是否匹配指定规则
func (o ArrayOp) Name() OpType {
	return Array
}

// ValidateValue validate filter array operator value
func (o ArrayOp) ValidateValue(v interface{}, opt *ExprOption) error {

	// filter array operator's value is the sub-rule to filter the array's field.
	subRule, ok := v.(IRuleFactory)
	if !ok {
		return fmt.Errorf("filter array operator's value(%+v) is not a rule type", v)
	}

	// validate filter array rule depth, then continues to validate children rule depth
	if opt == nil {
		return errors.New("validate option must be set")
	}

	if opt.MaxRulesDepth <= 1 {
		return fmt.Errorf("expression rules depth exceeds maximum")
	}

	childOpt := CloneExprOption(opt)
	childOpt.MaxRulesDepth = opt.MaxRulesDepth - 1

	if err := subRule.Validate(childOpt); err != nil {
		return fmt.Errorf("invalid value(%+v), err: %v", v, err)
	}

	return nil
}

// ToMgo convert the filter array operator's field and value to a mongo query condition.
func (o ArrayOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	subRule, ok := value.(IRuleFactory)
	if !ok {
		return nil, fmt.Errorf("filter array operator's value(%+v) is not a rule type", value)
	}

	parentOpt := &RuleOption{
		Parent:     field,
		ParentType: criteria.Array,
	}

	return subRule.ToMgo(parentOpt)
}

// Match checks if the first data matches the second data by this operator
func (o ArrayOp) Match(value1, value2 interface{}) (bool, error) {
	if value1 == nil {
		return false, errors.New("input value is nil")
	}

	subRule, ok := value2.(IRuleFactory)
	if !ok {
		return false, fmt.Errorf("filter array operator's value(%+v) is not a rule type", value2)
	}

	val := mapstr.MapStr{
		ArrayElement: value1,
	}

	parentOpt := &RuleOption{
		ParentType: criteria.Array,
	}

	return subRule.Match(val, parentOpt)
}
