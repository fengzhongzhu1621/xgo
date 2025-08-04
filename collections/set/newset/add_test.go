package newset

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

const N = 1000

func Test_AddConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := mapset.NewSet[int]()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	for _, i := range ints {
		if !s.Contains(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}

func TestMakeSetInt(*testing.T) {
	ints := []int{1, 2, 3}
	s := mapset.NewSet[int]()
	for _, i := range ints {
		s.Add(i)
	}
	fmt.Println(s) // Set{2, 3, 1}
}
