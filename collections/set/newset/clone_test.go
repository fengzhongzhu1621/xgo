package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_SetClone(t *testing.T) {
	a := mapset.NewSet[int]()
	a.Add(1)
	a.Add(2)

	b := a.Clone()

	if !a.Equal(b) {
		t.Error("Clones should be equal")
	}

	a.Add(3)
	if a.Equal(b) {
		t.Error("a contains one more element, they should not be equal")
	}

	c := a.Clone()
	c.Remove(1)

	if a.Equal(c) {
		t.Error("C contains one element less, they should not be equal")
	}
}
