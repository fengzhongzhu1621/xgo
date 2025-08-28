package slice

// Chunk 将整型数组按指定长度分组
func Chunk(slice []int64, chunkSize int) [][]int64 {
	var chunks [][]int64
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
