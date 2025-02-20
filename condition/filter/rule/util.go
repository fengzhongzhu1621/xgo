package rule

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/condition/filter/operator"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseJsonRule(raw []byte) (operator.IRuleFactory, error) {
	// rule with 'condition' key means that it is a combined rule
	if gjson.GetBytes(raw, "condition").Exists() {
		rule := new(CombinedRule)
		err := json.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into combined rule failed, err: %v", err)
		}
		return rule, nil
	}

	// rule with 'operator' key means that it is an atomic rule
	if gjson.GetBytes(raw, "operator").Exists() {
		rule := new(AtomRule)
		err := json.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into atomic rule failed, err: %v", err)
		}
		return rule, nil
	}

	return nil, errors.New("no rule is found")
}

func parseBsonRule(raw []byte) (operator.IRuleFactory, error) {
	// rule with 'condition' key means that it is a combined rule
	if _, ok := bson.Raw(raw).Lookup("condition").StringValueOK(); ok {
		rule := new(CombinedRule)
		err := bson.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into combined rule failed, err: %v", err)
		}
		return rule, nil
	}

	// rule with 'operator' key means that it is an atomic rule
	if _, ok := bson.Raw(raw).Lookup("operator").StringValueOK(); ok {
		rule := new(AtomRule)
		err := bson.Unmarshal(raw, rule)
		if err != nil {
			return nil, fmt.Errorf("unmarshal into atomic rule failed, err: %v", err)
		}
		return rule, nil
	}

	return nil, errors.New("no rule is found")
}
