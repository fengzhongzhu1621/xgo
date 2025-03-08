package excel

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fengzhongzhu1621/xgo/file/excel/style"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/xuri/excelize/v2"
)

const (
	excelSuffix  = ".xlsx"
	defaultSheet = "Sheet1"
	initIdx      = -1
	oneCellLen   = 1
)

type Excel struct {
	sync.RWMutex                                      // 保护对 Excel 实例的并发访问
	filePath        string                            // 存储 Excel 文件的路径
	file            *excelize.File                    // 持有一个 excelize.File 实例，用于操作 Excel 文件。
	writers         map[string]*excelize.StreamWriter // 管理多个 StreamWriter，可用于处理大文件的流式写入。
	hasDefaultSheet bool                              // 标识是否存在默认的工作表
}

// OperatorFunc 用于创建 Excel 实例时进行自定义初始化
type OperatorFunc func(excel *Excel) error

func NewExcel(opts ...OperatorFunc) (*Excel, error) {
	excel := &Excel{
		writers: make(map[string]*excelize.StreamWriter),
	}
	for _, opt := range opts {
		if err := opt(excel); err != nil {
			return nil, err
		}
	}

	return excel, nil
}

func WithFilePath(filePath string) OperatorFunc {
	return func(excel *Excel) error {
		excel.Lock()
		defer excel.Unlock()

		if !strings.HasSuffix(filePath, excelSuffix) {
			filePath = filePath + excelSuffix
		}

		excel.filePath = filePath
		return nil
	}
}

func WithKeepDefaultSheet() OperatorFunc {
	return func(excel *Excel) error {
		excel.Lock()
		defer excel.Unlock()

		excel.hasDefaultSheet = true

		return nil
	}
}

func WithOpenOrCreate() OperatorFunc {
	return func(excel *Excel) error {
		excel.Lock()
		defer excel.Unlock()

		if excel.filePath == "" {
			return errors.New("excel filePath can not be empty")
		}

		// 创建目录
		dirPath := filepath.Dir(excel.filePath)
		if _, err := os.Stat(dirPath); err != nil {
			if err := os.MkdirAll(dirPath, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
		}

		// 判断是否文件
		if !validator.IsFile(excel.filePath) {
			excel.file = excelize.NewFile()
			return nil
		}

		// 打开文件
		var err error
		excel.file, err = excelize.OpenFile(excel.filePath)
		if err != nil {
			return err
		}

		return nil
	}
}

func (excel *Excel) CreateSheet(sheet string) error {
	excel.Lock()
	defer excel.Unlock()
	if _, err := excel.file.NewSheet(sheet); err != nil {
		return err
	}

	return nil
}

func (excel *Excel) DeleteSheet(sheet string) error {
	excel.Lock()
	defer excel.Unlock()

	return excel.deleteSheet(sheet)
}

func (excel *Excel) deleteSheet(sheet string) error {
	return excel.file.DeleteSheet(sheet)
}

func (excel *Excel) SetAllColsWidth(sheet string, width float64) error {
	return excel.SetColWidth(sheet, excelize.MinColumns, excelize.MaxColumns, width)
}

func (excel *Excel) SetColWidth(sheet string, startCol, endCol int, width float64) error {
	excel.Lock()
	defer excel.Unlock()

	var err error
	if excel.writers[sheet] == nil {
		excel.writers[sheet], err = excel.file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	}

	if err := excel.writers[sheet].SetColWidth(startCol, endCol, width); err != nil {
		return err
	}

	return nil
}

// NewReader create io stream reader
func (excel *Excel) NewReader(sheet string) (*Reader, error) {
	excel.RLock()
	defer excel.RUnlock()

	if excel.file == nil {
		return nil, fmt.Errorf("excel file has not been created yet")
	}

	rows, err := excel.file.Rows(sheet)
	if err != nil {
		return nil, err
	}

	return &Reader{rows: rows, curIdx: initIdx}, nil
}

func (excel *Excel) MergeCell(sheet string, hCell, vCell string) error {
	excel.Lock()
	defer excel.Unlock()

	return excel.mergeCell(sheet, hCell, vCell)
}

func (excel *Excel) mergeCell(sheet string, hCell, vCell string) error {
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	var err error
	if excel.writers[sheet] == nil {
		excel.writers[sheet], err = excel.file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	}

	return excel.writers[sheet].MergeCell(hCell, vCell)
}

func (excel *Excel) Flush(sheets []string) error {
	excel.Lock()
	defer excel.Unlock()

	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	for _, sheet := range sheets {
		if excel.writers[sheet] == nil {
			continue
		}

		if err := excel.writers[sheet].Flush(); err != nil {
			return err
		}

		delete(excel.writers, sheet)
	}

	return nil
}

func (excel *Excel) save() error {
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	return excel.file.SaveAs(excel.filePath)
}

func (excel *Excel) Close() (err error) {
	excel.Lock()
	defer excel.Unlock()

	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	if !excel.hasDefaultSheet {
		if err = excel.deleteSheet(defaultSheet); err != nil {
			return
		}

		if err = excel.save(); err != nil {
			return
		}
	}

	defer func() {
		if err = excel.file.Close(); err != nil {
			return
		}
	}()

	return nil
}

// Clean 删除 excel文件
func (excel *Excel) Clean() error {
	excel.Lock()
	defer excel.Unlock()

	if err := os.Remove(excel.filePath); err != nil {
		return err
	}

	return nil
}

func (excel *Excel) StreamingWrite(sheet string, startIdx int, data [][]Cell) error {
	excel.Lock()
	defer excel.Unlock()
	if excel.file == nil {
		return fmt.Errorf("excel file has not been created yet")
	}

	// 初始化 StreamWriter 如果尚未初始化
	var err error
	if excel.writers[sheet] == nil {
		excel.writers[sheet], err = excel.file.NewStreamWriter(sheet)
		if err != nil {
			return err
		}
	}

	// 验证起始行索引
	startIdx++
	if startIdx < rowStartIdx {
		return fmt.Errorf("row start index is invalid, val: %d", startIdx)
	}

	for i := 0; i < len(data); i++ {
		// 列和行索引转换为单元格名称（如 "A1"），始终从第1行开始
		firstCell, err := excelize.CoordinatesToCellName(colStartIdx, startIdx+i)
		if err != nil {
			return err
		}

		cells := make([]interface{}, len(data[i]))
		for idx, val := range data[i] {
			cells[idx] = val.transfer()
		}

		// 设置行数据
		if err := excel.writers[sheet].SetRow(firstCell, cells); err != nil {
			return err
		}
	}

	return nil
}

func (excel *Excel) StreamingRead(sheet string) (result [][]string, err error) {
	excel.RLock()
	defer excel.RUnlock()
	if excel.file == nil {
		return nil, fmt.Errorf("excel file has not been created yet")
	}

	rows, err := excel.file.Rows(sheet)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			return
		}
	}()

	for rows.Next() {
		var rowData []string
		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		for _, cell := range cols {
			rowData = append(rowData, cell)
		}

		result = append(result, rowData)
	}

	return result, nil
}

