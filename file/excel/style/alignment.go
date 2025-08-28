package style

// Alignment directly maps the alignment settings of the cells.
// 单元格的对齐方式设置
type Alignment struct {
	Horizontal      string // 水平对齐方式（如 "left", "center", "right"）
	Indent          int    // 缩进级别
	JustifyLastLine bool   // 是否仅对最后一行进行对齐
	ReadingOrder    uint64 // v阅读顺序（具体含义依赖于 excelize 库的实现）
	RelativeIndent  int    // 相对缩进
	ShrinkToFit     bool   // 是否缩小字体以适应单元格宽度
	TextRotation    int    // 文本旋转角度
	Vertical        string // 垂直对齐方式（如 "top", "center", "bottom"）
	WrapText        bool   // 是否自动换行
}
