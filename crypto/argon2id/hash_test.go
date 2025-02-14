package argon2id

import (
	"fmt"
	"log"
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestArgon2idHash(t *testing.T) {
	// 创建哈希配置
	params := &argon2id.Params{
		Memory:      64 * 1024, // 内存消耗（字节）
		Iterations:  1,         // 迭代次数
		Parallelism: 2,         // 并行度
		SaltLength:  16,        // 盐值长度（字节）
		KeyLength:   32,        // 哈希输出长度（字节）
	}

	// 生成密码哈希
	hash, err := argon2id.CreateHash("xx123456xx", params)
	if err != nil {
		log.Fatalf("创建哈希失败: %v", err)
	}

	fmt.Println("生成的哈希:", hash)

	// 验证密码
	testPassword := "xx123456xx"
	match, err := argon2id.ComparePasswordAndHash(testPassword, hash)
	if err != nil {
		log.Fatalf("验证密码时出错: %v", err)
	}

	if match {
		fmt.Println("密码验证成功")
	} else {
		fmt.Println("密码验证失败")
	}

	// 尝试使用错误的密码
	wrongPassword := "wrongPassword"
	match, err = argon2id.ComparePasswordAndHash(hash, wrongPassword)
	if err != nil {
		log.Fatalf("验证密码时出错: %v", err)
	}

	if match {
		fmt.Println("密码验证成功")
	} else {
		fmt.Println("密码验证失败")
	}
}
