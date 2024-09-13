package ssl

import (
	"go.uber.org/zap"
)

type SSHKey struct {
	Name        string
	Type        int  // 非对称加密类型
	Length      uint // 密钥长度
	Fingerprint string
	PrivKey     string // 私钥
	PubKey      string // 公钥
	Comment     string
}

type SSHKey2 struct {
	PubKey  string
	PrivKey string
	Type    int
}

type SshKeyGenerator struct {
	Size    int // 队列长度
	Length  int // 密钥长度
	Queue   chan SSHKey2
	Type    int
	Running bool // 是否运行生成器
	loger   *zap.Logger
}
