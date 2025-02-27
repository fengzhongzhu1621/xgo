package xid

import (
	"fmt"
	"testing"

	"github.com/rs/xid"
)

func TestNew(t *testing.T) {
	guid := xid.New()
	// cuvauq164h7mlc53qgjg
	fmt.Println(guid.String())
}

func TestMeta(t *testing.T) {
	guid := xid.New()
	// [38 36 79]
	fmt.Println(guid.Machine()) // 获取机器 ID
	// 27393
	fmt.Println(guid.Pid()) // 获取进程 ID
	// 2025-02-26 14:06:38 +0800 CST
	fmt.Println(guid.Time()) // 获取时间戳
	// 4240343
	fmt.Println(guid.Counter()) // 获取计数器
}
