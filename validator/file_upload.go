package validator

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/fengzhongzhu1621/xgo/crypto/randutils"
	"github.com/h2non/filetype"
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

	// 2. 检查是否包含空字符 (\u0000)
	if strings.ContainsRune(path, '\x00') {
		return false
	}

	// 3. 检查路径是否包含非法 Unicode 字符（可选，如果仍然需要）
	if !utf8.ValidString(path) {
		return false
	}

	// 4. 解析相对路径（相对于 storagePath）
	relPath, err := filepath.Rel(storagePath, path)
	if err != nil {
		return false
	}

	// 5. 规范化路径（去除 ./ 和 ../）
	relPath = filepath.Clean(relPath)

	// 6. 检查路径穿越（更严格的方式）
	if strings.HasPrefix(relPath, "../") || relPath == ".." {
		return false
	}

	// 7. 检查非法字符（跨平台）
	illegalChars := `~$&*|<>?"'` // 基础非法字符
	if runtime.GOOS == "windows" {
		illegalChars += `:\/` // Windows 额外禁止的字符
	}
	if strings.ContainsAny(relPath, illegalChars) {
		return false
	}

	// 8. 检查保留文件名（跨平台）
	base := filepath.Base(relPath)
	switch strings.ToUpper(base) {
	case "CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return false // Windows 保留文件名
	case ".", "..":
		return false // Linux/Unix 保留文件名
	}

	// 9. 检查路径是否试图逃逸 storagePath（双重验证）
	absPath := filepath.Join(storagePath, relPath)
	if !strings.HasPrefix(absPath, filepath.Clean(storagePath)+string(filepath.Separator)) {
		return false
	}

	return true
}

// ValidateFileType 使用 filetype 库检测 MIME 类型
func ValidateFileType(file io.ReadSeeker, allowedMimeTypes map[string]bool) error {
	// 读取前 261 字节用于类型检测（filetype 库的推荐大小）
	buf := make([]byte, 261)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("读取文件头失败: %w", err)
	}

	if n == 0 {
		return fmt.Errorf("文件为空，无法检测 MIME 类型")
	}

	kind, err := filetype.Match(buf[:n])
	if err != nil {
		return fmt.Errorf("MIME 类型检测失败: %w", err)
	}

	if kind == filetype.Unknown {
		return fmt.Errorf("未知的文件类型")
	}

	// 检查 MIME 类型是否在允许的列表中，注意区分大小写
	mimeType := kind.MIME.Value
	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("不允许的文件类型: %s", mimeType)
	}

	// 重置文件指针到开头
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("重置文件指针失败: %w", err)
	}

	return nil
}
