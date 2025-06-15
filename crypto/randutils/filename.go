package randutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateRandomFileName 生成随机文件名
func GenerateRandomFileName(ext string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes) + ext, nil
}

func GenerateSafeFileNameUseUUID(original string) string {
	ext := filepath.Ext(original)
	uuidStr := uuid.New().String()
	return uuidStr + ext
}

// generateSafeFileName 生成一个安全的文件名
// 1. 去除路径信息（防止路径遍历攻击）
// 2. 添加哈希前缀避免冲突
// 3. 保留原始扩展名
func GenerateSafeFileNameUseHash(original string) string {
	if original == "" {
		return fmt.Sprintf("file_%x", md5.Sum([]byte(time.Now().String())))
	}

	// 1. 去除路径信息，只保留文件名部分
	name := filepath.Base(original)

	// 2. 提取扩展名
	ext := filepath.Ext(name)
	baseName := strings.TrimSuffix(name, ext)

	// 3. 生成哈希前缀（基于原始文件名 + 当前时间）
	hash := md5.Sum([]byte(baseName + time.Now().Format("20060102150405.000")))

	// 4. 组合成最终文件名: 哈希前缀 + 扩展名
	safeName := fmt.Sprintf("%x%s", hash, ext)

	return safeName
}
