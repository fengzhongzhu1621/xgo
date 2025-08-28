package file

import (
	"runtime"
	"strings"

	"github.com/fengzhongzhu1621/xgo/file/pathutils"
)

// SanitizedName 格式化文件名
func SanitizedName(filename string) string {
	// 去掉 C:
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	// 将 Windows 风格的路径中的反斜杠替换为正斜杠，而对于 Unix 风格的路径，它不执行任何操作
	filename = pathutils.SlashAndCleanPath(filename)

	// 去掉开头的 /
	filename = strings.TrimLeft(strings.ReplaceAll(filename, `\`, "/"), `/`)

	return filename
}
