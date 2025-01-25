package newset

import "testing"

func Test_RemoveSet(t *testing.T) {
	a := MakeSetInt([]int{6, 3, 1})

	a.Remove(3)

	if a.Cardinality() != 2 {
		t.Error("RemoveSet should only have 2 items in the set")
	}

	if !(a.Contains(6) && a.Contains(1)) {
		t.Error("RemoveSet should have only items 6 and 1 in the set")
	}

	a.Remove(6)
	a.Remove(1)

	if a.Cardinality() != 0 {
		t.Error("RemoveSet should be an empty set after removing 6 and 1")
	}
}

func Test_RemoveAllSet(t *testing.T) {
	a := MakeSetInt([]int{6, 3, 1, 8, 9})

	a.RemoveAll(3, 1)

	if a.Cardinality() != 3 {
		t.Error("RemoveAll should only have 2 items in the set")
	}

	if !a.Contains(6, 8, 9) {
		t.Error("RemoveAll should have only items (6,8,9) in the set")
	}

	a.RemoveAll(6, 8, 9)

	if a.Cardinality() != 0 {
		t.Error("RemoveSet should be an empty set after removing 6 and 1")
	}
}
