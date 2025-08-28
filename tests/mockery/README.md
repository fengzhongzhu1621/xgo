# 简介

使用了
github.com/stretchr/testify/mock


# 引用
go get github.com/vektra/mockery/v2/.../

# go:generate
```go
package example

//go:generate mockery --name Speaker --case snake
type Speaker interface {
	Say(string) string
}
```

# 执行测试
```go
func TestSay(t *testing.T) {
	// 初始化 mock 对象
	mockCtrl := &mocks.Speaker{}

	// 设置匹配条件和返回，注意默认即匹配 AnyTimes
	mockCtrl.On("Say", "hello").Return("world").Times(2)

	var speaker Speaker = mockCtrl

	actual1 := speaker.Say("hello")
	assert.Equal(t, "world", actual1)
	actual2 := speaker.Say("hello")
	assert.Equal(t, "world", actual2)
}

```