package style

import "github.com/xuri/excelize/v2"

// Protection directly maps the protection settings of the cells.
type Protection struct {
	Hidden bool // 如果设置为 true，则在单元格被保护时，其内容在编辑栏中不可见。
	Locked bool // 如果设置为 true，则在单元格被保护时，用户无法更改单元格的内容。
}

func (p *Protection) Convert() (*excelize.Protection, error) {
	return &excelize.Protection{
		Hidden: p.Hidden,
		Locked: p.Locked,
	}, nil
}
