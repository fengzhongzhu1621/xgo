// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bytesconv

import (
	"math/rand"
	"reflect"
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

/////////////////////////////////////////////////////////////////////////////
// 字符串转换为[]bytes
// StringToBytes converts string to byte slice without a memory allocation.
// 效率更高
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// Bytes converts string_utils to byte slice.
func Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func rawStrToBytes(s string) []byte {
	return []byte(s)
}

func SafeBytes(s string) []byte {
	return []byte(s)
}

/////////////////////////////////////////////////////////////////////////////
// []bytes转换为字符串
// BytesToString converts byte slice to string without a memory allocation.
// 效率更高
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String converts byte slice to string_utils.
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

func Atoi(b []byte) (int, error) {
	return strconv.Atoi(BytesToString(b))
}

func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(BytesToString(b), base, bitSize)
}

func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(BytesToString(b), base, bitSize)
}

func ParseFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(BytesToString(b), bitSize)
}
