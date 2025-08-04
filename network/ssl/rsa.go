package ssl

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// NewRSAKeyPair 生成RSA密钥对
func NewRSAKeyPair(keyLen int) (string, string, error) {
	// RSA密钥长度
	if keyLen == 0 {
		keyLen = DEFAULT_RSA_KEY_LENGTH
	}
	// 生成一个RSA私钥
	private, err := rsa.GenerateKey(rand.Reader, keyLen)
	if err != nil {
		return "", "", err
	}
	// 从私钥中提取公钥
	publicKey, err := ssh.NewPublicKey(&private.PublicKey)
	if err != nil {
		return "", "", err
	}
	// 将私钥编码为PEM格式的字节切片
	privateKey := pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(private),
		Type:  "RSA PRIVATE KEY",
	})

	// 返回编码后的私钥和公钥，以及一个nil错误
	// MarshalAuthorizedKey 将一个 SSH 公钥对象（通常是 *ssh.PublicKey 类型）转换为可以在 SSH 授权文件（如 ~/.ssh/authorized_keys）中使用的格式
	return string(privateKey), string(ssh.MarshalAuthorizedKey(publicKey)), nil
}

// NewRSAKey 生成RSA密钥
// param length: >=1024 && <= 16384
func NewRSAKey(length uint) (*pem.Block, ssh.PublicKey, error) {
	if length < 1024 || length > 16384 {
		return nil, nil, fmt.Errorf(
			"key length not supported: %d, supported values are between 1024 and 16384",
			length,
		)
	}
	// 生成一个RSA私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, int(length))
	if err != nil {
		return nil, nil, err
	}
	// 从私钥中提取公钥
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	// 将生成的私钥转换为 PEM 格式
	pemKey := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// 返回生成的 PEM 格式的私钥、SSH 公钥和错误信息
	return pemKey, publicKey, err
}

// NewECDSAKey 生成ECDSA密钥
// param length: 长度限制为256,384,521
func NewECDSAKey(length uint) (*pem.Block, ssh.PublicKey, error) {
	// 根据传入的 length 参数选择椭圆曲线
	var curve elliptic.Curve
	switch length {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, nil, fmt.Errorf("ECDSA key length not supported: %d[256/384/521]", length)
	}

	// 生成 ECDSA 私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// 从私钥中提取公钥
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	// 将生成的私钥转换为 PKCS8 格式，并封装为 PEM 格式
	marshaledKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	pemKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: marshaledKey,
	}
	if err != nil {
		return nil, nil, err
	}
	// 返回生成的 PEM 格式的私钥、SSH 公钥和错误信息
	return pemKey, publicKey, err
}

// NewEd25519Key 生成Ed25519密钥
func NewEd25519Key() (*pem.Block, ssh.PublicKey, error) {
	publicKeyEd25519, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	// convert priv key to x509 format
	marshaledKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	pemKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: marshaledKey,
	}
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := ssh.NewPublicKey(publicKeyEd25519)
	if err != nil {
		return nil, nil, err
	}
	return pemKey, publicKey, err
}
