package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

const (
	rowStartIdx = 1
	colStartIdx = 1
)

// Cell can be used directly in StreamWriter.SetRow to specify a style and a value.
type Cell struct {
	StyleID int         // 单元格的样式 ID
	Value   interface{} // 单元格的值
}

// transfer 将自定义的 Cell 结构体转换为 excelize.Cell 结构体，方便与 excelize 库交互。
func (c *Cell) transfer() *excelize.Cell {
	return &excelize.Cell{
		StyleID: c.StyleID,
		Value:   c.Value,
	}
}

// GetCellIdx get cell index
// 将从 0 开始计数的列和行索引转换为 Excel 的单元格名称（如 "A1"）。
// 允许调用者使用更直观的从 0 开始的索引，而无需关心 excelize 库从 1 开始的索引。
func GetCellIdx(col int, row int) (string, error) {
	// 由于第三方库的行和列不是从0开始，所以这里加上开始数，使调用者可以按照从0开始进行计数
	return excelize.CoordinatesToCellName(col+colStartIdx, row+rowStartIdx)
}

// GetSingleColSqref get single column sqref
// Example: GetSingleColSqref(0, 1, 2) // return A1:A2
// 获取单列的引用范围，如 "A1:A2"。方便生成对单一列多个行的引用，适用于批量操作或样式应用。
func GetSingleColSqref(col, startRow, endRow int) (string, error) {
	colNum, err := ColumnNumberToName(col)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%d:%s%d", colNum, startRow, colNum, endRow), nil
}

// ColumnNumberToName provides a function to convert the integer to Excel
// sheet column title.
// Example:	ColumnNumberToName(0) // returns "A", nil
// 将从 0 开始的列索引转换为 Excel 列名称（如 0 -> "A"）。
func ColumnNumberToName(col int) (string, error) {
	// 由于第三方库的列是从1开始，所以这里进行了+1操作，使调用者可以按照从0开始进行计数
	return excelize.ColumnNumberToName(col + 1)
}

// GetTotalRows get total rows 获取 Excel 表格的总行数。
func GetTotalRows() int {
	return excelize.TotalRows
}

// CellNameToCoordinates converts alphanumeric cell name to [X, Y] coordinates
// or returns an error.
// 将 Excel 的单元格名称（如 "A1"）转换为列和行的索引。
//
// Example:
//
//	CellNameToCoordinates("A1") // returns 1, 1, nil
//	CellNameToCoordinates("Z3") // returns 26, 3, nil
func CellNameToCoordinates(cell string) (int, int, error) {
	return excelize.CellNameToCoordinates(cell)
}

// CellMergeMsg cell merge message
// 表示单元格合并的信息，包括起始和结束的单元格引用。
type CellMergeMsg struct {
	start string
	end   string
}

// GetStartAxis returns the top left cell reference of merged range, for
// example: "C2". 返回合并区域的左上角单元格引用。
func (c *CellMergeMsg) GetStartAxis() string {
	return c.start
}

// GetEndAxis returns the bottom right cell reference of merged range, for
// example: "D4". 返回合并区域的右下角单元格引用。
func (c *CellMergeMsg) GetEndAxis() string {
	return c.end
}
