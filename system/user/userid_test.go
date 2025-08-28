package user

import (
	"fmt"
	"os/user"
	"testing"
)

func TestGetUserId(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("无法获取当前用户信息:", err)
		return
	}

	// 可以通过操作系统命令：id 查询
	// 当前用户的 UID 是: 501
	fmt.Println("当前用户的 UID 是:", currentUser.Uid)
}
