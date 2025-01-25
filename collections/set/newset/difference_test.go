package newset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_SetDifference(t *testing.T) {
	a := mapset.NewSet[int]()
	a.Add(1)
	a.Add(2)
	a.Add(3)

	b := mapset.NewSet[int]()
	b.Add(1)
	b.Add(3)
	b.Add(4)
	b.Add(5)
	b.Add(6)
	b.Add(99)

	// Set{2}
	c := a.Difference(b)
	fmt.Println(c)

	if !(c.Cardinality() == 1 && c.Contains(2)) {
		t.Error("the difference of set a to b is the set of 1 item: 2")
	}
}

func Test_SetSymmetricDifference(t *testing.T) {
	a := mapset.NewSet[int]()
	a.Add(1)
	a.Add(2)
	a.Add(3)
	a.Add(45)

	b := mapset.NewSet[int]()
	b.Add(1)
	b.Add(3)
	b.Add(4)
	b.Add(5)
	b.Add(6)
	b.Add(99)

	// Set{99, 2, 45, 4, 5, 6}
	c := a.SymmetricDifference(b)
	fmt.Println(c)

	if !(c.Cardinality() == 6 && c.Contains(2) && c.Contains(45) && c.Contains(4) && c.Contains(5) && c.Contains(6) && c.Contains(99)) {
		t.Error("the symmetric difference of set a to b is the set of 6 items: 2, 45, 4, 5, 6, 99")
	}
}
