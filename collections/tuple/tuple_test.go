package tuple

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestT(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	r1 := lo.T2("a", 1)
	r2 := lo.T3[string, int, float32]("b", 2, 3.0)
	r3 := lo.T4[string, int, float32]("c", 3, 4.0, true)
	r4 := lo.T5[string, int, float32]("d", 4, 5.0, false, "e")
	r5 := lo.T6[string, int, float32]("f", 5, 6.0, true, "g", 7)
	r6 := lo.T7[string, int, float32]("h", 6, 7.0, false, "i", 8, 9.0)
	r7 := lo.T8[string, int, float32]("j", 7, 8.0, true, "k", 9, 10.0, false)
	r8 := lo.T9[string, int, float32]("l", 8, 9.0, false, "m", 10, 11.0, true, "n")

	is.Equal(r1, lo.Tuple2[string, int]{A: "a", B: 1})
	is.Equal(r2, lo.Tuple3[string, int, float32]{A: "b", B: 2, C: 3.0})
	is.Equal(r3, lo.Tuple4[string, int, float32, bool]{A: "c", B: 3, C: 4.0, D: true})
	is.Equal(
		r4,
		lo.Tuple5[string, int, float32, bool, string]{A: "d", B: 4, C: 5.0, D: false, E: "e"},
	)
	is.Equal(
		r5,
		lo.Tuple6[string, int, float32, bool, string, int]{
			A: "f",
			B: 5,
			C: 6.0,
			D: true,
			E: "g",
			F: 7,
		},
	)
	is.Equal(
		r6,
		lo.Tuple7[string, int, float32, bool, string, int, float64]{
			A: "h",
			B: 6,
			C: 7.0,
			D: false,
			E: "i",
			F: 8,
			G: 9.0,
		},
	)
	is.Equal(
		r7,
		lo.Tuple8[string, int, float32, bool, string, int, float64, bool]{
			A: "j",
			B: 7,
			C: 8.0,
			D: true,
			E: "k",
			F: 9,
			G: 10.0,
			H: false,
		},
	)
	is.Equal(
		r8,
		lo.Tuple9[string, int, float32, bool, string, int, float64, bool, string]{
			A: "l",
			B: 8,
			C: 9.0,
			D: false,
			E: "m",
			F: 10,
			G: 11.0,
			H: true,
			I: "n",
		},
	)
}
