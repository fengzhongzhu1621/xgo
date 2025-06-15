package validator

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/fengzhongzhu1621/xgo/crypto/randutils"

	"github.com/samber/lo"
)

type FileValidator struct {
	MaxSize    int64    // 文件的最大值
	AllowTypes []string // 允许的文件类型，注意都是大写
}

func NewFileValidator(maxSize int64, allowTypes []string) *FileValidator {
	allowTypesList := lo.Map(allowTypes, func(x string, _ int) string {
		return strings.ToUpper(x)
	})
	obj := &FileValidator{
		MaxSize:    maxSize,
		AllowTypes: allowTypesList,
	}

	return obj
}

func (v *FileValidator) ValidateFile(file *multipart.FileHeader) (string, error) {
	// 检查文件大小
	if file.Size > v.MaxSize {
		return "", errors.New("file size exceeds limit")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !lo.Contains(v.AllowTypes, strings.ToUpper(ext)) {
		return "", errors.New("file type not allowed")
	}

	// 生成随机文件名
	randomName, err := randutils.GenerateRandomFileName(ext)
	if err != nil {
		return "", err
	}

	return randomName, nil
}

// IsPathValid 验证路径安全性
func IsPathValid(storagePath, path string) bool {
	// 1. 检查路径是否为空
	if path == "" {
		return false
	}

	// 2. 检查路径是否包含非法 Unicode 字符
	// 检查字符串是否是合法的 UTF-8 编码，但它不会拒绝包含 \u0000（空字符）的字符串，
	// 因为 \u0000 本身是合法的 Unicode 码点（尽管在文件路径中通常不允许）。
	if !utf8.ValidString(path) {
		return false
	}
	// 检查是否包含空字符 (\u0000)
	if strings.ContainsRune(path, '\x00') {
		return false
	}

	// 3. 解析相对路径（相对于 storagePath）
	relPath, err := filepath.Rel(storagePath, path)
	if err != nil {
		return false
	}

	// 4. 规范化路径（去除 ./ 和 ../）
	relPath = filepath.Clean(relPath)

	// 5. 检查路径穿越（更严格的方式）
	if strings.HasPrefix(relPath, "../") || relPath == ".." {
		return false
	}

	// 6. 检查非法字符（跨平台）
	illegalChars := `~$&*|<>?"'` // 基础非法字符
	if runtime.GOOS == "windows" {
		illegalChars += `:\/` // Windows 额外禁止的字符
	}
	if strings.ContainsAny(relPath, illegalChars) {
		return false
	}

	// 7. 检查保留文件名（跨平台）
	base := filepath.Base(relPath)
	switch strings.ToUpper(base) {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return false // Windows 保留文件名
	case ".", "..":
		return false // Linux/Unix 保留文件名
	}

	// 8. 检查路径是否试图逃逸 storagePath（双重验证）
	absPath := filepath.Join(storagePath, relPath)
	if !strings.HasPrefix(absPath, filepath.Clean(storagePath)+string(filepath.Separator)) {
		return false
	}

	return true
}
