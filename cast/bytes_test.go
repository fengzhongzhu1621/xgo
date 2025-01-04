package cast

import (
	"math/rand"
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

var testString = "Albert Einstein: Logic will get you from A to B. Imagination will take you everywhere."
var testBytes = []byte(testString)

// go test -v

func TestBytesToString(t *testing.T) {
	data := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		rand.Read(data)
		if rawBytesToStr(data) != BytesToString(data) {
			t.Fatal("don't match")
		}
	}
}

// TestLancetBytesToString 将字节切片转换为字符串而无需进行内存分配。
// func BytesToString(bytes []byte) string
func TestLancetBytesToString(t *testing.T) {
	bytes := []byte{'a', 'b', 'c'}
	result := strutil.BytesToString(bytes)
	assert.Equal(t, "abc", result)
}

// go test -v -run=none -bench=^BenchmarkBytesConv -benchmem=true

func BenchmarkBytesConvBytesToStrRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawBytesToStr(testBytes)
	}
}

func BenchmarkBytesConvBytesToStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesToString(testBytes)
	}
}
