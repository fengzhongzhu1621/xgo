package style

import (
	"github.com/xuri/excelize/v2"
)

// Style excel style
type Style struct {
	Fill   *Fill    // 单元格的填充样式
	Border []Border // 单元格的多个边框样式
	Font   *Font    // 单元格的字体样式
}

// Convert 将 Style 结构体转换为 excelize.Style
func (s *Style) Convert() (*excelize.Style, error) {
	// 创建一个新的 excelize.Style 实例。
	style := new(excelize.Style)

	if s.Fill != nil {
		fill, err := s.Fill.Convert()
		if err != nil {
			return nil, err
		}
		style.Fill = fill
	}

	if s.Font != nil {
		font, err := s.Font.Convert()
		if err != nil {
			return nil, err
		}
		style.Font = font
	}

	if s.Border != nil {
		for _, border := range s.Border {
			excelBorder, err := border.Convert()
			if err != nil {
				return nil, err
			}
			style.Border = append(style.Border, excelBorder)
		}
	}

	return style, nil
}
