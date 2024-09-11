package stringutils

// SliceRemoveDuplicates 从一个字符串切片（[]string）中移除重复的元素，并返回一个新的不包含重复元素的切片。
func SliceRemoveDuplicates(s []string) []string {
	m := make(map[string]bool)
	result := []string{}
	for _, v := range s {
		// 判断 key 是否存在
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}
