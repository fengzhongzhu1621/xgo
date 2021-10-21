package string_utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"xgo/utils/bytesconv"
)

/*
 * 字符替换器
 */
var HeaderNewlineToSpace = strings.NewReplacer("\n", " ", "\r", " ")
var HeaderDashToUnderscore = strings.NewReplacer("-", "_")

// 根据分隔符分割字符串
func Head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

////////////////////////////////////////////////////////
// 字符串拼接

// Deprecated: 低效的字符串拼接
func StringPlus(p []string) string {
	var s string
	l := len(p)
	for i := 0; i < l; i++ {
		s += p[i]
	}
	return s
}

// Deprecated: 低效的字符串拼接
func StringFmt(p []interface{}) string {
	return fmt.Sprint(p...)
}

// Deprecated: 低效的字符串拼接
func StringBuffer(p []string) string {
	var b bytes.Buffer
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func StringJoin(p []string) string {
	return strings.Join(p, "")
}

func StringBuilder(p []string) string {
	var b strings.Builder
	l := len(p)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

func StringBuilderEx(p []string, cap int) string {
	var b strings.Builder
	l := len(p)
	// 实现分配足够的内容，减少运行时的内存分配
	b.Grow(cap)
	for i := 0; i < l; i++ {
		b.WriteString(p[i])
	}
	return b.String()
}

// 切片去重
// 空struct不占内存空间，使用它来实现我们的函数空间复杂度是最低的。
func RemoveDuplicateElement(items []string) []string {
	result := make([]string, 0, len(items))
	// 定义集合
	set := map[string]struct{}{}
	for _, item := range items {
		if _, ok := set[item]; !ok {
			// 如果集合中不存在此元素，则加入集合
			set[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// 去掉字符串的后缀
func TrimRight(str string, substring string) string {
	idx := strings.LastIndex(str, substring)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

// 去掉字符串的前缀
func TrimLeft(str string, substring string) string {
	return strings.TrimPrefix(str, substring)
}

func UnicodeTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

func UnicodeUnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

// 返回数组最后一个元素
func Last(list []string) string {
	return list[len(list)-1]
}

// Deprecated: 切片比较
func CompareStringSliceReflect(a, b []string) bool {
	return reflect.DeepEqual(a, b)
}

// 切片比较
func CompareStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	// Golang提供BCE特性，即Bounds-checking elimination
	// 通过b = b[:len(a)]处的bounds check能够明确保证v != b[i]中的b[i]不会出现越界错误，从而避免了b[i]中的越界检查从而提高效率
	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// Deprecated: 翻转切片 panic if s is not a slice
func ReflectReverseSlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Deprecated: 翻转切片，返回一个新的切片，有copy的耗损
func ReverseSliceGetNew(s []string) []string {
	a := make([]string, len(s))
	copy(a, s)

	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}

	return a
}

// 翻转切片，值会改变，性能最高
func ReverseSlice(a []string) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

// 获得随机字符串
func GenerateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}


func Str2map(s string, sep1 string, sep2 string) map[string]string {
	if "" == s {
		return nil
	}
	spe1List := strings.Split(s, sep1)
	if len(spe1List) <= 0 {
		return nil
	}
	m := make(map[string]string)
	for _, sub := range spe1List {
		spe2List := strings.SplitN(sub, sep2, 2)
		num := len(spe2List)
		if num == 1 {
			m[spe2List[0]] = ""
		} else if num > 1 {
			m[spe2List[0]] = spe2List[1]
		}
	}
	return m
}

func MergeGetAndPostParamWithKey(queryParam map[string]string, postParam map[string]string, key string, keyName string) string {
	m := make(map[string]string)
	if len(queryParam) > 0 {
		for k, v := range queryParam {
			m[k] = v
		}
	}
	if len(postParam) > 0 {
		for k, v := range postParam {
			m[k] = v
		}
	}

	// 获取数组的key，排序
	keyList := make([]string, 0, len(m))
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)

	// 排序后的数组
	params := ""
	for _, key := range keyList {
		value, _ := m[key]
		if value != "" {
			params += key + "=" + value + "&"
		}
	}
	// 添加key参数
	params += keyName + "=" + key
	return params
}

// 计算字符串的MD5值
// 同 echo -n "123456789" | md5sum
func Md5(src string) string {
	md5ctx := md5.New()
	md5ctx.Write([]byte(src))
	cipher := md5ctx.Sum(nil)
	value := hex.EncodeToString(cipher)
	return value
}

// 字符串转换为小写，在转化前先判断是否包含大写字符，比strings.ToLower性能高
func ToLower(s string) string {
	if isLower(s) {
		return s
	}

	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return bytesconv.BytesToString(b)
}

func isLower(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}
