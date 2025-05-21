package math

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func MaxInt64(a, b int64) int64 {
	if a < b {
		return b
	} else {
		return a
	}
}
