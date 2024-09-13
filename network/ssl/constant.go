package ssl

import "regexp"

var OpensslVersionRegex = regexp.MustCompile(`(?i)OpenSSH_([0-9]+\.[0-9]+)`)
var OpensslVersionRegexV2 = regexp.MustCompile(`(?i)OpenSSH_([0-9]+)\.([0-9]+)`)

// 非对称加密算法类型
const (
	CA_CERT_TYPE_RSA     = iota + 1
	CA_CERT_TYPE_DSA     // 由于安全性问题，已被逐渐淘汰，兼容性较差。
	CA_CERT_TYPE_ECDSA   // 基于椭圆曲线数字签名算法
	CA_CERT_TYPE_ED25519 // 基于爱德华曲线的数字签名算法
)

const (
	// 非对称密钥长度
	DEFAULT_RSA_KEY_LENGTH = 2048
)
