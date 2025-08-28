package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_SetIsSuperset(t *testing.T) {
	a := mapset.NewSet[int]()
	a.Add(9)
	a.Add(5)
	a.Add(2)
	a.Add(1)
	a.Add(11)

	b := mapset.NewSet[int]()
	b.Add(5)
	b.Add(2)
	b.Add(11)

	// true
	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}

	b.Add(42)

	if a.IsSuperset(b) {
		t.Error("set a should not be a superset of set b because set b has a 42")
	}
}

func Test_SetIsProperSuperset(t *testing.T) {
	a := mapset.NewSet[int]()
	a.Add(5)
	a.Add(2)
	a.Add(11)

	b := mapset.NewSet[int]()
	b.Add(2)
	b.Add(5)
	b.Add(11)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}
	if a.IsProperSuperset(b) {
		t.Error("set a should not be a proper superset of set b (they're equal)")
	}

	a.Add(9)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}
	if !a.IsProperSuperset(b) {
		t.Error("set a not be a proper superset of set b because set a has a 9")
	}

	b.Add(42)

	if a.IsSuperset(b) {
		t.Error("set a should not be a superset of set b because set b has a 42")
	}
	if a.IsProperSuperset(b) {
		t.Error("set a should not be a proper superset of set b because set b has a 42")
	}
}
