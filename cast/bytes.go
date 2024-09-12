package cast

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// []bytes转换为字符串
// BytesToString converts byte slice to string without a memory allocation.
// 效率更高.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// []bytes转换为字符串
// BytesToString converts byte slice to string without a memory allocation.
// 效率更高.
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func rawBytesToStr(b []byte) string {
	return string(b)
}

func SafeString(b []byte) string {
	return string(b)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func Itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// 将字节数组转换为int类型.
func Atoi(b []byte) (int, error) {
	return strconv.Atoi(BytesToString(b))
}

// 将字节数组转换为int64.
func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(BytesToString(b), base, bitSize)
}

// 将字节数组转换为uint64.
func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(BytesToString(b), base, bitSize)
}

// 将字节数组转换为float64.
func ParseFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(BytesToString(b), bitSize)
}
