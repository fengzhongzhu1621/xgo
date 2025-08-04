package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_SetUnion(t *testing.T) {
	a := mapset.NewSet[int]()

	b := mapset.NewSet[int]()
	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	b.Add(5)

	c := a.Union(b)

	if c.Cardinality() != 5 {
		t.Error("set c is unioned with an empty set and therefore should have 5 elements in it")
	}

	d := mapset.NewSet[int]()
	d.Add(10)
	d.Add(14)
	d.Add(0)

	e := c.Union(d)
	if e.Cardinality() != 8 {
		t.Error("set e should have 8 elements in it after being unioned with set c to d")
	}

	f := mapset.NewSet[int]()
	f.Add(14)
	f.Add(3)

	g := f.Union(e)
	if g.Cardinality() != 8 {
		t.Error(
			"set g should still have 8 elements in it after being unioned with set f that has duplicates",
		)
	}
}
