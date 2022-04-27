package fastrand

// 获得指定范围的随机整数 [min, max).
func RangeIntn(min, max int) int {
	return Intn(max-min) + min
}
