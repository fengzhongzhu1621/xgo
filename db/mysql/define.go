package mysql

type Database struct {
	ID       string
	Host     string
	Port     int
	User     string
	Password string
	Name     string

	// 最大连接数
	MaxOpenConns int
	// 最大空闲连接数
	MaxIdleConns int
	// 单个连接的最大生命周期
	ConnMaxLifetimeSecond int

	// 是否打印sql 语句
	DebugMode bool
}

// CommonQueryConditionMap 通用查询条件
type CommonQueryConditionMap map[string]interface{}
