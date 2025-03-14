package color

import "fmt"

type Color int

const (
	// 定义前景色
	Cyan    Color = 36
	Green   Color = 32
	Magenta Color = 35
	Red     Color = 31
	Yellow  Color = 33
	// 定义背景色
	BackgroundRed Color = 41
)

// ColorResetEscape 颜色重置配置
var ColorResetEscape = "\033[0m"

// ColorEscape 生成前景色颜色转义
func ColorEscape(color Color) string {
	return fmt.Sprintf("\033[0;%dm", color)
}

// Colorize 使用指定颜色显示 msg，然后重置
func Colorize(msg string, color Color) string {
	return ColorEscape(color) + msg + ColorResetEscape
}
