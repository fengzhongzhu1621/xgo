package upload

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

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
