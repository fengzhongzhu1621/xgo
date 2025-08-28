package ssl

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

func TestNewECDSAKey(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "256",
			args: args{length: 256},
		},
		{
			name: "334",
			args: args{length: 384},
		},
		{
			name: "521",
			args: args{length: 521},
		},
		{
			name:    "other",
			args:    args{length: 11111},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建密钥生成器
			gen := NewSshKeyGenerator(zap.NewNop(), CA_CERT_TYPE_ECDSA, tt.args.length, 1024)
			defer gen.Close()

			var sshKey *SSHKey2
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				for {
					select {
					case <-ctx.Done():
						// 2 秒后结束循环
						return
					default:
						// 获得生成密钥
						sshKey = gen.Next()
						if sshKey == nil {
							// 队列为空
							time.Sleep(100 * time.Millisecond)
							continue
						}
						return
					}
				}
			}()

			if tt.wantErr {
				require.Nil(t, sshKey)
				_, _, err := NewECDSAKey(uint(tt.args.length))
				assert.Error(t, err)
				return
			}
			require.NotNil(t, sshKey)

			// 解码 PEM 格式的私钥
			// 将私钥字符串转换为 pem.Block 对象
			privPEM := sshKey.PrivateKey // 读取私钥文件
			privBlock, _ := pem.Decode([]byte(privPEM))
			// 解析 PKCS#8 私钥
			privKey, _ := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
			// 解析 SSH 授权密钥文件（通常是 ~/.ssh/authorized_keys 文件）中的公钥条目
			publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey.PubKey))
			assert.NoError(t, err)

			// 验证公钥是否与私钥匹配
			assert.EqualValues(t, &(privKey.(*ecdsa.PrivateKey).PublicKey), publicKey)
		})
	}
}
