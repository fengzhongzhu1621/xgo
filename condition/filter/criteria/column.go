package criteria

// FieldType define the table's field data type.
type FieldType string

// MapStringType 自定义 map[string]string
type MapStringType map[string]string

type NumericType int64

const (
	// Numeric means this field is numeric data type.
	Numeric FieldType = "numeric"

	// Time means this field is Time data type.
	Time FieldType = "time"

	// Timestamp means this field is timestamp data type.
	Timestamp FieldType = "timestamp"

	// Boolean means this field is boolean data type.
	Boolean FieldType = "bool"
	// String means this field is string data type.
	String FieldType = "string"
	// MapString means this field is map string type.there is a special map for
	// container scenarios, in which both key and value must be of string type,
	// such as label, taints, etc. At this time, these variables need to be set
	// to this type to facilitate judgment and verification.
	MapString FieldType = "mapString"
	// Array means this field is array data type.
	Array FieldType = "array"
	// Object means this field is object data type.
	Object FieldType = "object"
	// Enum means this field is enum type.
	Enum FieldType = "enum"
	// Note: subsequent support for other types can be added here.
	// after adding a type, pay attention to adding a verification
	// function for this type synchronously. special attention is
	// paid to whether the array elements also need to synchronize support for this type.
)
