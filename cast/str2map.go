package cast

import (
	"strings"
)

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
