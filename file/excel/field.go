package excel

// FieldType field type
type FieldType string

const (
	Decimal FieldType = "decimal"
	Bool    FieldType = "bool"
	// Enum 谨慎使用，excel本身限制单元格下拉列表的总大小不超过255字符，如果超过，会报错；
	// 如果下拉列表总大小需要超过255字符，可以使用Ref类型引用另一个sheet的一列值作为下拉列表
	Enum FieldType = "enum"
	Ref  FieldType = "ref" // 类型引用
)
