package mask

// MaskRealName 用户名脱敏
func MaskRealName(realName string) string {
	// 支持中文
	runeRealNameSlice := []rune(realName)
	if n := len(runeRealNameSlice); n >= 2 {
		if n == 2 {
			// AB -> A*
			return string(append(runeRealNameSlice[0:1], rune('*')))
		}

		// 第一个字符不隐藏
		newRealName := runeRealNameSlice[0:1]
		// 计算需要隐藏的字符的数量，最后一个字符不隐藏，中间字符全部隐藏
		count := n - 2
		for temp := 1; temp <= count; temp++ {
			newRealName = append(newRealName, rune('*'))
		}
		return string(append(newRealName, runeRealNameSlice[n-1]))

	}

	return realName
}
