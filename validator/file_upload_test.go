package validator

import "testing"

func TestIsPathValid(t *testing.T) {
	storagePath := "/safe/storage"

	tests := []struct {
		path     string
		expected bool
	}{
		{"/safe/storage/file.txt", true},        // 合法路径
		{"/safe/storage/../file.txt", false},    // 路径穿越
		{"/safe/storage/../../file.txt", false}, // 路径穿越
		{"/safe/storage/COM1", false},           // Windows 保留文件名
		{"/safe/storage/.", false},              // Linux 保留文件名
		{"/safe/storage/file?.txt", false},      // 非法字符
		{"/safe/storage/file*.txt", false},      // 非法字符
		{"/safe/storage/file\u0000.txt", false}, // 非法 Unicode
		{"/safe/storage/./file.txt", true},      // 规范化后合法
		{"/safe/storage/file.txt/..", false},    // 路径穿越
	}

	for _, tt := range tests {
		if got := IsPathValid(storagePath, tt.path); got != tt.expected {
			t.Errorf("isPathValid(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}
