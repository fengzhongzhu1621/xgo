package user

import (
	"fmt"
	"os/user"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	username := "root"
	u, err := user.Lookup(username)
	if err != nil {
		fmt.Printf("查找用户失败: %v\n", err)
		return
	}

	fmt.Printf("用户名: %s\n", u.Username) // root
	fmt.Printf("UID: %s\n", u.Uid)      // 0
	fmt.Printf("家目录: %s\n", u.HomeDir)  // var/root
}
