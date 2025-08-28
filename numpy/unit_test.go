package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/formatter"
	"github.com/dustin/go-humanize"
	"github.com/gookit/goutil/fmtutil"
	"github.com/gookit/goutil/mathutil"
	"github.com/stretchr/testify/assert"
)

// 返回以十进制（基数为1000）标准下人类可读的字节大小。精度参数指定小数点后的位数，默认为4位。
// func DecimalBytes(size float64, precision ...int) string
func TestDecimalBytes(t *testing.T) {
	result1 := formatter.DecimalBytes(1000)
	result2 := formatter.DecimalBytes(1024)
	result3 := formatter.DecimalBytes(1234567)
	result4 := formatter.DecimalBytes(1234567, 3)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// 1KB
	// 1.024KB
	// 1.2346MB
	// 1.235MB
}

// 返回以二进制（基数为1024）标准下人类可读的字节大小。精度参数指定小数点后的位数，默认为4位。
// func BinaryBytes(size float64, precision ...int) string
func TestBinaryBytes(t *testing.T) {
	{
		result1 := formatter.BinaryBytes(1024)
		result2 := formatter.BinaryBytes(1024 * 1024)
		result3 := formatter.BinaryBytes(1234567)
		result4 := formatter.BinaryBytes(1234567, 2)

		fmt.Println(result1) // 1KiB
		fmt.Println(result2) // 1MiB
		fmt.Println(result3) // 1.1774MiB
		fmt.Println(result4) // 1.18MiB
	}

	{
		result1 := humanize.Bytes(1024)
		result2 := humanize.Bytes(1024 * 1024)
		result3 := humanize.Bytes(1234567)
		result4 := humanize.Bytes(1234567)

		fmt.Println(result1) // 1.0 kB
		fmt.Println(result2) // 1.0 MB
		fmt.Println(result3) // 1.2 MB
		fmt.Println(result4) // 1.2 MB
	}

	{
		tests := []struct {
			args uint64
			want string
		}{
			{346, "346B"},
			{3467, "3.39K"},
			{346778, "338.65K"},
			{12346778, "11.77M"},
			{1200346778, "1.12G"},
		}

		for _, tt := range tests {
			assert.Equal(t, tt.want, mathutil.DataSize(tt.args))
			// SizeToString 是 DataSize 的别名
			assert.Equal(t, tt.want, fmtutil.SizeToString(tt.args))
		}
	}
}

// func ParseDecimalBytes(size string) (uint64, error)
func TestParseDecimalBytes(t *testing.T) {
	result1, _ := formatter.ParseDecimalBytes("12")
	result2, _ := formatter.ParseDecimalBytes("12k")
	result3, _ := formatter.ParseDecimalBytes("12 Kb")
	result4, _ := formatter.ParseDecimalBytes("12.2 kb")

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// 12
	// 12000
	// 12000
	// 12200
}

// func ParseBinaryBytes(size string) (uint64, error)
func TestParseBinaryBytes(t *testing.T) {
	{
		result1, _ := formatter.ParseBinaryBytes("12")
		result2, _ := formatter.ParseBinaryBytes("12ki")
		result3, _ := formatter.ParseBinaryBytes("12 KiB")
		result4, _ := formatter.ParseBinaryBytes("12.2 kib")

		fmt.Println(result1)
		fmt.Println(result2)
		fmt.Println(result3)
		fmt.Println(result4)

		// Output:
		// 12
		// 12288
		// 12288
		// 12492
	}

	{
		// 将表示字节大小的人类可读字符串（如 "1GB"、"12mb" 等）转换为对应的无符号整数字节数 (uint64)。
		// 根据最后一个字符（转换为小写后）确定对应的乘数：
		// 'k' 或 'K'：1 << 10（1024）
		// 'm' 或 'M'：1 << 20（1,048,576）
		// 'g' 或 'G'：1 << 30（1,073,741,824）
		// 't' 或 'T'：1 << 40（1,099,511,627,776）
		// 'p' 或 'P'：1 << 50（1,125,899,906,842,624）
		examples := []string{"1GB", "512MB", "1024", "2TB", "invalid", "3.5PB"}
		for _, str := range examples {
			bytes := fmtutil.StringToByte(str)
			fmt.Printf("ToByteSize(%s) = %d bytes\n", str, bytes)
		}
		// ToByteSize(1GB) = 1073741824 bytes
		// ToByteSize(512MB) = 536870912 bytes
		// ToByteSize(1024) = 1024 bytes
		// ToByteSize(2TB) = 2199023255552 bytes
		// ToByteSize(invalid) = 0 bytes
		// ToByteSize(3.5PB) = 3940649673949184 bytes
	}
}
