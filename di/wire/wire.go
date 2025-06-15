//go:build wireinject
// +build wireinject

package wire

import "github.com/google/wire"

func InitMission(name string) Mission {
	wire.Build(NewMonster, NewPlayer, NewMission)
	return Mission{}
}

func InitializeServer(dsn string) *Server {
	wire.Build(NewDatabase, NewServer)
	return &Server{}
}

// ProviderSet 是所有依赖项的提供者集合
var ProviderSet = wire.NewSet(
	NewService,
	NewMySQLDatabase,
	wire.InterfaceValue(new(IDatabase), new(MySQLDatabase)), // 将 *MySQLDatabase 实现 IDatabase 接口
)

// InitializeService 使用 ProviderSet 创建 Service 实例
func InitializeService() (*Service, error) {
	wire.Build(ProviderSet)
	return &Service{}, nil
}
