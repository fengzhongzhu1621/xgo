package newset

import (
	mapset "github.com/deckarep/golang-set/v2"
)

func MakeSetInt(ints []int) mapset.Set[int] {
	s := mapset.NewSet[int]()
	for _, i := range ints {
		s.Add(i)
	}
	return s
}

func MakeSetIntWithAppend(ints ...int) mapset.Set[int] {
	s := mapset.NewSet[int]()
	s.Append(ints...)
	return s
}
