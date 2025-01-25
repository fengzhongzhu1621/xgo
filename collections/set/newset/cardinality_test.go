package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestCardinality(t *testing.T) {
	a := mapset.NewSet[int]()
	if a.Cardinality() != 0 {
		t.Error("NewSet should start out as an empty set")
	}

	assertEqual(mapset.NewSet([]int{}...), mapset.NewSet[int](), t)
	assertEqual(mapset.NewSet([]int{1}...), mapset.NewSet(1), t)
	assertEqual(mapset.NewSet([]int{1, 2}...), mapset.NewSet(1, 2), t)
	assertEqual(mapset.NewSet([]string{"a"}...), mapset.NewSet("a"), t)
	assertEqual(mapset.NewSet([]string{"a", "b"}...), mapset.NewSet("a", "b"), t)
}
