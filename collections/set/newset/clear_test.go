package newset

import "testing"

func Test_ClearSet(t *testing.T) {
	a := MakeSetInt([]int{2, 5, 9, 10})

	a.Clear()

	if a.Cardinality() != 0 {
		t.Error("ClearSet should be an empty set")
	}
}
