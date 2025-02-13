package randutils

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		rand.Seed(time.Now().UnixNano())

		result1 := lo.Sample([]string{"a", "b", "c"})
		result2 := lo.Sample([]string{})

		is.True(lo.Contains([]string{"a", "b", "c"}, result1))
		is.Equal(result2, "")
	}

	{
		intSlice := []int{1, 2, 3, 4, 5, 6}
		intVal := arrutil.GetRandomOne(intSlice)
		intVal1 := arrutil.GetRandomOne(intSlice)
		for intVal == intVal1 {
			intVal1 = arrutil.GetRandomOne(intSlice)
		}
	}

}

func TestSamples(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	rand.Seed(time.Now().UnixNano())

	result1 := lo.Samples([]string{"a", "b", "c"}, 3)
	result2 := lo.Samples([]string{}, 3)

	sort.Strings(result1)

	is.Equal(result1, []string{"a", "b", "c"})
	is.Equal(result2, []string{})

	type myStrings []string
	allStrings := myStrings{"", "foo", "bar"}
	nonempty := lo.Samples(allStrings, 2)
	is.IsType(nonempty, allStrings, "type preserved")
}
