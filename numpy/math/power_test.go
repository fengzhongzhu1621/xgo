package math

import "testing"

func TestPowerToRoundUp(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		// 正常情况测试 2^power >= n
		{"n=1 should return 0", 1, 0},         // 2^0 = 1 >= 1
		{"n=2 should return 1", 2, 1},         // 2^1 = 2 >= 2
		{"n=3 should return 2", 3, 2},         // 2^2 = 4 >= 3
		{"n=4 should return 2", 4, 2},         // 2^2 = 4 >= 4
		{"n=5 should return 3", 5, 3},         // 2^3 = 8 >= 5
		{"n=8 should return 3", 8, 3},         // 2^3 = 8 >= 8
		{"n=9 should return 4", 9, 4},         // 2^4 = 16 >= 9
		{"n=16 should return 4", 16, 4},       // 2^4 = 16 >= 16
		{"n=17 should return 5", 17, 5},       // 2^5 = 32 >= 17
		{"n=1024 should return 10", 1024, 10}, // 2^10 = 1024 >= 1024

		// 边界情况测试
		{"n=0 should return 0", 0, 0},       // 特殊情况处理
		{"n=-1 should return 0", -1, 0},     // 负数处理
		{"n=-100 should return 0", -100, 0}, // 大负数处理
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PowerToRoundUp(tt.input); got != tt.expected {
				t.Errorf("powerToRoundUp(%d) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func BenchmarkPowerToRoundUp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PowerToRoundUp(1000)
	}
}
