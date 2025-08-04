package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_AddSetNoDuplicate(t *testing.T) {
	a := MakeSetInt([]int{7, 5, 3, 7})

	if a.Cardinality() != 3 {
		t.Error("AddSetNoDuplicate set should have 3 elements since 7 is a duplicate")
	}

	if !a.Contains(7) || !a.Contains(5) || !a.Contains(3) {
		t.Error("AddSetNoDuplicate set should have a 7, 5, and 3 in it.")
	}
}

func Test_ContainsMultipleSet(t *testing.T) {
	a := MakeSetInt([]int{8, 6, 7, 5, 3, 0, 9})

	if !a.Contains(8, 6, 7, 5, 3, 0, 9) {
		t.Error("ContainsAll should contain Jenny's phone number")
	}

	if a.Contains(8, 6, 11, 5, 3, 0, 9) {
		t.Error("ContainsAll should not have all of these numbers")
	}
}

func Test_ContainsOneSet(t *testing.T) {
	a := mapset.NewSet[int]()

	a.Add(71)

	if !a.ContainsOne(71) {
		t.Error("ContainsSet should contain 71")
	}

	a.Remove(71)

	if a.ContainsOne(71) {
		t.Error("ContainsSet should not contain 71")
	}

	a.Add(13)
	a.Add(7)
	a.Add(1)

	if !a.ContainsOne(13) || !a.ContainsOne(7) || !a.ContainsOne(1) {
		t.Error("ContainsSet should contain 13, 7, 1")
	}
}

func Test_ContainsAnySet(t *testing.T) {
	a := mapset.NewSet[int]()

	a.Add(71)

	if !a.ContainsAny(71) {
		t.Error("ContainsSet should contain 71")
	}

	if !a.ContainsAny(71, 10) {
		t.Error("ContainsSet should contain 71 or 10")
	}

	a.Remove(71)

	if a.ContainsAny(71) {
		t.Error("ContainsSet should not contain 71")
	}

	if a.ContainsAny(71, 10) {
		t.Error("ContainsSet should not contain 71 or 10")
	}

	a.Add(13)
	a.Add(7)
	a.Add(1)

	if !(a.ContainsAny(13, 17, 10)) {
		t.Error("ContainsSet should contain 13, 17, or 10")
	}
}

//func Test_ContainsAnyElement(t *testing.T) {
//	a := mapset.NewSet[int]()
//	a.Add(1)
//	a.Add(3)
//	a.Add(5)
//
//	b := mapset.NewSet[int]()
//	a.Add(2)
//	a.Add(4)
//	a.Add(6)
//
//	if ret := a.ContainsAnyElement(b); ret {
//		t.Errorf("set a not contain any element in set b")
//	}
//
//	a.Add(10)
//
//	if ret := a.ContainsAnyElement(b); ret {
//		t.Errorf("set a not contain any element in set b")
//	}
//
//	b.Add(10)
//
//	if ret := a.ContainsAnyElement(b); !ret {
//		t.Errorf("set a contain 10")
//	}
//}
