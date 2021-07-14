package string_utils

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
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
	actual := GenerateId()
	s, err := strconv.ParseUint(actual, 10, 64)
	assert.Equal(t, err, nil)
	assert.Equal(t, s > 0, true)
}

func BenchmarkCompareStringSliceReflect(b *testing.B) {
	sliceA := []string{"a", "b", "c", "d", "e"}
	sliceB := []string{"e", "d", "c", "b", "a"}
	for n := 0; n < b.N; n++ {
		CompareStringSliceReflect(sliceA, sliceB)
	}
}

func BenchmarkCompareStringSlice(b *testing.B) {
	sliceA := []string{"a", "b", "c", "d", "e"}
	sliceB := []string{"e", "d", "c", "b", "a"}
	for n := 0; n < b.N; n++ {
		CompareStringSlice(sliceA, sliceB)
	}
}

func BenchmarkReverseReflectSlice(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReflectReverseSlice(names)
	}
}

func BenchmarkReverseSlice(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReverseSlice(names)
	}
}

func BenchmarkReverseSliceNew(b *testing.B) {
	names := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i := 0; i < b.N; i++ {
		ReverseSliceGetNew(names)
	}
}
