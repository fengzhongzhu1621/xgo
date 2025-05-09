package color

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestColor(t *testing.T) {
	// 支持打印有颜色的字符串
	color.Blue("hello %s", "world")

	colorPrint := color.New()
	colorPrint.Add(color.FgRed)   // 添加前景色
	colorPrint.Add(color.Italic)  // 设置为斜体
	colorPrint.Add(color.BgGreen) // 设置背景色

	colorPrint.Println("hello world")

	beeStrUp := "beeStrUp"
	beeStrDown := "beeStrDown"
	// 设置前景色为紫红色
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset()
	fmt.Println(beeStrUp)
	// 设置前景色为绿色
	color.Set(color.FgGreen, color.Bold)
	fmt.Println(beeStrDown)
}
