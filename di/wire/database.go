package wire

// IDatabase 定义数据库操作的接口
type IDatabase interface {
	Query() string
}

// ////////////////////////////////////////////////////////
// Database 定义数据库结构体
type Database struct {
	DSN string
}

// NewDatabase 创建一个新的数据库实例
func NewDatabase(dsn string) *Database {
	return &Database{DSN: dsn}
}

// ////////////////////////////////////////////////////////
type MySQLDatabase struct{}

func (m *MySQLDatabase) Query() string {
	return "Executing MySQL query"
}

// NewMySQLDatabase 是 MySQLDatabase 的提供者函数
func NewMySQLDatabase() *MySQLDatabase {
	return &MySQLDatabase{}
}
