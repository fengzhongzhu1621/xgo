package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// GenerateHmacSignedString 生成签名串.
func GenerateHmacSignedString(algorithm string, key []byte, p []byte) (string, error) {
	var h hash.Hash
	if algorithm == "sha256" {
		h = hmac.New(sha256.New, key)
	} else {
		// 默认算法
		h = hmac.New(sha256.New, key)
	}
	_, err := h.Write(p)
	if err != nil {
		return "", err
	}
	// base64编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
