package newset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_ToSlice(t *testing.T) {
	s := MakeSetInt([]int{1, 2, 3})
	setAsSlice := s.ToSlice()
	if len(setAsSlice) != s.Cardinality() {
		t.Errorf("Set length is incorrect: %v", len(setAsSlice))
	}

	for _, i := range setAsSlice {
		if !s.Contains(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}

func Test_NewSetFromMapKey_Ints(t *testing.T) {
	m := map[int]int{
		5: 500,
		2: 300,
	}

	s := mapset.NewSetFromMapKeys(m)
	fmt.Println(s) // Set{5, 2}

	if len(m) != s.Cardinality() {
		t.Errorf("Length of Set is not the same as the map. Expected: %d. Actual: %d", len(m), s.Cardinality())
	}

	for k := range m {
		if !s.Contains(k) {
			t.Errorf("Element %d not found in map: %v", k, m)
		}
	}
}
