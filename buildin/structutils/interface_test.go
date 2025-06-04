package structutils

// 在结构体中定义接口
// 在 Go 语言中，在结构体中定义接口类型的字段是一种非常常见的做法，它主要体现了 Go 的「面向接口编程」思想，能够带来解耦、可测试、易扩展等好处。
//
// | 作用               | 说明                                                                 |
// |--------------------|----------------------------------------------------------------------|
// | 实现依赖抽象       | 结构体不依赖具体实现，而是依赖一个接口，符合「依赖倒置原则」         |
// | 提高可测试性       | 可以注入 Mock 或 Stub 实现，便于单元测试                             |
// | 增强灵活性与可扩展性 | 可以在运行时替换不同的实现，而不影响使用该结构体的代码               |
// | 解耦代码           | 结构体与具体实现解耦，便于模块化设计和团队协作                       |
//

import (
	"fmt"
	"testing"
)

// 1️⃣ 定义一个数据库操作的接口
type IUserRepository interface {
	GetUserByID(id int) (*User3, error)
}

// 2️⃣ 定义 User 结构体
type User3 struct {
	ID   int
	Name string
}

// /////////////////////////////////////////////////////////////////
// 3️⃣ 定义一个具体的数据库实现（比如 MySQL）
type MySQLUserRepository struct{}

func (m *MySQLUserRepository) GetUserByID(id int) (*User3, error) {
	// 模拟从数据库查询
	return &User3{ID: id, Name: "User from MySQL"}, nil
}

// /////////////////////////////////////////////////////////////////
// 4️⃣ 定义 UserService，**结构体中包含接口类型的字段**
// UserService 不关心数据是从 MySQL、PostgreSQL、MongoDB 还是 Mock 中获取的；
// 它只依赖 UserRepository 接口，符合「面向接口编程」原则；
// 你可以轻松替换 repo 的实现，而不需要修改 UserService 的代码。
type UserService struct {
	repo IUserRepository // 这里就是「在结构体中定义接口」，依赖注入点 ✅
}

// UserService 的业务方法
func (s *UserService) GetUserName(id int) (string, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

// /////////////////////////////////////////////////////////////////
// 5️⃣ 使用示例
func TestStructIncludeInterface(t *testing.T) {
	// 使用真实数据库实现
	mysqlRepo := &MySQLUserRepository{}
	userService := &UserService{repo: mysqlRepo} // 通过构造函数注入依赖

	name, err := userService.GetUserName(1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("User name:", name) // 输出：User name: User from MySQL
}
