package excel

import (
	"sync"

	"github.com/xuri/excelize/v2"
)

// Reader 以流的方式读取 Excel 文件中的数据
type Reader struct {
	sync.RWMutex                // 提供读写锁机制，允许多个并发读取操作，但写入操作是互斥的。
	rows         *excelize.Rows // 行索引，表示 Excel 文件中的一组行数据，通常通过 Sheet.Rows() 方法获取。
	curIdx       int            // 当前行索引，用于跟踪读取进度。默认 -1
}

// Next 移动到下一行，并返回是否存在该行
func (r *Reader) Next() bool {
	r.RLock()
	defer r.RUnlock()

	r.curIdx++
	return r.rows.Next()
}

// CurRow 获取当前行的所有列值
func (r *Reader) CurRow() ([]string, error) {
	r.Lock()
	defer r.Unlock()

	columns, err := r.rows.Columns()
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// Close 关闭行读取器，释放资源
func (r *Reader) Close() error {
	r.Lock()
	defer r.Unlock()

	return r.rows.Close()
}

// GetCurIdx 获取当前的行索引
func (r *Reader) GetCurIdx() int {
	r.RLock()
	defer r.RUnlock()

	return r.curIdx
}