func (excel *Excel) Save() error {
	excel.Lock()
	defer excel.Unlock()

	return excel.save()
}

func (excel *Excel) AddValidation(sheet string, param *ValidationParam) error {
	excel.Lock()
	defer excel.Unlock()

	validation, err := NewDataValidation(param)
	if err != nil {
		return err
	}

	if err := excel.file.AddDataValidation(sheet, validation); err != nil {
		return err
	}

	return nil
}

func (excel *Excel) NewStyle(style *style.Style) (int, error) {
	excelStyle, err := style.Convert()
	if err != nil {
		return 0, err
	}

	return excel.file.NewStyle(excelStyle)
}

func (excel *Excel) MergeSameRowCell(sheet string, colIdx, rowIdx, length int) error {

	if length == oneCellLen {
		return nil
	}

	hCell, err := GetCellIdx(colIdx, rowIdx)
	if err != nil {
		return err
	}

	vCell, err := GetCellIdx(colIdx+length-1, rowIdx)
	if err != nil {
		return err
	}

	excel.Lock()
	defer excel.Unlock()

	if err := excel.mergeCell(sheet, hCell, vCell); err != nil {
		return err
	}

	return nil
}

func (excel *Excel) MergeSameColCell(sheet string, colIdx, rowIdx, height int) error {
	if height == oneCellLen {
		return nil
	}

	hCell, err := GetCellIdx(colIdx, rowIdx)
	if err != nil {
		return err
	}

	vCell, err := GetCellIdx(colIdx, rowIdx+height-1)
	if err != nil {
		return err
	}

	excel.Lock()
	defer excel.Unlock()

	if err := excel.mergeCell(sheet, hCell, vCell); err != nil {
		return err
	}

	return nil
}

func (excel *Excel) GetMergeCellMsg(sheet string) ([]CellMergeMsg, error) {
	if excel.file == nil {
		return nil, fmt.Errorf("excel file has not been created yet")
	}

	cells, err := excel.file.GetMergeCells(sheet)
	if err != nil {
		return nil, fmt.Errorf("get merge cell failed, sheet: %s, err: %v", sheet, err)
	}

	result := make([]CellMergeMsg, len(cells))
	for idx, cell := range cells {
		start := cell.GetStartAxis()
		end := cell.GetEndAxis()

		result[idx] = CellMergeMsg{start: start, end: end}
	}

	return result, nil
}

func (excel *Excel) IsSheetExist(sheet string) (bool, error) {
	excel.RLock()
	defer excel.RUnlock()

	if excel.file == nil {
		return false, fmt.Errorf("excel file has not been created yet")
	}

	idx, err := excel.file.GetSheetIndex(sheet)
	if err != nil {
		return false, err
	}

	return idx > -1, nil
}
