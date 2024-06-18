package table

import (
	"fmt"
	"testing"

	"github.com/bndr/gotabulate"
)

func TestRenderTable(t *testing.T) {
	// []interface{} 是一个切片类型，它可以存储任意类型的值。这意味着你可以在一个 []interface{} 切片中存储不同类型的元素，
	// 如整数、浮点数、字符串、布尔值、结构体等。
	// 这使得 []interface{} 成为处理不确定类型数据的强大工具。
	row_1 := []interface{}{"a", 10, "fine"}
	row_2 := []interface{}{"b", 20, "fine"}

	// 创建二维切片
	table := gotabulate.Create([][]interface{}{row_1, row_2})

	// 设置表头
	table.SetHeaders([]string{"name", "age", "status"})

	// 设置空字符串的显示方式
	table.SetEmptyString("None")

	// 设置单元格内容对齐方式
	table.SetAlign("center")

	// 打印表格
	fmt.Println(table.Render("grid"))
}
