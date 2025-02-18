package rule

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var _ operator.IRuleFactory = new(CombinedRule)

// CombinedRule is the compound query rule combined by many rules.
type CombinedRule struct {
	Condition operator.LogicOperator  `json:"condition" bson:"condition"`
	Rules     []operator.IRuleFactory `json:"rules" bson:"rules"`
}

// WithType return the combined rule's tye.
func (cr *CombinedRule) WithType() operator.RuleType {
	return CombinedType
}

// Validate the combined rule
func (cr *CombinedRule) Validate(opt *operator.ExprOption) error {
	if err := cr.Condition.Validate(); err != nil {
		return err
	}

	if len(cr.Rules) == 0 {
		return errors.New("combined rules shouldn't be empty")
	}

	if opt == nil {
		return errors.New("validate option must be set")
	}

	if len(cr.Rules) > int(opt.MaxRulesLimit) {
		return fmt.Errorf("rules elements number exceeds limit: %d", opt.MaxRulesLimit)
	}

	// validate combined rule depth, then continues to validate children rule depth
	if opt.MaxRulesDepth <= 1 {
		return fmt.Errorf("expression rules depth exceeds maximum")
	}

	childOpt := operator.CloneExprOption(opt)
	childOpt.MaxRulesDepth = opt.MaxRulesDepth - 1

	for _, one := range cr.Rules {
		if err := one.Validate(childOpt); err != nil {
			return err
		}
	}

	return nil
}

// RuleFields get combined rule's fields
func (cr *CombinedRule) RuleFields() []string {
	fields := make([]string, 0)
	for _, rule := range cr.Rules {
		fields = append(fields, rule.RuleFields()...)
	}
	return fields
}

// ToMgo convert the combined rule to a mongo query condition.
func (cr *CombinedRule) ToMgo(opt ...*operator.RuleOption) (map[string]interface{}, error) {
	if err := cr.Condition.Validate(); err != nil {
		return nil, err
	}

	if len(cr.Rules) == 0 {
		return nil, errors.New("combined rules shouldn't be empty")
	}

	filters := make([]map[string]interface{}, 0)
	for idx, rule := range cr.Rules {
		filter, err := rule.ToMgo(opt...)
		if err != nil {
			return nil, fmt.Errorf("rules[%d] is invalid, err: %v", idx, err)
		}
		filters = append(filters, filter)
	}

	switch cr.Condition {
	case operator.Or:
		return map[string]interface{}{operator.DBOR: filters}, nil
	case operator.And:
		return map[string]interface{}{operator.DBAND: filters}, nil
	default:
		return nil, fmt.Errorf("unexpected operator %s", cr.Condition)
	}
}

// Match checks if the input data matches this combined rule
func (cr *CombinedRule) Match(data operator.MatchedData, opts ...*operator.RuleOption) (bool, error) {
	if err := cr.Condition.Validate(); err != nil {
		return false, err
	}

	if len(cr.Rules) == 0 {
		return false, errors.New("combined rules shouldn't be empty")
	}

	for idx, rule := range cr.Rules {
		matched, err := rule.Match(data, opts...)
		if err != nil {
			log.Errorf("match rules[%d] failed, err: %v, data: %+v", idx, err, data)
		}

		switch cr.Condition {
		case operator.Or:
			if matched {
				return true, nil
			}
		case operator.And:
			if !matched {
				return false, nil
			}
		}
	}

	switch cr.Condition {
	case operator.Or:
		return false, nil
	case operator.And:
		return true, nil
	default:
		return false, fmt.Errorf("unexpected operator %s", cr.Condition)
	}
}

type jsonCombinedRuleBroker struct {
	Condition operator.LogicOperator `json:"condition"`
	Rules     []json.RawMessage      `json:"rules"`
}

// UnmarshalJSON unmarshal the json raw message to AtomRule
func (cr *CombinedRule) UnmarshalJSON(raw []byte) error {
	broker := new(jsonCombinedRuleBroker)

	err := json.Unmarshal(raw, broker)
	if err != nil {
		return fmt.Errorf("unmarshal into combined rule failed, err: %v", err)
	}

	cr.Condition = broker.Condition
	cr.Rules = make([]operator.IRuleFactory, len(broker.Rules))

	for idx, rawRule := range broker.Rules {
		rule, err := parseJsonRule(rawRule)
		if err != nil {
			return fmt.Errorf("parse rules[%d] %s failed, err: %v", idx, string(rawRule), err)
		}
		cr.Rules[idx] = rule
	}

	return nil
}

type bsonCombinedRuleBroker struct {
	Condition operator.LogicOperator `bson:"condition"`
	Rules     []bson.Raw             `bson:"rules"`
}

// MarshalBSON marshal the bson raw message to CombinedRule
func (cr *CombinedRule) MarshalBSON() ([]byte, error) {
	// right now bson will panic if MarshalBSON is defined using a value receiver and called by a nil pointer
	// TODO this is compatible for nil pointer, but struct marshalling is not supported, find a way to support both
	if cr == nil {
		return bson.Marshal(map[string]interface{}(nil))
	}

	b := bsonCombinedRuleBroker{
		Condition: cr.Condition,
		Rules:     make([]bson.Raw, len(cr.Rules)),
	}

	for index, value := range cr.Rules {
		bsonVal, err := bson.Marshal(value)
		if err != nil {
			return nil, err
		}
		b.Rules[index] = bsonVal
	}

	return bson.Marshal(b)
}

// UnmarshalBSON unmarshal the bson raw message to CombinedRule
func (cr *CombinedRule) UnmarshalBSON(raw []byte) error {
	broker := new(bsonCombinedRuleBroker)

	err := bson.Unmarshal(raw, broker)
	if err != nil {
		return fmt.Errorf("unmarshal into combined rule failed, err: %v", err)
	}

	cr.Condition = broker.Condition
	cr.Rules = make([]operator.IRuleFactory, len(broker.Rules))

	for idx, rawRule := range broker.Rules {
		rule, err := parseBsonRule(rawRule)
		if err != nil {
			return fmt.Errorf("parse rules[%d] %s failed, err: %v", idx, string(rawRule), err)
		}
		cr.Rules[idx] = rule
	}

	return nil
}
