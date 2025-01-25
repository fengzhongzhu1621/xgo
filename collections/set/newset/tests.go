package newset

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func assertEqual[T comparable](a, b mapset.Set[T], t *testing.T) {
	if !a.Equal(b) {
		t.Errorf("%v != %v\n", a, b)
	}
}
