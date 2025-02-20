package expression

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/fengzhongzhu1621/xgo/condition/filter/rule"
)

const (
	// MaxRulesDepth defines the maximum number of rules depth
	MaxRulesDepth = uint(3)
)

// Expression is to build a query expression
// 结构体（struct）可以包含接口（interface）类型的字段。
// 这种设计允许结构体实现多态性，即结构体可以根据不同的接口实现来执行不同的行为。
// 接口字段使得结构体能够持有实现了特定接口的任何类型的值，从而提高代码的灵活性和可扩展性。
//
// 嵌入interface作为struct的一个匿名成员，就可以假设这个struct就是此成员interface的一个实现，
// 而不管struct是否已经实现interface所定义的函数。
type Expression struct {
	operator.IRuleFactory
}

// Validate if the expression is valid or not.
func (exp Expression) Validate(opt *operator.ExprOption) error {
	if opt == nil {
		return errors.New("expression's validate option must be set")
	}

	if exp.IRuleFactory == nil {
		return errors.New("expression should not be nil")
	}

	if opt.MaxRulesDepth > opt.MaxRulesDepth {
		return fmt.Errorf("max rule depth exceeds maximum limit: %d", opt.MaxRulesDepth)
	}

	return exp.IRuleFactory.Validate(opt)
}

// MarshalJSON marshal Expression into json value
func (exp Expression) MarshalJSON() ([]byte, error) {
	if exp.IRuleFactory != nil {
		return json.Marshal(exp.IRuleFactory)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON unmarshal Expression from json value
func (exp *Expression) UnmarshalJSON(raw []byte) error {
	rule, err := rule.ParseJsonRule(raw)
	if err != nil {
		return fmt.Errorf("parse rule(%s) failed, err: %v", string(raw), err)
	}

	exp.IRuleFactory = rule
	return nil
}

// String convert expression to string, used for log
func (exp *Expression) String() string {
	if exp == nil {
		return "null"
	}

	jsonVal, err := exp.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("marshal json failed, err: %v", err)
	}

	return string(jsonVal)
}
