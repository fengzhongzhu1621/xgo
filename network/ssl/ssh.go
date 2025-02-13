package ssl

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// NewSSHKey 生成SSH密钥
func NewSSHKey(keyType int, length uint) (*SSHKey, error) {
	key := SSHKey{
		Type:   keyType,
		Length: length,
	}

	var err error
	var pemKey *pem.Block
	var publicKey ssh.PublicKey

	// 生成非对称密钥，包含公钥和私钥
	switch keyType {
	case CA_CERT_TYPE_RSA:
		pemKey, publicKey, err = NewRSAKey(length)
	case CA_CERT_TYPE_ECDSA:
		pemKey, publicKey, err = NewECDSAKey(length)
	case CA_CERT_TYPE_ED25519:
		pemKey, publicKey, err = NewEd25519Key()
	default:
		return nil, fmt.Errorf("key type not supported: %q[rsa/ecdsa/ed25519]", key.Type)
	}
	if err != nil {
		return nil, err
	}

	// 将生成的私钥编码为PEM格式，存储在buf中
	buf := bytes.NewBufferString("")
	if err = pem.Encode(buf, pemKey); err != nil {
		return nil, err
	}

	// 私钥
	key.PrivateKey = buf.String()
	// 公钥 MarshalAuthorizedKey 将一个 SSH 公钥对象（通常是 *ssh.PublicKey 类型）转换为可以在 SSH 授权文件（如 ~/.ssh/authorized_keys）中使用的格式
	key.PubKey = strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publicKey)))

	return &key, nil
}

// NewSshKeyGenerator 创建 ssh 密钥生成器
// - keyType 非对称算法类型
// - keyLength 密钥队列长度
func NewSshKeyGenerator(logger *zap.Logger, keyType, keyLength, size int) *SshKeyGenerator {
	switch keyType {
	case CA_CERT_TYPE_RSA:
	case CA_CERT_TYPE_ECDSA:
	case CA_CERT_TYPE_ED25519:
		break
	default:
		panic("keytype error")
	}

	// 创建密钥生成器
	generator := &SshKeyGenerator{
		Size:    size, // 密钥队列的长度
		Running: true,
		Length:  keyLength,
		Type:    keyType,
		logger:  logger,
		Queue:   make(chan SSHKey2, size),
	}

	// 启动密钥生成器
	go generator.run()

	return generator
}

// run 循环生成密钥
func (g *SshKeyGenerator) run() {
	defer func() { recover() }()

	for {
		var (
			err error
		)
		// 生成密钥
		sshKey, err := NewSSHKey(g.Type, uint(g.Length))
		if err != nil {
			g.logger.Error("[SshKeyGenerator]Generator RSA key pairs error", zap.String("errmsg", err.Error()))
			// 失败重试 in case of cpu 100%
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// 将生成的密钥放入队列
		queueValue := SSHKey2{
			PubKey:     sshKey.PubKey,
			PrivateKey: sshKey.PrivateKey,
			Type:       g.Type,
		}
		g.Queue <- queueValue
	}
}

// Next 从队列获取一个生成的密钥
func (g *SshKeyGenerator) Next() *SSHKey2 {
	select {
	case it := <-g.Queue:
		return &it
	default:
		return nil
	}
}

// Close 关闭生成器
func (g *SshKeyGenerator) Close() {
	g.Running = false
	close(g.Queue)
}
