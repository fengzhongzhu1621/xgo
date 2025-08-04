package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_SetIntersect(t *testing.T) {
	a := mapset.NewSet[int]()
	a.Add(1)
	a.Add(3)
	a.Add(5)

	b := mapset.NewSet[int]()
	a.Add(2)
	a.Add(4)
	a.Add(6)

	c := a.Intersect(b)

	if c.Cardinality() != 0 {
		t.Error("set c should be the empty set because there is no common items to intersect")
	}

	a.Add(10)
	b.Add(10)

	d := a.Intersect(b)

	if d.Cardinality() != 1 || !d.Contains(10) {
		t.Error("set d should have a size of 1 and contain the item 10")
	}
}
