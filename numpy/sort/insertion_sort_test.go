package sort

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/algorithm"
)

type people struct {
	Name string
	Age  int
}

// PeopleAageComparator sort people slice by age field
type peopleAgeComparator struct{}

// Compare implements github.com/duke-git/lancet/constraints/constraints.go/Comparator
func (pc *peopleAgeComparator) Compare(v1 any, v2 any) int {
	p1, _ := v1.(people)
	p2, _ := v2.(people)

	// ascending order
	if p1.Age < p2.Age {
		return -1
	} else if p1.Age > p2.Age {
		return 1
	}

	return 0
}

// TestInsertionSort 使用插入排序算法对切片进行排序。参数comparator应实现constraints.Comparator。
// func InsertionSort[T any](slice []T, comparator constraints.Comparator)
func TestInsertionSort(t *testing.T) {
	peoples := []people{
		{Name: "a", Age: 20},
		{Name: "b", Age: 10},
		{Name: "c", Age: 17},
		{Name: "d", Age: 8},
		{Name: "e", Age: 28},
	}

	comparator := &peopleAgeComparator{}

	algorithm.InsertionSort(peoples, comparator)

	fmt.Println(peoples)
}
