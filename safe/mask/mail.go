package mask

import (
	"strings"

	"github.com/duke-git/lancet/v2/strutil"
)

// MaskEmail 隐藏邮箱ID的中间部分
func MaskEmail(address string) string {
	// 获得邮箱ID
	id := address[0:strings.LastIndex(address, "@")]
	// 获得邮箱域名
	domain := address[strings.LastIndex(address, "@"):]

	// 邮箱ID小于等于1个字符，不处理
	if len(id) <= 1 {
		return address
	}

	// 根据邮箱ID的长度，决定星号的位置
	switch len(id) {
	case 2:
		id = strutil.Concat(2, id[0:1], "*")
	case 3:
		id = strutil.Concat(3, id[0:1], "*", id[2:])
	case 4:
		id = strutil.Concat(4, id[0:1], "**", id[3:])
	default:
		// 4个字符以上的邮箱ID，中间使用星号填充
		masks := strings.Repeat("*", len(id)-4)
		id = strutil.Concat(0, id[0:2], masks, id[len(id)-2:])
	}

	// 拼接邮箱
	return strutil.Concat(0, id, domain)
}
