package client

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	redisV9 "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
)

func TestSSHClient(t *testing.T) {
	// 定义了 SSH 客户端的配置
	sshConfig := &ssh.ClientConfig{
		User: "root",                                     // SSH 登录用户名
		Auth: []ssh.AuthMethod{ssh.Password("password")}, // 使用密码认证
		// 忽略主机密钥检查, 在测试环境中可以接受，但在生产环境中应使用 ssh.FixedHostKey 或类似机制验证主机密钥。
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second, // 连接超时时间
	}

	// 建立到远程服务器的 SSH 连接
	sshClient, err := ssh.Dial("tcp", "remoteIP:22", sshConfig)
	if err != nil {
		log.Fatalf("Failed to establish SSH connection: %v", err)
	}
	defer sshClient.Close()

	rdb := redisV9.NewClient(&redisV9.Options{
		// Redis 服务器的地址，这里使用 net.JoinHostPort 将 127.0.0.1 和 6379 组合成 "127.0.0.1:6379"
		Addr: net.JoinHostPort("127.0.0.1", "6379"),
		// 通过 SSH 隧道连接到 Redis 服务器
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial(network, addr)
		},
		// SSH不支持超时设置，在这里禁用
		// 虽然 SSH 隧道不支持超时设置，但可以在 Redis 客户端中设置其他超时
		// （如 DialTimeout、PoolTimeout 等）以避免长时间阻塞。
		ReadTimeout:  -1,
		WriteTimeout: -1,
	})

	// 测试 Redis 连接
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}
	log.Printf("Redis ping response: %s", pong)
}
