package ssl

import "crypto/tls"

// ClientTslConfNoVerity InsecureSkipVerify是一个布尔类型的字段，用于控制是否跳过TLS服务器证书验证。
// 当设置为true时，客户端将不会验证服务器的TLS证书是否有效，这可能会导致安全风险，因为中间人攻击者可能会伪造证书。
func ClientTslConfNoVerity() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}
