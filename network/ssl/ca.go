package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// GenerateCACertificatePEM 用于生成自签名CA（证书颁发机构）证书
func GenerateCACertificatePEM() (certPEM, keyPEM []byte, err error) {
	// 生成CA的密钥对
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("生成RSA密钥出错: %v", err)
	}
	pub := &priv.PublicKey

	// 创建CA证书模板：创建一个x509.Certificate的实例作为CA证书的模板
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization: []string{"Example CA"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),

		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 使用模板和私钥创建CA证书
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("创建CA证书失败: %v", err)
	}

	// 将证书转换为PEM格式
	certBuf := new(bytes.Buffer)
	pem.Encode(certBuf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	certPEM = certBuf.Bytes()

	// 将私钥转换为PEM格式
	keyBytes := x509.MarshalPKCS1PrivateKey(priv)
	keyBuf := new(bytes.Buffer)
	pem.Encode(keyBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})
	keyPEM = keyBuf.Bytes()

	return certPEM, keyPEM, nil
}
