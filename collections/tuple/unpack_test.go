package tuple

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		tuple := lo.Tuple2[string, int]{A: "a", B: 1}

		r1, r2 := lo.Unpack2(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)

		r1, r2 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
	}

	{
		tuple := lo.Tuple3[string, int, float64]{A: "a", B: 1, C: 1.0}

		r1, r2, r3 := lo.Unpack3(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)

		r1, r2, r3 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
	}

	{
		tuple := lo.Tuple4[string, int, float64, bool]{A: "a", B: 1, C: 1.0, D: true}

		r1, r2, r3, r4 := lo.Unpack4(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)

		r1, r2, r3, r4 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
	}

	{
		tuple := lo.Tuple5[string, int, float64, bool, string]{
			A: "a",
			B: 1,
			C: 1.0,
			D: true,
			E: "b",
		}

		r1, r2, r3, r4, r5 := lo.Unpack5(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)

		r1, r2, r3, r4, r5 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
	}

	{
		tuple := lo.Tuple6[string, int, float64, bool, string, int]{
			A: "a",
			B: 1,
			C: 1.0,
			D: true,
			E: "b",
			F: 2,
		}

		r1, r2, r3, r4, r5, r6 := lo.Unpack6(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)

		r1, r2, r3, r4, r5, r6 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
	}

	{
		tuple := lo.Tuple7[string, int, float64, bool, string, int, float64]{
			A: "a",
			B: 1,
			C: 1.0,
			D: true,
			E: "b",
			F: 2,
			G: 3.0,
		}

		r1, r2, r3, r4, r5, r6, r7 := lo.Unpack7(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
		is.Equal(3.0, r7)

		r1, r2, r3, r4, r5, r6, r7 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
		is.Equal(3.0, r7)
	}

	{
		tuple := lo.Tuple8[string, int, float64, bool, string, int, float64, bool]{
			A: "a",
			B: 1,
			C: 1.0,
			D: true,
			E: "b",
			F: 2,
			G: 3.0,
			H: true,
		}

		r1, r2, r3, r4, r5, r6, r7, r8 := lo.Unpack8(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
		is.Equal(3.0, r7)
		is.Equal(true, r8)

		r1, r2, r3, r4, r5, r6, r7, r8 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
		is.Equal(3.0, r7)
		is.Equal(true, r8)
	}

	{
		tuple := lo.Tuple9[string, int, float64, bool, string, int, float64, bool, string]{
			A: "a",
			B: 1,
			C: 1.0,
			D: true,
			E: "b",
			F: 2,
			G: 3.0,
			H: true,
			I: "c",
		}

		r1, r2, r3, r4, r5, r6, r7, r8, r9 := lo.Unpack9(tuple)

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
		is.Equal(3.0, r7)
		is.Equal(true, r8)
		is.Equal("c", r9)

		r1, r2, r3, r4, r5, r6, r7, r8, r9 = tuple.Unpack()

		is.Equal("a", r1)
		is.Equal(1, r2)
		is.Equal(1.0, r3)
		is.Equal(true, r4)
		is.Equal("b", r5)
		is.Equal(2, r6)
		is.Equal(3.0, r7)
		is.Equal(true, r8)
		is.Equal("c", r9)
	}
}
