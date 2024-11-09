package ginkgo

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func mockInputData() ([]Gopher, error) {
	inputData := []Gopher{
		{
			Name:   "菜刀",
			Gender: "男",
			Age:    18,
		},
		{
			Name:   "小西瓜",
			Gender: "女",
			Age:    19,
		},
		{
			Name:   "机器铃砍菜刀",
			Gender: "男",
			Age:    17,
		},
		{
			Name:   "小菜刀",
			Gender: "男",
			Age:    20,
		},
	}
	return inputData, nil
}

var _ = Describe("validate", func() {

	// 每个测试例执行前执行该段代码
	BeforeEach(func() {
		By("当测试不通过时，我会在这里打印一个消息 【BeforeEach】")
	})

	// 构造 mock 数据
	inputData, err := mockInputData()

	Describe("校验输入数据", func() {

		Context("当获取数据没有错误发生时", func() {
			It("它应该是接收数据成功了的", func() {
				// gomega.Expect(err)：这部分告诉 gomega 库我们想要对变量 err 进行断言。
				// .Should(gomega.BeNil())：这部分是一个断言，检查 err 是否为 nil。如果 err 不是 nil，测试将会失败，并且会输出一个描述性的错误信息。
				gomega.Expect(err).Should(gomega.BeNil())
			})
		})

		Context("当获取的数据校验失败时", func() {
			It("当数据校验返回错误为：名字太短，不能小于3 时", func() {
				gomega.Expect(Validate(inputData[0])).Should(gomega.MatchError("名字太短，不能小于3"))
			})

			It("当数据校验返回错误为：只要男的 时", func() {
				gomega.Expect(Validate(inputData[1])).Should(gomega.MatchError("只要男的"))
			})

			It("当数据校验返回错误为：岁数太小，不能小于18 时", func() {
				gomega.Expect(Validate(inputData[2])).Should(gomega.MatchError("岁数太小，不能小于18"))
			})
		})

		Context("当获取的数据校验成功时", func() {
			It("通过了数据校验", func() {
				gomega.Expect(Validate(inputData[3])).Should(gomega.BeNil())
			})
		})
	})

	// 每个测试例执行后执行该段代码
	AfterEach(func() {
		By("当测试不通过时，我会在这里打印一个消息 【AfterEach】")
	})
})
