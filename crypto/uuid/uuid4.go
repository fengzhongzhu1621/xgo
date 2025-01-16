package uuid

import (
	"github.com/fengzhongzhu1621/xgo/crypto/randutils"
	"github.com/google/uuid"
)

// GetUUIDv4 获取v4版本UUID
func GetUUIDv4() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// GetLocalUUIDv4 获取UUID，无法生成则通过本地随机一个字符串
func GetLocalUUIDv4() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return randutils.RandomString(24)
	}
	return id.String()
}

// IsValidUUID 判断UUID是否有效
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
