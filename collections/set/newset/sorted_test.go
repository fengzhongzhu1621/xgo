package newset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_Sorted(t *testing.T) {
	test := func(t *testing.T, ctor func(vals ...string) mapset.Set[string]) {
		set := ctor("apple", "banana", "pear")
		sorted := mapset.Sorted(set)
		fmt.Println(sorted)
		if len(sorted) != set.Cardinality() {
			t.Errorf(
				"Length of slice is not the same as the set. Expected: %d. Actual: %d",
				set.Cardinality(),
				len(sorted),
			)
		}

		if sorted[0] != "apple" {
			t.Errorf("Element 0 was not equal to apple: %s", sorted[0])
		}
		if sorted[1] != "banana" {
			t.Errorf("Element 1 was not equal to banana: %s", sorted[1])
		}
		if sorted[2] != "pear" {
			t.Errorf("Element 2 was not equal to pear: %s", sorted[2])
		}
	}

	t.Run("Safe", func(t *testing.T) {
		test(t, mapset.NewSet[string]) // [apple banana pear]
	})
	t.Run("Unsafe", func(t *testing.T) { // [apple banana pear]
		test(t, mapset.NewThreadUnsafeSet[string])
	})
}
