package cast

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

func TestStringToBytes(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := RandStringBytesMaskImprSrcSB(64)
		if !bytes.Equal(rawStrToBytes(s), StringToBytes(s)) {
			t.Fatal("don't match")
		}
	}
}

func BenchmarkBytesConvStrToBytesRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawStrToBytes(testString)
	}
}

func BenchmarkBytesConvStrToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringToBytes(testString)
	}
}

// TestLancetStringToBytes 将字符串转换为字节切片而无需进行内存分配。
// func StringToBytes(str string) (b []byte)
func TestLancetStringToBytes(t *testing.T) {
	result1 := strutil.StringToBytes("abc")
	result2 := reflect.DeepEqual(result1, []byte{'a', 'b', 'c'})

	// [97 98 99]
	fmt.Println(result1)
	assert.Equal(t, true, result2)
}
