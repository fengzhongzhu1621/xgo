package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_SetIsSubset(t *testing.T) {
	a := MakeSetInt([]int{1, 2, 3, 5, 7})

	b := mapset.NewSet[int]()
	b.Add(3)
	b.Add(5)
	b.Add(7)

	if !b.IsSubset(a) {
		t.Error("set b should be a subset of set a")
	}

	b.Add(72)

	if b.IsSubset(a) {
		t.Error(
			"set b should not be a subset of set a because it contains 72 which is not in the set of a",
		)
	}
}

func Test_SetIsProperSubset(t *testing.T) {
	a := MakeSetInt([]int{1, 2, 3, 5, 7})
	b := MakeSetInt([]int{7, 5, 3, 2, 1})

	// true
	if !a.IsSubset(b) {
		t.Error("set a should be a subset of set b")
	}
	// false
	if a.IsProperSubset(b) {
		t.Error("set a should not be a proper subset of set b (they're equal)")
	}

	b.Add(72)

	// true
	if !a.IsSubset(b) {
		t.Error("set a should be a subset of set b")
	}
	// false
	if !a.IsProperSubset(b) {
		t.Error("set a should be a proper subset of set b")
	}
}
