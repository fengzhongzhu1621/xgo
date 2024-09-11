package numpy

// RemoveDuplicatesInt64 去重整型数组
func RemoveDuplicatesInt64(numbers []int64) []int64 {
	uniqMap := make(map[int64]struct{})
	result := []int64{}

	for _, v := range numbers {
		if v >= 0 {
			if _, exists := uniqMap[v]; exists {
				continue
			}
			uniqMap[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
