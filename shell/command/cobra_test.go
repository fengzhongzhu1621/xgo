package command

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

// rootCmd是CLI应用程序的根命令。它是所有其他子命令的父命令，并且通常包含应用程序的主要逻辑或设置
// 定义了一个名为myapp的命令，它有一个简短的描述和一个更详细的描述。Run字段是一个函数，当执行myapp命令而没有指定任何子命令时，这个函数将被调用。
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
   examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 这里放置rootCmd的执行逻辑
		fmt.Println("Executing root command...")
	},
}

// 定义了一个名为serve的子命令，并将其添加到rootCmd中。现在，用户可以通过运行myapp serve来执行这个子命令。
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server...")
		// 服务器启动逻辑
	},
}

var verbose bool

func TestCobra(t *testing.T) {
	// 定义全局标志，这些标志对所有子命令都可用
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
