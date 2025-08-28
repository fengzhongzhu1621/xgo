package numpy

import (
	"strings"
	"unicode"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/numpy/math"
)

// ParseSizeInBytes converts strings like 1GB or 12 mb into an unsigned integer number of bytes.
// 将 KB, MB, GB 转换为字节数
func ParseSizeInBytes(sizeStr string) uint {
	// 去掉空白字符
	sizeStr = strings.TrimSpace(sizeStr)
	// 获得单位的索引
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		// 如果单位是字节
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			// 如果单位是KB, MB, GB
			if lastChar > 1 {
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size := cast.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

	return math.SafeMul(uint(size), multiplier)
}
