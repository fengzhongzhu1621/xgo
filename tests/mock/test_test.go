package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGoMock(t *testing.T) {
	// 创建 gomock 控制器，用来记录后续操作信息
	// t: 一个实现了 TestReporter 接口的对象。通常情况下，你可以使用 testing.T 或 testing.B 作为参数
	// returns: 返回一个新的 *Controller 对象，用于管理模拟对象的生命周期。
	mockCtl := gomock.NewController(t)
	// 创建一个 MyInterForMock 接口的 mock 实例
	mockInter := NewMockMyInterForMock(mockCtl)
	// 设置 MyInterForMock.GetName() 的返回结果
	mockInter.EXPECT().GetName(1).Return("one")

	// 执行函数
	actual := GetUser(mockInter, 1)
	assert.Equal(t, "one", actual)

	// TODO: doesn't match the argument at index 0.
	// actual = GetUser(mockInter, 2)
	// assert.NotEqual(t, "one", actual)

	mockInter.EXPECT().GetName(2).Return("two")
	actual = GetUser(mockInter, 2)
	assert.Equal(t, "two", actual)
}
