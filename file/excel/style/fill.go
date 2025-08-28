package style

import "github.com/xuri/excelize/v2"

type FillType string

const (
	// Pattern pattern fill type 图案填充类型
	Pattern FillType = "pattern"
	// Gradient gradient fill type 渐变填充类型
	Gradient FillType = "gradient"
)

// Fill directly maps the fill settings of the cells.
type Fill struct {
	Type    FillType // 填充类型，使用 FillType 枚举（Pattern 或 Gradient）
	Pattern int      // 图案填充的样式编号（具体含义依赖于 excelize 库的实现）
	Color   []string // 填充颜色的字符串切片（通常为颜色代码，如 ["#FFFFFF", "#000000"]）
	Shading int      // 渐变填充的阴影或角度设置（具体含义依赖于 excelize 库的实现）
}

func (f *Fill) Convert() (excelize.Fill, error) {
	return excelize.Fill{
		Type:    string(f.Type),
		Pattern: f.Pattern,
		Color:   f.Color,
		Shading: f.Shading,
	}, nil
}
