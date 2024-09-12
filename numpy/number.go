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

// Chunk 将整型数组按指定长度分组
func Chunk(slice []int64, chunkSize int) [][]int64 {
	var (
		chunks [][]int64
	)
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
