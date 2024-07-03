package file

import (
	"errors"
	"runtime"
	"strings"

	"github.com/fengzhongzhu1621/xgo/file/pathutils"
)

// 格式化文件名
func SanitizedName(filename string) string {
	// 去掉 C:
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows" {
		filename = filename[2:]
	}
	// 将 Windows 风格的路径中的反斜杠替换为正斜杠，而对于 Unix 风格的路径，它不执行任何操作
	filename = pathutils.SlashAndCleanPath(filename)

	// 去掉开头的 /
	filename = strings.TrimLeft(strings.Replace(filename, `\`, "/", -1), `/`)

	return filename
}


// 判断文件名是否包含特殊字符
func CheckFilename(name string) error {
	// 用于检查一个字符串中是否包含指定的任何字符集。
	if strings.ContainsAny(name, "\\/:*<>|") {
		// 自定义错误
		return errors.New("Name should not contains \\/:*<>|")
	}
	return nil
}
