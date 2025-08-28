package newset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_PopSafe(t *testing.T) {
	a := mapset.NewSet[string]()

	a.Add("a")
	a.Add("b")
	a.Add("c")
	a.Add("d")
	fmt.Println(a) // Set{a, b, c, d}

	aPop := func() (v string) {
		v, _ = a.Pop()
		return
	}

	captureSet := mapset.NewSet[string]()
	captureSet.Add(aPop())
	captureSet.Add(aPop())
	captureSet.Add(aPop())
	captureSet.Add(aPop())
	finalNil := aPop()

	if captureSet.Cardinality() != 4 {
		t.Error("unexpected captureSet cardinality; should be 4")
	}

	if a.Cardinality() != 0 {
		t.Error("unepxected a cardinality; should be zero")
	}

	if !captureSet.Contains("c", "a", "d", "b") {
		t.Error("unexpected result set; should be a,b,c,d (any order is fine")
	}

	if finalNil != "" {
		t.Error("when original set is empty, further pops should result in nil")
	}
}
