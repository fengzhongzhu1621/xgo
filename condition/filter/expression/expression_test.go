package expression

import (
	"encoding/json"
	"testing"

	"github.com/fengzhongzhu1621/xgo/condition/filter/rule"
	"go.mongodb.org/mongo-driver/bson"
)

func TestJsonMarshal(t *testing.T) {
	ruleJson, err := json.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr := Expression{
		exampleRule,
	}
	exprJson, err := json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprJson,
			ruleJson,
		)
		return
	}
}

func TestJsonMarshalNil(t *testing.T) {
	// check if nil expression json equals nil combined rule json
	var rule1 *rule.CombinedRule
	ruleJson, err := json.Marshal(rule1)
	if err != nil {
		t.Error(err)
		return
	}

	var expr *Expression
	exprJson, err := json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprJson,
			ruleJson,
		)
		return
	}

	// check if expression with nil combined rule json equals nil combined rule json
	expr = &Expression{
		IRuleFactory: rule1,
	}
	exprJson, err = json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprJson,
			ruleJson,
		)
		return
	}

	// check if expression with nil atom rule json equals nil atom rule json
	var atomRule *rule.AtomRule
	ruleJson, err = json.Marshal(atomRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr = &Expression{
		IRuleFactory: atomRule,
	}
	exprJson, err = json.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleJson) != string(exprJson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprJson,
			ruleJson,
		)
		return
	}
}

func TestJsonUnmarshal(t *testing.T) {
	exampleExpr := Expression{
		IRuleFactory: exampleRule,
	}
	exprJson, err := json.Marshal(exampleExpr)
	if err != nil {
		t.Error(err)
		return
	}

	expr := new(Expression)
	err = json.Unmarshal(exprJson, expr)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestBsonMarshal(t *testing.T) {
	// TODO test bson marshal Expression value as well as pointer value if bson supports nil pointer with MarshalBSON
	ruleBson, err := bson.Marshal(exampleRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr := &Expression{
		IRuleFactory: exampleRule,
	}
	exprBson, err := bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprBson,
			ruleBson,
		)
		return
	}
}

func TestBsonMarshalNil(t *testing.T) {
	// check if nil expression bson equals nil combined rule bson
	var rule1 *rule.CombinedRule
	ruleBson, err := bson.Marshal(rule1)
	if err != nil {
		t.Error(err)
		return
	}

	var expr *Expression
	exprBson, err := bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprBson,
			ruleBson,
		)
		return
	}

	// check if expression with nil combined rule bson equals nil combined rule bson
	expr = &Expression{
		IRuleFactory: rule1,
	}
	exprBson, err = bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprBson,
			ruleBson,
		)
		return
	}

	// check if expression with nil atom rule bson equals nil atom rule bson
	var atomRule *rule.AtomRule
	ruleBson, err = bson.Marshal(atomRule)
	if err != nil {
		t.Error(err)
		return
	}

	expr = &Expression{
		IRuleFactory: atomRule,
	}
	exprBson, err = bson.Marshal(expr)
	if err != nil {
		t.Error(err)
		return
	}

	if string(ruleBson) != string(exprBson) {
		t.Errorf(
			"expression marshal result %s is not equal to rule marshal result %s",
			exprBson,
			ruleBson,
		)
		return
	}
}

func TestBsonUnmarshal(t *testing.T) {
	exampleExpr := &Expression{
		IRuleFactory: exampleRule,
	}

	exprBson, err := bson.Marshal(exampleExpr)
	if err != nil {
		t.Error(err)
		return
	}

	expr := new(Expression)
	err = bson.Unmarshal(exprBson, expr)
	if err != nil {
		t.Error(err)
		return
	}
}
