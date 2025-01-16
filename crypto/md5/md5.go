package md5

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/pkg/errors"
)

// GetMD5Hash 生成 md5 字符串
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GetBytesMD5Hash 生成 md5 字符串
func GetBytesMD5Hash(text []byte) string {
	hash := md5.Sum(text)
	return hex.EncodeToString(hash[:])
}

// GetSliceMD5Hash 生成 md5 字符串
func GetSliceMD5Hash(old interface{}) (string, error) {
	switch trans := old.(type) {
	case []string:
		return GetMD5Hash(strings.Join(trans, ",")), nil
	case []int64:
		// 将整型数组转换为字符串后合并为一个字符串
		return GetMD5Hash(cast.ArrayInt64ToString(trans, ",")), nil
	default:
		return "", errors.New("illegal type")
	}
}
