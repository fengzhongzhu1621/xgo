package stringutils

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fengzhongzhu1621/xgo/str/bytesconv"
)

// Last 返回数组最后一个元素.
func Last(list []string) string {
	return list[len(list)-1]
}

// LastChar 获得最后一个字符.
func LastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// GetValueInBraces 获得大括号中间的值.
func GetValueInBraces(key string) string {
	if s := strings.IndexByte(key, '{'); s > -1 {
		if e := strings.IndexByte(key[s+1:], '}'); e > 0 {
			return key[s+1 : s+e+1]
		}
	}
	return key
}

// RemoveDuplicateElement 切片去重
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

// TrimRight 去掉字符串的后缀.
func TrimRight(str string, substring string) string {
	idx := strings.LastIndex(str, substring)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

// TrimLeft 去掉字符串的前缀.
func TrimLeft(str string, substring string) string {
	return strings.TrimPrefix(str, substring)
}

// ToLower 字符串转换为小写，在转化前先判断是否包含大写字符，比strings.ToLower性能高.
func ToLower(s string) string {
	// 判断字符串是否包含小写字母
	if IsLower(s) {
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
	// []bytes转换为字符串
	return bytesconv.BytesToString(b)
}

// UnicodeTitle 首字母大写.
func UnicodeTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

// UnicodeUnTitle 首字母小写.
func UnicodeUnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

// ChangeInitialCase 将字符串的首字符转换为指定格式
func ChangeInitialCase(s string, mapper func(rune) rune) string {
	if s == "" {
		return s
	}
	// 返回第一个utf8字符，n是字符的长度，即返回首字符
	r, n := utf8.DecodeRuneInString(s)
	// 根据mapper方法转换首字符
	return string(mapper(r)) + s[n:]
}

func Str2map(s string, sep1 string, sep2 string) map[string]string {
	if s == "" {
		return nil
	}
	spe1List := strings.Split(s, sep1)
	if len(spe1List) == 0 {
		return nil
	}
	m := make(map[string]string)
	for _, sub := range spe1List {
		splitNum := 2
		spe2List := strings.SplitN(sub, sep2, splitNum)
		num := len(spe2List)
		if num == 1 {
			m[spe2List[0]] = ""
		} else if num > 1 {
			m[spe2List[0]] = spe2List[1]
		}
	}
	return m
}

func MergeGetAndPostParamWithKey(queryParam map[string]string,
	postParam map[string]string, key string, keyName string) string {
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
	for k := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)

	// 排序后的数组
	params := ""
	for _, key := range keyList {
		if value := m[key]; value != "" {
			params += key + "=" + value + "&"
		}
	}
	// 添加key参数
	params += keyName + "=" + key
	return params
}