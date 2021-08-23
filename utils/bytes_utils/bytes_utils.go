package bytes_utils

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
	"xgo/utils/bytesconv"
)

/**
 * 判断前缀和后缀是否全部匹配
 */
func HasPrefixAndSuffix(s, prefix []byte, suffix []byte) bool {
	return bytes.HasPrefix(s, prefix) && bytes.HasSuffix(s, suffix)
}

/**
 * 去掉结尾换行符
 * trimEOL cuts unixy style \n and windowsy style \r\n suffix from the string
 */
func TrimEOL(b []byte) []byte {
	lns := len(b)
	if lns > 0 && b[lns-1] == '\n' {
		lns--
		if lns > 0 && b[lns-1] == '\r' {
			lns--
		}
	}
	return b[:lns]
}

////////////////////////////////////////////////////////////////////////////////////
// 字符数组拼接
func AppendArg(b []byte, v interface{}) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		return AppendUTF8String(b, bytesconv.Bytes(v))
	case []byte:
		return AppendUTF8String(b, v)
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	}
}

func AppendUTF8String(dst []byte, src []byte) []byte {
	dst = append(dst, src...)
	return dst
}


//////////////////////////////////////////////////////////////////////////////////////
// 将字节数组转换为其它类型

func Atoi(b []byte) (int, error) {
	return strconv.Atoi(bytesconv.BytesToString(b))
}

func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(bytesconv.BytesToString(b), base, bitSize)
}

func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(bytesconv.BytesToString(b), base, bitSize)
}

func ParseFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(bytesconv.BytesToString(b), bitSize)
}
