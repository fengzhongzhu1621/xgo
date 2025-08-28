package lock_dir

import (
	"io"
	"os"
	"testing"
)

func TestOpenRoot(t *testing.T) {
	configdir := "../../config"

	_ = os.MkdirAll(configdir, 0o755)

	// os.OpenRoot 打开并锁定目录
	root, err := os.OpenRoot(configdir)
	if err != nil {
		panic(err)
	}

	// root.OpenFile 打开文件， 可以是基于 root 的相对路径
	f, err := root.OpenFile("config.yaml", os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 不支持 os.ReadFile 这样的快捷方式。
	// 使用 io 读取文件
	_, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}
}
