package mask

// MaskPhone 手机号脱敏
func MaskPhone(phone string) string {
	if n := len(phone); n >= 8 {
		// 隐去 4 位地区码
		return phone[:n-8] + "****" + phone[n-4:]
	}

	return phone
}
