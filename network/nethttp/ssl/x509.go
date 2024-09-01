package ssl

import (
	"crypto/x509"
	"fmt"
	"os"
)

// LoadX509Cert 读取 CA 证书文件
func LoadX509Cert(caFile string) (*x509.CertPool, error) {
	// 读取自定义的 CA 证书文件
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("读取 CA 证书失败: %s", err)
	}
	// 创建一个新的 CertPool
	certPool := x509.NewCertPool()
	// 将 CA 证书添加到 CertPool
	if ok := certPool.AppendCertsFromPEM(caCert); ok != true {
		return nil, fmt.Errorf("append ca cert failed")
	}

	return certPool, nil
}
