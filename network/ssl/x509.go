package ssl

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// These are utilities for working with files containing ECDSA public and
// private keys. See this helpful doc for how to generate them:
// https://wiki.openssl.org/index.php/Command_Line_Elliptic_Curve_Operations
//
// The quick cheat sheet below.
//  1. Generate an ECDSA-P256 private key
//     openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
//  2. Generate public key from private key
//     openssl ec -in ecprivatekey.pem -pubout -out ecpubkey.pem

// GetCertPool gets CertPool information.
func GetCertPool(caCertFile string) (*x509.CertPool, error) {
	// root means to use the root ca certificate installed on the machine to
	// verify the peer, if not root, use the input ca file to verify peer.
	if caCertFile == "root" {
		return nil, nil
	}

	// 读取自定义的 CA 证书文件
	ca, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("read ca file error: %w", err)
	}

	// 创建一个新的 CertPool
	certPool := x509.NewCertPool()
	// 将 CA 证书添加到 CertPool
	if !certPool.AppendCertsFromPEM(ca) {
		return nil, errors.New("AppendCertsFromPEM fail")
	}
	return certPool, nil
}

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
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
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
	if passwd != "" {
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
	// 	客户端证书通常包含以下核心元素：
	//
	// 版本号 - 证书格式版本
	// 序列号 - 证书的唯一标识符
	// 签名算法 - 用于签发证书的算法
	// 颁发者 - 签发证书的CA信息
	// 有效期 - 证书的起止时间
	// 主体信息 - 证书持有者的标识信息
	// 公钥信息 - 包含公钥和算法参数
	// 扩展字段 - 如密钥用法、扩展密钥用法等
	// 签名值 - CA对证书内容的数字签名
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

// LoadEcdsaPublicKey reads an ECDSA public key from an X509 encoding stored in a PEM encoding.
func LoadEcdsaPublicKey(buf []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(buf)

	if block == nil {
		return nil, errors.New("no PEM data block found")
	}
	// The public key is loaded via a generic loader. We use X509 key format,
	// which supports multiple types of keys.
	keyIface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error loading public key: %w", err)
	}

	// Now, we're assuming the key content is ECDSA, and converting.
	publicKey, ok := keyIface.(*ecdsa.PublicKey)
	if !ok {
		// The cast failed, we might have loaded an RSA file or something.
		return nil, errors.New("file contents were not an ECDSA public key")
	}
	return publicKey, nil
}

// LoadEcdsaPrivateKey reads an ECDSA private key from an X509 encoding stored in a PEM encoding
func LoadEcdsaPrivateKey(buf []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(buf)

	if block == nil {
		return nil, errors.New("no PEM data block found")
	}

	// At this point, we've got a valid PEM data block. PEM is just an encoding,
	// and we're assuming this encoding contains X509 key material.
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error loading private ECDSA key: %w", err)
	}
	return privateKey, nil
}

// StoreEcdsaPublicKey writes an ECDSA public key to a PEM encoding
func StoreEcdsaPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {
	encodedKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error x509 encoding public key: %w", err)
	}
	pemEncodedKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: encodedKey,
	})
	return pemEncodedKey, nil
}

// StoreEcdsaPrivateKey writes an ECDSA private key to a PEM encoding
func StoreEcdsaPrivateKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	encodedKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error x509 encoding private key: %w", err)
	}
	pemEncodedKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: encodedKey,
	})
	return pemEncodedKey, nil
}
