package mock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	// 使用defer关键字确保在测试结束时调用Finish()方法。
	// 这个方法会检查所有的期望是否都已经满足，如果没有，则测试失败。
	// 不要忘记在测试结束时调用ctrl.Finish()来验证所有的期望是否都已满足
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))
	// 如果你不关心传入的具体参数，可以使用gomock.Any()来匹配任意参数
	m.EXPECT().Get(gomock.Any()).Return(100, errors.New("not exist"))

	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}

func TestGetFromDB2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)
	GetFromDB(m, "ABC")
	GetFromDB(m, "DEF")
	// 第三次失败 has already been called the max number of times
	// GetFromDB(m, "EFG")
}

func TestGetFromDB3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
	o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)

	// 用于指定模拟对象方法调用顺序的一个函数。
	// 这个函数确保在测试中，o1 和 o2 表示的方法调用会按照指定的顺序发生。
	gomock.InOrder(o1, o2)

	GetFromDB(m, "Tom")
	GetFromDB(m, "Sam")
}
