package flagutils

import (
	"flag"
	"fmt"
	"testing"

	"github.com/spf13/pflag"
)

func TestAddGoFlagSet(t *testing.T) {
	// 创建一个新的 pflag.FlagSet
	fs1 := pflag.NewFlagSet("fs1", pflag.ExitOnError)
	fs1.String("flag1", "default1", "This is flag1 from fs1")

	// 创建一个 Go 标准库 flag.FlagSet
	fs2 := flag.NewFlagSet("fs2", flag.ExitOnError)
	fs2.String("flag2", "default2", "This is flag2 from fs2")

	// 将 fs2 添加到 fs1 中
	fs1.AddGoFlagSet(fs2)

	// 解析命令行参数
	fs1.Parse([]string{"--flag1=value1", "--flag2=value2"})

	// 输出解析后的参数值
	fmt.Println("flag1:", fs1.Lookup("flag1").Value.String())
	fmt.Println("flag2:", fs2.Lookup("flag2").Value.String())
}
