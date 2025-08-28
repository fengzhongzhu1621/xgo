package datetime

import "time"

func Duration(str string) (time.Duration, error) {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return time.Duration(0), err
	}
	return dur, nil
}

// Milliseconds 将 time.Duration 类型转换为毫秒（float64）的函数
func Milliseconds(t time.Duration) float64 {
	// t/1e6 将纳秒值除以 1,000,000（即 1e6），得到毫秒的整数部分
	// 3.5ms 的纳秒表示为 3,500,000ns，3,500,000/1e6 = 3（整数部分）
	//
	// t%1e6：取纳秒值除以 1e6 的余数，得到不足 1ms 的纳秒部分
	// 例如 3,500,000ns % 1e6 = 500,000ns。
	// float64(t%1e6)/1e6：将余数转换为浮点数并除以 1e6，得到毫秒的小数部分
	// 500,000 / 1e6 = 0.5，最终结果为 3 + 0.5 = 3.5ms。
	return float64(t/1e6) + float64(t%1e6)/1e6
}
