package style

import "github.com/xuri/excelize/v2"

// Font directly maps the font settings of the fonts.
type Font struct {
	Bold         bool    // 是否加粗
	Italic       bool    // 是否斜体
	Underline    string  // 下划线样式（如 "single", "double"）
	Family       string  // 字体家族（如 "Arial", "Times New Roman"）
	Size         float64 // 字体大小（以磅为单位）
	Strike       bool    // 是否删除线
	Color        string  // 字体颜色（通常为十六进制颜色代码，如 "#000000"）
	ColorIndexed int     // 索引颜色（基于颜色索引表）
	ColorTheme   *int    // 主题颜色（指向主题颜色索引的指针）
	ColorTint    float64 // 颜色色调（用于调整主题颜色的色调）
	VertAlign    string  // 垂直对齐方式（如 "baseline", "subscript", "superscript"）
}

func (f *Font) Convert() (*excelize.Font, error) {
	return &excelize.Font{
		Bold:         f.Bold,
		Italic:       f.Italic,
		Underline:    f.Underline,
		Family:       f.Family,
		Size:         f.Size,
		Strike:       f.Strike,
		Color:        f.Color,
		ColorIndexed: f.ColorIndexed,
		ColorTheme:   f.ColorTheme,
		ColorTint:    f.ColorTint,
		VertAlign:    f.VertAlign,
	}, nil
}
