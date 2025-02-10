package flagutils

import (
	"flag"
	"os"

	"github.com/fengzhongzhu1621/xgo/version"
	"github.com/spf13/pflag"
)

// InitFlags normalizes, parses, then logs the command line flags
func InitFlags() {
	// 设置命令行参数名称规范化函数的函数，可以使用规范化函数将参数名称统一为一种格式。
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
	// 用于将一个 *flag.FlagSet 添加到全局的 flag.CommandLine 中。这样，你可以在一个程序中使用多个 FlagSet，并将它们组合在一起。
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// 添加帮助
	var help bool
	pflag.CommandLine.BoolVarP(&help, "help", "h", false, "show help info")

	// 添加版本
	ver := pflag.CommandLine.Bool("version", false, "print version information")

	// 解析命令行参数
	pflag.Parse()

	if help {
		// 打印所有已注册的命令行参数及其默认值
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *ver {
		version.ShowVersion()
		os.Exit(0)
	}
}
