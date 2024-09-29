package assert

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// 定义结构体，继承 suite.Suite，即可拥有 Suite 的接口方法
type exampleTestSuite struct {
	suite.Suite

	// 放置公共的变量或者对象，不用每个 Test 方法都重复创建
	a int
}

// 初始化操作，每个 Test 函数执行前初始化
func (s *exampleTestSuite) SetupTest() {
	s.a = 1
}

// 单元测试用例
func (s *exampleTestSuite) TestExample1() {
	s.Equal(1, s.a)
}

func (s *exampleTestSuite) TestExample2() {
	s.Equal(1, s.a)
}

// 测试套件的执行，执行 Suite 对象的所有 Test 方法
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(exampleTestSuite))
}
