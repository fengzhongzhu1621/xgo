package debuglog

import "fmt"

// color is the type for log color.
type color uint8

// fontColor returns the font color of the log.
func (c color) fontColor() uint8 {
	return uint8(c)
}

const (
	colorRed     color = 31 // Red color.
	colorMagenta color = 35 // Magenta color.
)

// levelColorMap maps log levels to colors.
var levelColorMap = map[logLevel]color{
	debugLevel: colorMagenta,
	errorLevel: colorRed,
}

// getLogFormat formats the log output with different colors based on user settings.
// 根据日志级别和用户设置，为日志输出添加颜色控制字符，从而在支持颜色显示的终端中实现彩色日志输出。
func getLogFormat(level logLevel, enableColor bool, format string) string {
	if !enableColor {
		// Default scheme.
		return format
	}
	// 设置颜色前缀，%d 动态插入颜色代码（如 31 为红色）
	preColor := "\033[1;%dm" // Control character for color display.
	// 重置颜色后缀，避免后续文本继承颜色
	sufColor := "\033[0m" // Color display end character.

	if v, ok := levelColorMap[level]; ok {
		return fmt.Sprintf(preColor, v.fontColor()) + format + sufColor
	}
	return format
}
