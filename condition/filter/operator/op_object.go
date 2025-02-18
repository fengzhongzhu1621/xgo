package operator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/maps/mapstr"
	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
	jsonx "github.com/fengzhongzhu1621/xgo/crypto/encoding/json"
)

// ObjectOp is filter object operator
type ObjectOp OpType

// Name is filter object operator name
func (o ObjectOp) Name() OpType {
	return Object
}

// ValidateValue validate filter object operator value
func (o ObjectOp) ValidateValue(v interface{}, opt *ExprOption) error {
	// filter object operator's value is the sub-rule to filter the object's field.
	subRule, ok := v.(IRuleFactory)
	if !ok {
		return fmt.Errorf("filter object operator's value(%+v) is not a rule type", v)
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

// ToMgo convert the filter object operator's field and value to a mongo query condition.
func (o ObjectOp) ToMgo(field string, value interface{}) (map[string]interface{}, error) {
	if len(field) == 0 {
		return nil, errors.New("field is empty")
	}

	subRule, ok := value.(IRuleFactory)
	if !ok {
		return nil, fmt.Errorf("filter object operator's value(%+v) is not a rule type", value)
	}

	parentOpt := &RuleOption{
		Parent:     field,
		ParentType: criteria.Object,
	}

	return subRule.ToMgo(parentOpt)
}

// Match checks if the first data matches the second data by this operator
func (o ObjectOp) Match(value1, value2 interface{}) (bool, error) {
	subRule, ok := value2.(IRuleFactory)
	if !ok {
		return false, fmt.Errorf("filter object operator's value(%+v) is not a rule type", value2)
	}

	switch t := value1.(type) {
	case MatchedData:
		return subRule.Match(t)
	case map[string]interface{}:
		return subRule.Match(mapstr.MapStr(t))
	case mapstr.MapStr:
		return subRule.Match(mapstr.MapStr(t))
	case string:
		return subRule.Match(jsonx.JsonString(t))
	case json.RawMessage:
		return subRule.Match(jsonx.JsonString(t))
	default:
		return false, fmt.Errorf("filter object operator's input value(%+v) is not an object type", value1)
	}
}
