package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestGoMock(t *testing.T) {
	// 创建 gomock 控制器，用来记录后续操作信息
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
