package stringutils

import (
	"strconv"
	"strings"
	"testing"

	"github.com/fengzhongzhu1621/xgo/crypto/randutils"
	"github.com/stretchr/testify/assert"
)

func TestHead(t *testing.T) {
	s := "abc__def"
	sep := "__"
	left, right := Head(s, sep)
	assert.Equal(t, left, "abc")
	assert.Equal(t, right, "def")
}

func TestRemoveDuplicateElement(t *testing.T) {
	items := []string{"a", "b", "a"}
	dropDuplicatedItems := RemoveDuplicateElement(items)
	assert.Equal(t, dropDuplicatedItems, []string{"a", "b"})
}

func TestReflectReverseSlice(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	ReflectReverseSlice(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, names)
}

func TestReverseSliceGetNew(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	newNames := ReverseSliceGetNew(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, newNames)
}

func TestReverseSlice(t *testing.T) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	ReverseSlice(names)
	expected := []string{"g", "f", "e", "d", "c", "b", "a"}
	assert.Equal(t, expected, names)
}

func TestGenerateId(t *testing.T) {
	actual := randutils.GenerateID()
	s, err := strconv.ParseUint(actual, 10, 64)
	assert.Equal(t, err, nil)
	assert.Equal(t, s > 0, true)
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

func TestMd5(t *testing.T) {
	src := "123456789"
	actual := randutils.Md5(src)
	expect := "25f9e794323b453885f5181f1b624d0b"
	assert.Equal(t, expect, actual)
}

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

func BenchmarkToLower(b *testing.B) {
	string1 := "ABCDEFGHIJKLMN"
	string2 := "abcdefghijklmn"
	for i := 0; i < b.N; i++ {
		ToLower(string1)
		ToLower(string2)
	}
}

func BenchmarkBuildInToLower(b *testing.B) {
	string1 := "ABCDEFGHIJKLMN"
	string2 := "abcdefghijklmn"
	for i := 0; i < b.N; i++ {
		strings.ToLower(string1)
		strings.ToLower(string2)
	}
}
