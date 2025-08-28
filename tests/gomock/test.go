// gomock 是一个 Go 语言的测试框架，用于编写单元测试时模拟和测试依赖于外部服务的代码。
// 它允许你创建模拟对象（Mock Objects），这些对象可以预设期望的行为，以便在测试时模拟外部依赖，
// 通常使用它对代码中的那些接口类型进行mock。
//
// 安装 mockgen 工具
// go install go.uber.org/mock/mockgen@latest
// mockgen -version
// go get go.uber.org/mock/gomock
//
// 在使用 mockgen 生成模拟对象（Mock Objects）时，通常需要指定三个主要参数：
// mockgen -source 需要mock的文件名 -destination 生成的mock文件名 -package 生成mock文件的包名
// * source：这是你想要生成模拟对象的接口定义所在的文件路径。
// * destination：这是你想要生成模拟对象代码的目标路径。
// * package：这是生成代码的包名。
//
// mockgen -source test.go -destination mock_test.go -package mock

package mock

type MyInterForMock interface {
	GetName(id int) string
}

func GetUser(m MyInterForMock, id int) string {
	user := m.GetName(id)
	return user
}

type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
