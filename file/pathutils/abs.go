package pathutils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fengzhongzhu1621/xgo/file/homedir"
)

// AbsPathify 将路径转换为绝对路径.
func AbsPathify(inPath string) string {
	// 如果是home路径，则转换为绝对路径
	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		inPath = homedir.UserHomeDir() + inPath[5:]
	}

	// 路径模板字符串替换
	inPath = os.ExpandEnv(inPath)

	// 判断路径是否是决定路径
	if filepath.IsAbs(inPath) {
		// 清理路径中的多余字符，比如 /// 或 ../ 或 ./
		return filepath.Clean(inPath)
	}

	// 转换为绝对路径
	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}
