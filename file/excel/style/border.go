package style

import "github.com/xuri/excelize/v2"

type borderType string

// 边框的位置
const (
	Left   borderType = "left"
	Right  borderType = "right"
	Top    borderType = "top"
	Bottom borderType = "bottom"
)

// Border directly maps the border settings of the cells.
// 单元格的边框设置
type Border struct {
	Type  borderType // 边框的位置，使用 borderType 枚举（Left, Right, Top, Bottom）
	Color string     // 边框颜色的字符串表示（通常为颜色代码，如 "#000000"）
	Style int        // 边框样式的整数表示（具体样式依赖于 excelize 库的实现）
}

func (b *Border) Convert() (excelize.Border, error) {
	return excelize.Border{
		Type:  string(b.Type),
		Color: b.Color,
		Style: b.Style,
	}, nil
}
