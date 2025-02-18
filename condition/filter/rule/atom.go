package rule

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/fengzhongzhu1621/xgo/condition/filter/criteria"
	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/fengzhongzhu1621/xgo/str/stringutils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// UnknownType means it's an unknown type.
	UnknownType operator.RuleType = "Unknown"
	// AtomType means it's an AtomRule
	AtomType operator.RuleType = "AtomRule"
	// CombinedType means it's a CombinedRule
	CombinedType operator.RuleType = "CombinedRule"
)

var _ operator.IRuleFactory = new(AtomRule)

// AtomRule is the basic query rule.
type AtomRule struct {
	Field    string          `json:"field" bson:"field"`
	Operator operator.OpType `json:"operator" bson:"operator"`
	Value    interface{}     `json:"value" bson:"value"`
}

// WithType return the atom rule's type.
func (ar *AtomRule) WithType() operator.RuleType {
	return AtomType
}

// Validate this atom rule is valid or not
func (ar *AtomRule) Validate(opt *operator.ExprOption) error {
	if len(ar.Field) == 0 {
		return errors.New("field is empty")
	}

	// validate operator
	if err := ar.Operator.Validate(); err != nil {
		return err
	}

	if opt == nil {
		return errors.New("validate option must be set")
	}

	// ignore rule fields validation, only validate the operator's value
	if opt.IgnoreRuleFields {
		if err := operator.GetOperator(ar.Operator).ValidateValue(ar.Value, opt); err != nil {
			return fmt.Errorf("%s validate failed, %v", ar.Field, err)
		}
		return nil
	}

	if len(opt.RuleFields) == 0 {
		return errors.New("validate rule fields option must be set")
	}

	typ, exist := opt.RuleFields[ar.Field]
	if !exist {
		return fmt.Errorf("rule field: %s is not exist in the expr option", ar.Field)
	}

	if err := ar.validateValueWithType(opt, typ); err != nil {
		return err
	}

	return nil
}

func (ar *AtomRule) validateValueWithType(opt *operator.ExprOption, typ criteria.FieldType) error {
	childOpt := operator.CloneExprOption(opt)

	// TODO confirm how to deal with object and array and mapstr
	switch ar.Operator {
	case operator.Object:
		if typ != criteria.Object && typ != criteria.MapString {
			return fmt.Errorf("%s is of %s type, should not use operator: %s", ar.Field, typ, ar.Operator)
		}
	case operator.Array:
		if typ != criteria.Array {
			return fmt.Errorf("%s is of %s type, should not use operator: %s", ar.Field, typ, ar.Operator)
		}
	default:
		if err := criteria.ValidateFieldValue(ar.Value, typ); err != nil {
			return fmt.Errorf("invalid %s's value, %v", ar.Field, err)
		}
	}

	switch typ {
	case criteria.Object, criteria.Array:
		ruleFields := make(map[string]criteria.FieldType)
		for field, typ := range opt.RuleFields {
			if strings.HasPrefix(field, ar.Field+".") {
				ruleFields[strings.TrimPrefix(field, ar.Field+".")] = typ
			}
		}
		childOpt.RuleFields = ruleFields

		if err := operator.GetOperator(ar.Operator).ValidateValue(ar.Value, childOpt); err != nil {
			return fmt.Errorf("%s validate failed, %v", ar.Field, err)
		}
	case criteria.MapString:
		childOpt.IgnoreRuleFields = true
	}

	// validate the operator's value
	if err := operator.GetOperator(ar.Operator).ValidateValue(ar.Value, childOpt); err != nil {
		return fmt.Errorf("%s validate failed, %v", ar.Field, err)
	}

	return nil
}

// RuleFields get atom rule's field
func (ar *AtomRule) RuleFields() []string {
	switch ar.Operator {
	// TODO confirm how to deal with these
	case operator.Object, operator.Array:
		// filter object and array operator's fields are its sub-rule fields with its prefix.
		subRule, ok := ar.Value.(operator.IRuleFactory)
		if !ok {
			log.Errorf("%s operator's value(%+v) is not a rule type", ar.Operator, ar.Value)
			return []string{ar.Field}
		}

		subFields := subRule.RuleFields()

		fields := make([]string, len(subFields))
		for idx, field := range subFields {
			fields[idx] = ar.Field + "." + field
		}

		return fields
	}
	return []string{ar.Field}
}

// ToMgo convert this atom rule to a mongo query condition.
func (ar *AtomRule) ToMgo(opts ...*operator.RuleOption) (map[string]interface{}, error) {
	if len(opts) > 0 && opts[0] != nil {
		opt := opts[0]
		if len(opt.Parent) == 0 {
			return nil, errors.New("parent is empty")
		}

		switch opt.ParentType {
		case criteria.Object:
			// 由于目前使用版本的mongodb不支持key中包含的.的查询，所以在这里进行编码；注：需保证数据在存入db的时候，将.以编码的方式存入
			if strings.Contains(ar.Field, ".") {
				ar.Field = stringutils.EncodeDot(ar.Field)
			}
			// add object parent field as prefix to generate object filter rules
			return operator.GetOperator(ar.Operator).ToMgo(opt.Parent+"."+ar.Field, ar.Value)
		case criteria.Array:
			switch ar.Field {
			case operator.ArrayElement:
				// filter array element, matches if any of the elements matches the filter
				return operator.GetOperator(ar.Operator).ToMgo(opt.Parent, ar.Value)
			default:
				return nil, fmt.Errorf("filter array field %s is invalid", ar.Field)
			}
		default:
			return nil, fmt.Errorf("parent type %s is invalid", opt.ParentType)
		}
	}

	return operator.GetOperator(ar.Operator).ToMgo(ar.Field, ar.Value)
}

