package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
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

	// 返回证书池
	return certPool, nil
}

// LoadX509Certificates 从指定的客户端证书和私钥文件中加载 TLS 证书，并在需要时解密私钥，返回一个 *tls.Certificate 对象
func LoadX509Certificates(certFile, keyFile, passwd string) (*tls.Certificate, error) {
	// 读取私钥文件
	privateKey, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	// 如果提供了密码，尝试解密私钥
	if "" != passwd {
		// 首先使用 pem.Decode() 函数解码 PEM 格式的私钥
		privatePem, _ := pem.Decode(privateKey)
		if privatePem == nil {
			return nil, fmt.Errorf("decode private key failed")
		}
		// 如果解码失败，返回错误信息。然后使用 x509.DecryptPEMBlock() 函数解密私钥
		privateDecryptPem, err := x509.DecryptPEMBlock(privatePem, []byte(passwd))
		if err != nil {
			return nil, err
		}
		// 将解密后的私钥重新编码为 PEM 格式
		privateKey = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateDecryptPem,
		})
	}

	// 读取客户端证书
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	// 将证书数据和私钥数据组合成一个 tls.Certificate 对象
	tlsCert, err := tls.X509KeyPair(certData, privateKey)
	if err != nil {
		return nil, err
	}

	return &tlsCert, nil
}
