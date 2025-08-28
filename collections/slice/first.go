package slice

// FirstNotEmptyString return the first string in slice strs that is not empty
func FirstNotEmptyString(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}