// Match checks if the input data matches this atomic rule
func (ar *AtomRule) Match(data operator.MatchedData, opts ...*operator.RuleOption) (bool, error) {
	value, err := data.GetValue(ar.Field)
	if err != nil {
		return false, fmt.Errorf("get value by %s field failed, err: %v", ar.Field, err)
	}

	if len(opts) > 0 && opts[0] != nil {
		opt := opts[0]

		switch opt.ParentType {
		case criteria.Array:
			// filter array element, matches if any of the elements matches the filter
			switch reflect.TypeOf(value).Kind() {
			case reflect.Array:
			case reflect.Slice:
			default:
				return false, fmt.Errorf("filter array value(%+v) is not of array type", value)
			}

			v := reflect.ValueOf(value)
			length := v.Len()
			if length == 0 {
				return false, errors.New("value is empty")
			}

			for i := 0; i < length; i++ {
				item := v.Index(i).Interface()

				matched, err := operator.GetOperator(ar.Operator).Match(item, ar.Value)
				if err != nil {
					return false, fmt.Errorf("filter array element(%+v) failed, err: %v", item, err)
				}

				if matched {
					return true, nil
				}
			}
			return false, nil
		default:
			return false, fmt.Errorf("parent type %s is invalid", opt.ParentType)
		}
	}

	matched, err := operator.GetOperator(ar.Operator).Match(value, ar.Value)
	if err != nil {
		return false, fmt.Errorf("match field %s value %v failed, err: %v", ar.Field, value, err)
	}

	return matched, nil
}

type jsonAtomRuleBroker struct {
	Field    string          `json:"field"`
	Operator operator.OpType `json:"operator"`
	Value    json.RawMessage `json:"value"`
}

// UnmarshalJSON unmarshal the json raw message to AtomRule
func (ar *AtomRule) UnmarshalJSON(raw []byte) error {
	br := new(jsonAtomRuleBroker)
	err := json.Unmarshal(raw, br)
	if err != nil {
		return err
	}

	ar.Field = br.Field
	ar.Operator = br.Operator

	if len(br.Value) == 0 {
		return nil
	}

	switch br.Operator {
	case operator.In, operator.NotIn:
		// in and nin operator's value should be an array.
		array := make([]interface{}, 0)
		if err := json.Unmarshal(br.Value, &array); err != nil {
			return err
		}

		ar.Value = array
		return nil
	case operator.Object, operator.Array:
		// filter object and array operator's value should be a rule.
		subRule, err := parseJsonRule(br.Value)
		if err != nil {
			return err
		}
		ar.Value = subRule
		return nil
	}

	to := new(interface{})
	if err := json.Unmarshal(br.Value, to); err != nil {
		return err
	}
	ar.Value = *to

	return nil
}

type bsonAtomRuleBroker struct {
	Field    string          `bson:"field"`
	Operator operator.OpType `bson:"operator"`
	Value    bson.RawValue   `bson:"value"`
}

type bsonAtomRuleCopier struct {
	Field    string          `bson:"field"`
	Operator operator.OpType `bson:"operator"`
	Value    interface{}     `bson:"value"`
}

// MarshalBSON marshal the AtomRule to bson raw message
func (ar *AtomRule) MarshalBSON() ([]byte, error) {
	// right now bson will panic if MarshalBSON is defined using a value receiver and called by a nil pointer
	// TODO this is compatible for nil pointer, but struct marshalling is not supported, find a way to support both
	if ar == nil {
		return bson.Marshal(map[string]interface{}(nil))
	}

	b := bsonAtomRuleCopier{
		Field:    ar.Field,
		Operator: ar.Operator,
		Value:    ar.Value,
	}
	return bson.Marshal(b)
}

// UnmarshalBSON unmarshal the bson raw message to AtomRule
func (ar *AtomRule) UnmarshalBSON(raw []byte) error {
	br := new(bsonAtomRuleBroker)
	err := bson.Unmarshal(raw, br)
	if err != nil {
		return err
	}

	ar.Field = br.Field
	ar.Operator = br.Operator
	switch br.Operator {
	case operator.In, operator.NotIn:
		// in and nin operator's value should be an array.
		array := make([]interface{}, 0)
		if err := br.Value.Unmarshal(&array); err != nil {
			return err
		}

		ar.Value = array
		return nil
	case operator.Object, operator.Array:
		// filter object and array operator's value should be a rule.
		subRule, err := parseBsonRule(br.Value.Document())
		if err != nil {
			return err
		}
		ar.Value = subRule
		return nil
	}

	to := new(interface{})
	if err := br.Value.Unmarshal(to); err != nil {
		return err
	}
	ar.Value = *to

	return nil
}
