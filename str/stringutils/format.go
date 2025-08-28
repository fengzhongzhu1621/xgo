package stringutils

import (
	"sort"
	"strings"
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

func MergeGetAndPostParamWithKey(queryParam map[string]string,
	postParam map[string]string, key string, keyName string,
) string {
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
