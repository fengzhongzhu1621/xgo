package numpy

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/formatter"
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
	result1 := formatter.BinaryBytes(1024)
	result2 := formatter.BinaryBytes(1024 * 1024)
	result3 := formatter.BinaryBytes(1234567)
	result4 := formatter.BinaryBytes(1234567, 2)

	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)

	// Output:
	// 1KiB
	// 1MiB
	// 1.1774MiB
	// 1.18MiB
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
