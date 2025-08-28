package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValueInBraces(t *testing.T) {
	s := "ab{cd}ef{gh}i"
	actual := GetValueInBraces(s)
	expect := "cd"
	assert.Equal(t, expect, actual)

	s = "cd"
	actual = GetValueInBraces(s)
	expect = "cd"
	assert.Equal(t, expect, actual)
}

func TestMergeGetAndPostParamWithKey(t *testing.T) {
	queryParam := map[string]string{"b": "2", "a": "1"}
	postParam := map[string]string{"c": "3", "d": "4"}
	key := "123456789"
	keyName := "_key"
	actual := MergeGetAndPostParamWithKey(queryParam, postParam, key, keyName)
	expect := "a=1&b=2&c=3&d=4&_key=123456789"
	assert.Equal(t, expect, actual)
}
