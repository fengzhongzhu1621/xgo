package newset

import (
	mapset "github.com/deckarep/golang-set/v2"
)

func MakeUnsafeSetInt(ints []int) mapset.Set[int] {
	s := mapset.NewThreadUnsafeSet[int]()
	for _, i := range ints {
		s.Add(i)
	}
	return s
}

func MakeUnsafeSetIntWithAppend(ints ...int) mapset.Set[int] {
	s := mapset.NewThreadUnsafeSet[int]()
	s.Append(ints...)
	return s
}
