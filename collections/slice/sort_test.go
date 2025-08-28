package slice

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/araujo88/lambda-go/pkg/sortgroup"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
)

// TestSort Sorts a slice of any ordered type(number or string),
// use quick sort algrithm. Default sort order is ascending (asc),
// if want descending order, set param `sortOrder` to `desc`. Ordered type: number(all ints uints floats) or string.
// func Sort[T constraints.Ordered](slice []T, sortOrder ...string)
func TestSort(t *testing.T) {
	{
		s := []int{4, 2, 3, 1}
		sort.Ints(s)
		fmt.Println(s) // [1 2 3 4]
	}

	{
		numbers := []int{1, 4, 3, 2, 5}
		slice.Sort(numbers)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, numbers)

		slice.Sort(numbers, "desc")
		assert.Equal(t, []int{5, 4, 3, 2, 1}, numbers)

		strings := []string{"a", "d", "c", "b", "e"}
		slice.Sort(strings)
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, strings)

		slice.Sort(strings, "desc")
		assert.Equal(t, []string{"e", "d", "c", "b", "a"}, strings)
	}

	{
		tests := []struct {
			name  string
			slice []int
			want  []int
		}{
			{"sort positive numbers", []int{5, 3, 4, 1, 2}, []int{1, 2, 3, 4, 5}},
			{"sort including negatives", []int{-1, -3, 2, 1, 0}, []int{-3, -1, 0, 1, 2}},
			{"empty slice", []int{}, []int{}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				sortgroup.SortSlice(tt.slice) // Modify the slice in-place
				if !reflect.DeepEqual(tt.slice, tt.want) {
					t.Errorf("sortSlice() resulted in %v, want %v", tt.slice, tt.want)
				}
			})
		}
	}

	{
		strings := []string{"a", "d", "c", "b", "e"}
		sort.Strings(strings)
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, strings)
	}
}

// TestSortBy Sorts the slice in ascending order as determined by the less function. This sort is not guaranteed to be stable.
// func SortBy[T any](slice []T, less func(a, b T) bool)
func TestSortBy(t *testing.T) {
	type User struct {
		Name string
		Age  uint
	}

	{
		numbers := []int{5, 3, 8, 1}

		// 降序排序
		sort.Slice(numbers, func(i, j int) bool {
			return numbers[i] > numbers[j]
		})
		fmt.Println("降序:", numbers) // 输出: [8 5 3 1]

		// 升序排序
		sort.Slice(numbers, func(i, j int) bool {
			return numbers[i] < numbers[j]
		})
		fmt.Println("升序:", numbers) // 输出: [1 3 5 8]
	}

	{
		numbers := []int{1, 4, 3, 2, 5}

		slice.SortBy(numbers, func(a, b int) bool {
			return a < b
		})
		assert.Equal(t, []int{1, 2, 3, 4, 5}, numbers)

		users := []User{
			{Name: "a", Age: 21},
			{Name: "b", Age: 15},
			{Name: "c", Age: 100},
		}

		slice.SortBy(users, func(a, b User) bool {
			return a.Age < b.Age
		})

		assert.Equal(t, []User{
			{Name: "b", Age: 15},
			{Name: "a", Age: 21},
			{Name: "c", Age: 100},
		}, users)
	}

	{
		users := []User{
			{Name: "a", Age: 21},
			{Name: "b", Age: 15},
			{Name: "c", Age: 100},
		}

		// Sort by age, keeping original order or equal elements.
		sort.SliceStable(users, func(i, j int) bool {
			return users[i].Age < users[j].Age
		})
		fmt.Println(users) // [{b 15} {a 21} {c 100}]
	}
}

// TestSortByField Sort struct slice by field.
// Slice element should be struct, `field` param type should be int, uint, string, or bool.
// Default sort type is ascending (asc), if descending order, set `sortType` param to desc.
// func SortByField(slice any, field string, sortType ...string) error
func TestSortByField(t *testing.T) {
	type User struct {
		Name string
		Age  uint
	}

	users := []User{
		{Name: "a", Age: 21},
		{Name: "b", Age: 15},
		{Name: "c", Age: 100},
	}

	err := slice.SortByField(users, "Age", "desc")
	if err != nil {
		return
	}

	assert.Equal(t, []User{
		{Name: "c", Age: 100},
		{Name: "a", Age: 21},
		{Name: "b", Age: 15},
	}, users)
}

type Person struct {
	Name string
	Age  int
}

// ByAge implements sort.Interface based on the Age field.
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func TestSortInterface(t *testing.T) {
	// 只要某一个数据结构实现了 Len() int，Less(i, j int) bool 和 Swap(i, j int) 这三个方法，那么就可以使用 sort.Sort 来排序
	//
	//	type Interface interface {
	//	    // Len is the number of elements in the collection.
	//	    Len() int
	//	    // Less reports whether the element with
	//	    // index i should sort before the element with index j.
	//	    Less(i, j int) bool
	//	    // Swap swaps the elements with indexes i and j.
	//	    Swap(i, j int)
	//	}

	family := []Person{
		{"Alice", 23},
		{"Eve", 2},
		{"Bob", 25},
	}
	sort.Sort(ByAge(family))
	fmt.Println(family) // [{Eve 2} {Alice 23} {Bob 25}]
}
