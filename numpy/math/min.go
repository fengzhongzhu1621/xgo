package math

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}
