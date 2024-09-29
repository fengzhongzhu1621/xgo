package goom

import (
	"testing"

	mocker "github.com/tencent/goom"
	"github.com/tencent/goom/arg"
)

// foo 函数定义如下
func foo(i int) int {
	//...
	return 0
}

// bar 多参数函数
func bar(i interface{}, j int) int {
	//...
	return 0
}

func TestMockFunc(t *testing.T) {
	// mock示例
	// 创建当前包的mocker
	mock := mocker.Create()

	// mock函数foo并设定返回值为1
	mock.Func(foo).Return(1)
	s.Equal(1, foo(0), "return result check")

	// 可搭配When使用: 参数匹配时返回指定值
	mock.Func(foo).When(1).Return(2)
	s.Equal(2, foo(1), "when result check")

	// 使用arg.In表达式,当参数为1、2时返回值为100
	mock.Func(foo).When(arg.In(1, 2)).Return(100)
	s.Equal(100, foo(1), "when in result check")
	s.Equal(100, foo(2), "when in result check")

	// 按顺序依次返回(等价于gomonkey的Sequence)
	mock.Func(foo).Returns(1, 2, 3)
	s.Equal(1, foo(0), "returns result check")
	s.Equal(2, foo(0), "returns result check")
	s.Equal(3, foo(0), "returns result check")

	// mock函数foo，使用Apply方法设置回调函数
	// 注意: Apply和直接使用Return都可以实现mock，两种方式二选一即可
	// Apply可以在桩函数内部实现自己的逻辑，比如根据不同参数返回不同值等等。
	mock.Func(foo).Apply(func(int) int {
		return 1
	})
	s.Equal(1, foo(0), "apply callback check")

	// 忽略第一个参数, 当第二个参数为1、2时返回值为100
	mock.Func(bar).When(arg.Any(), arg.In(1, 2)).Return(100)
	s.Equal(100, bar(-1, 1), "any param result check")
	s.Equal(100, bar(0, 1), "any param result check")
	s.Equal(100, bar(1, 2), "any param result check")
	s.Equal(100, bar(999, 2), "any param result check")
}
