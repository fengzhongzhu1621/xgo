package ginkgo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// 运行 go test 或 ginkgo 的时候 Go 测试执行器会执行这个函数
func TestGinkgo(t *testing.T) {
	// 将Ginkgo的Fail函数传递给Gomega，Fail函数用于标记测试失败，这是Ginkgo和Gomega唯一的交互点
	// 如果Gomega断言失败，就会调用Fail进行处理
	RegisterFailHandler(Fail)
	// 启动测试套件
	// 告诉 Ginkgo 开始这个测试套件。如果任意 specs（说明）失败了，Ginkgo 会自动使 testing.T 失败。
	RunSpecs(t, "Ginkgo Suite")
}
