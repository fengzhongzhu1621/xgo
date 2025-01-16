package tuple

import (
	"testing"

	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"
)

func TestZip(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.Zip2(
		[]string{"a", "b"},
		[]int{1, 2},
	)

	r2 := lo.Zip3(
		[]string{"a", "b", "c"},
		[]int{1, 2, 3},
		[]int{4, 5, 6},
	)

	r3 := lo.Zip4(
		[]string{"a", "b", "c", "d"},
		[]int{1, 2, 3, 4},
		[]int{5, 6, 7, 8},
		[]bool{true, true, true, true},
	)

	r4 := lo.Zip5(
		[]string{"a", "b", "c", "d", "e"},
		[]int{1, 2, 3, 4, 5},
		[]int{6, 7, 8, 9, 10},
		[]bool{true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5},
	)

	r5 := lo.Zip6(
		[]string{"a", "b", "c", "d", "e", "f"},
		[]int{1, 2, 3, 4, 5, 6},
		[]int{7, 8, 9, 10, 11, 12},
		[]bool{true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06},
	)

	r6 := lo.Zip7(
		[]string{"a", "b", "c", "d", "e", "f", "g"},
		[]int{1, 2, 3, 4, 5, 6, 7},
		[]int{8, 9, 10, 11, 12, 13, 14},
		[]bool{true, true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07},
		[]int8{1, 2, 3, 4, 5, 6, 7},
	)

	r7 := lo.Zip8(
		[]string{"a", "b", "c", "d", "e", "f", "g", "h"},
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		[]int{9, 10, 11, 12, 13, 14, 15, 16},
		[]bool{true, true, true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08},
		[]int8{1, 2, 3, 4, 5, 6, 7, 8},
		[]int16{1, 2, 3, 4, 5, 6, 7, 8},
	)

	r8 := lo.Zip9(
		[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int{10, 11, 12, 13, 14, 15, 16, 17, 18},
		[]bool{true, true, true, true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09},
		[]int8{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int16{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int32{1, 2, 3, 4, 5, 6, 7, 8, 9},
	)

	is.Equal(r1, []lo.Tuple2[string, int]{
		{A: "a", B: 1},
		{A: "b", B: 2},
	})

	is.Equal(r2, []lo.Tuple3[string, int, int]{
		{A: "a", B: 1, C: 4},
		{A: "b", B: 2, C: 5},
		{A: "c", B: 3, C: 6},
	})

	is.Equal(r3, []lo.Tuple4[string, int, int, bool]{
		{A: "a", B: 1, C: 5, D: true},
		{A: "b", B: 2, C: 6, D: true},
		{A: "c", B: 3, C: 7, D: true},
		{A: "d", B: 4, C: 8, D: true},
	})

	is.Equal(r4, []lo.Tuple5[string, int, int, bool, float32]{
		{A: "a", B: 1, C: 6, D: true, E: 0.1},
		{A: "b", B: 2, C: 7, D: true, E: 0.2},
		{A: "c", B: 3, C: 8, D: true, E: 0.3},
		{A: "d", B: 4, C: 9, D: true, E: 0.4},
		{A: "e", B: 5, C: 10, D: true, E: 0.5},
	})

	is.Equal(r5, []lo.Tuple6[string, int, int, bool, float32, float64]{
		{A: "a", B: 1, C: 7, D: true, E: 0.1, F: 0.01},
		{A: "b", B: 2, C: 8, D: true, E: 0.2, F: 0.02},
		{A: "c", B: 3, C: 9, D: true, E: 0.3, F: 0.03},
		{A: "d", B: 4, C: 10, D: true, E: 0.4, F: 0.04},
		{A: "e", B: 5, C: 11, D: true, E: 0.5, F: 0.05},
		{A: "f", B: 6, C: 12, D: true, E: 0.6, F: 0.06},
	})

	is.Equal(r6, []lo.Tuple7[string, int, int, bool, float32, float64, int8]{
		{A: "a", B: 1, C: 8, D: true, E: 0.1, F: 0.01, G: 1},
		{A: "b", B: 2, C: 9, D: true, E: 0.2, F: 0.02, G: 2},
		{A: "c", B: 3, C: 10, D: true, E: 0.3, F: 0.03, G: 3},
		{A: "d", B: 4, C: 11, D: true, E: 0.4, F: 0.04, G: 4},
		{A: "e", B: 5, C: 12, D: true, E: 0.5, F: 0.05, G: 5},
		{A: "f", B: 6, C: 13, D: true, E: 0.6, F: 0.06, G: 6},
		{A: "g", B: 7, C: 14, D: true, E: 0.7, F: 0.07, G: 7},
	})

	is.Equal(r7, []lo.Tuple8[string, int, int, bool, float32, float64, int8, int16]{
		{A: "a", B: 1, C: 9, D: true, E: 0.1, F: 0.01, G: 1, H: 1},
		{A: "b", B: 2, C: 10, D: true, E: 0.2, F: 0.02, G: 2, H: 2},
		{A: "c", B: 3, C: 11, D: true, E: 0.3, F: 0.03, G: 3, H: 3},
		{A: "d", B: 4, C: 12, D: true, E: 0.4, F: 0.04, G: 4, H: 4},
		{A: "e", B: 5, C: 13, D: true, E: 0.5, F: 0.05, G: 5, H: 5},
		{A: "f", B: 6, C: 14, D: true, E: 0.6, F: 0.06, G: 6, H: 6},
		{A: "g", B: 7, C: 15, D: true, E: 0.7, F: 0.07, G: 7, H: 7},
		{A: "h", B: 8, C: 16, D: true, E: 0.8, F: 0.08, G: 8, H: 8},
	})

	is.Equal(r8, []lo.Tuple9[string, int, int, bool, float32, float64, int8, int16, int32]{
		{A: "a", B: 1, C: 10, D: true, E: 0.1, F: 0.01, G: 1, H: 1, I: 1},
		{A: "b", B: 2, C: 11, D: true, E: 0.2, F: 0.02, G: 2, H: 2, I: 2},
		{A: "c", B: 3, C: 12, D: true, E: 0.3, F: 0.03, G: 3, H: 3, I: 3},
		{A: "d", B: 4, C: 13, D: true, E: 0.4, F: 0.04, G: 4, H: 4, I: 4},
		{A: "e", B: 5, C: 14, D: true, E: 0.5, F: 0.05, G: 5, H: 5, I: 5},
		{A: "f", B: 6, C: 15, D: true, E: 0.6, F: 0.06, G: 6, H: 6, I: 6},
		{A: "g", B: 7, C: 16, D: true, E: 0.7, F: 0.07, G: 7, H: 7, I: 7},
		{A: "h", B: 8, C: 17, D: true, E: 0.8, F: 0.08, G: 8, H: 8, I: 8},
		{A: "i", B: 9, C: 18, D: true, E: 0.9, F: 0.09, G: 9, H: 9, I: 9},
	})
}

func TestZipBy(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.ZipBy2(
		[]string{"a", "b"},
		[]int{1, 2},
		func(a string, b int) lo.Tuple2[string, int] {
			return lo.T2(a, b)
		},
	)

	r2 := lo.ZipBy3(
		[]string{"a", "b", "c"},
		[]int{1, 2, 3},
		[]int{4, 5, 6},
		func(a string, b int, c int) lo.Tuple3[string, int, int] {
			return lo.T3(a, b, c)
		},
	)

	r3 := lo.ZipBy4(
		[]string{"a", "b", "c", "d"},
		[]int{1, 2, 3, 4},
		[]int{5, 6, 7, 8},
		[]bool{true, true, true, true},
		func(a string, b int, c int, d bool) lo.Tuple4[string, int, int, bool] {
			return lo.T4(a, b, c, d)
		},
	)

	r4 := lo.ZipBy5(
		[]string{"a", "b", "c", "d", "e"},
		[]int{1, 2, 3, 4, 5},
		[]int{6, 7, 8, 9, 10},
		[]bool{true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5},
		func(a string, b int, c int, d bool, e float32) lo.Tuple5[string, int, int, bool, float32] {
			return lo.T5(a, b, c, d, e)
		},
	)

	r5 := lo.ZipBy6(
		[]string{"a", "b", "c", "d", "e", "f"},
		[]int{1, 2, 3, 4, 5, 6},
		[]int{7, 8, 9, 10, 11, 12},
		[]bool{true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06},
		func(a string, b int, c int, d bool, e float32, f float64) lo.Tuple6[string, int, int, bool, float32, float64] {
			return lo.T6(a, b, c, d, e, f)
		},
	)

	r6 := lo.ZipBy7(
		[]string{"a", "b", "c", "d", "e", "f", "g"},
		[]int{1, 2, 3, 4, 5, 6, 7},
		[]int{8, 9, 10, 11, 12, 13, 14},
		[]bool{true, true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07},
		[]int8{1, 2, 3, 4, 5, 6, 7},
		func(a string, b int, c int, d bool, e float32, f float64, g int8) lo.Tuple7[string, int, int, bool, float32, float64, int8] {
			return lo.T7(a, b, c, d, e, f, g)
		},
	)

	r7 := lo.ZipBy8(
		[]string{"a", "b", "c", "d", "e", "f", "g", "h"},
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		[]int{9, 10, 11, 12, 13, 14, 15, 16},
		[]bool{true, true, true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08},
		[]int8{1, 2, 3, 4, 5, 6, 7, 8},
		[]int16{1, 2, 3, 4, 5, 6, 7, 8},
		func(a string, b int, c int, d bool, e float32, f float64, g int8, h int16) lo.Tuple8[string, int, int, bool, float32, float64, int8, int16] {
			return lo.T8(a, b, c, d, e, f, g, h)
		},
	)

	r8 := lo.ZipBy9(
		[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int{10, 11, 12, 13, 14, 15, 16, 17, 18},
		[]bool{true, true, true, true, true, true, true, true, true},
		[]float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
		[]float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09},
		[]int8{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int16{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int32{1, 2, 3, 4, 5, 6, 7, 8, 9},
		func(a string, b int, c int, d bool, e float32, f float64, g int8, h int16, i int32) lo.Tuple9[string, int, int, bool, float32, float64, int8, int16, int32] {
			return lo.T9(a, b, c, d, e, f, g, h, i)
		},
	)

	is.Equal(r1, []lo.Tuple2[string, int]{
		{A: "a", B: 1},
		{A: "b", B: 2},
	})

	is.Equal(r2, []lo.Tuple3[string, int, int]{
		{A: "a", B: 1, C: 4},
		{A: "b", B: 2, C: 5},
		{A: "c", B: 3, C: 6},
	})

	is.Equal(r3, []lo.Tuple4[string, int, int, bool]{
		{A: "a", B: 1, C: 5, D: true},
		{A: "b", B: 2, C: 6, D: true},
		{A: "c", B: 3, C: 7, D: true},
		{A: "d", B: 4, C: 8, D: true},
	})

	is.Equal(r4, []lo.Tuple5[string, int, int, bool, float32]{
		{A: "a", B: 1, C: 6, D: true, E: 0.1},
		{A: "b", B: 2, C: 7, D: true, E: 0.2},
		{A: "c", B: 3, C: 8, D: true, E: 0.3},
		{A: "d", B: 4, C: 9, D: true, E: 0.4},
		{A: "e", B: 5, C: 10, D: true, E: 0.5},
	})

	is.Equal(r5, []lo.Tuple6[string, int, int, bool, float32, float64]{
		{A: "a", B: 1, C: 7, D: true, E: 0.1, F: 0.01},
		{A: "b", B: 2, C: 8, D: true, E: 0.2, F: 0.02},
		{A: "c", B: 3, C: 9, D: true, E: 0.3, F: 0.03},
		{A: "d", B: 4, C: 10, D: true, E: 0.4, F: 0.04},
		{A: "e", B: 5, C: 11, D: true, E: 0.5, F: 0.05},
		{A: "f", B: 6, C: 12, D: true, E: 0.6, F: 0.06},
	})

	is.Equal(r6, []lo.Tuple7[string, int, int, bool, float32, float64, int8]{
		{A: "a", B: 1, C: 8, D: true, E: 0.1, F: 0.01, G: 1},
		{A: "b", B: 2, C: 9, D: true, E: 0.2, F: 0.02, G: 2},
		{A: "c", B: 3, C: 10, D: true, E: 0.3, F: 0.03, G: 3},
		{A: "d", B: 4, C: 11, D: true, E: 0.4, F: 0.04, G: 4},
		{A: "e", B: 5, C: 12, D: true, E: 0.5, F: 0.05, G: 5},
		{A: "f", B: 6, C: 13, D: true, E: 0.6, F: 0.06, G: 6},
		{A: "g", B: 7, C: 14, D: true, E: 0.7, F: 0.07, G: 7},
	})

	is.Equal(r7, []lo.Tuple8[string, int, int, bool, float32, float64, int8, int16]{
		{A: "a", B: 1, C: 9, D: true, E: 0.1, F: 0.01, G: 1, H: 1},
		{A: "b", B: 2, C: 10, D: true, E: 0.2, F: 0.02, G: 2, H: 2},
		{A: "c", B: 3, C: 11, D: true, E: 0.3, F: 0.03, G: 3, H: 3},
		{A: "d", B: 4, C: 12, D: true, E: 0.4, F: 0.04, G: 4, H: 4},
		{A: "e", B: 5, C: 13, D: true, E: 0.5, F: 0.05, G: 5, H: 5},
		{A: "f", B: 6, C: 14, D: true, E: 0.6, F: 0.06, G: 6, H: 6},
		{A: "g", B: 7, C: 15, D: true, E: 0.7, F: 0.07, G: 7, H: 7},
		{A: "h", B: 8, C: 16, D: true, E: 0.8, F: 0.08, G: 8, H: 8},
	})

	is.Equal(r8, []lo.Tuple9[string, int, int, bool, float32, float64, int8, int16, int32]{
		{A: "a", B: 1, C: 10, D: true, E: 0.1, F: 0.01, G: 1, H: 1, I: 1},
		{A: "b", B: 2, C: 11, D: true, E: 0.2, F: 0.02, G: 2, H: 2, I: 2},
		{A: "c", B: 3, C: 12, D: true, E: 0.3, F: 0.03, G: 3, H: 3, I: 3},
		{A: "d", B: 4, C: 13, D: true, E: 0.4, F: 0.04, G: 4, H: 4, I: 4},
		{A: "e", B: 5, C: 14, D: true, E: 0.5, F: 0.05, G: 5, H: 5, I: 5},
		{A: "f", B: 6, C: 15, D: true, E: 0.6, F: 0.06, G: 6, H: 6, I: 6},
		{A: "g", B: 7, C: 16, D: true, E: 0.7, F: 0.07, G: 7, H: 7, I: 7},
		{A: "h", B: 8, C: 17, D: true, E: 0.8, F: 0.08, G: 8, H: 8, I: 8},
		{A: "i", B: 9, C: 18, D: true, E: 0.9, F: 0.09, G: 9, H: 9, I: 9},
	})
}
