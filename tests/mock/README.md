# mockgen

## 生成 mock 文件

```
mockgen -source=./person/male.go -destination=./mock/male_mock.go -package=mock
mockgen -source=./foo.go -destination=./mock_foo.go -package=gomock
```

- source：设置需要模拟（mock）的接口文件
- destination：设置 mock 文件输出的地方，若不设置则打印到标准输出中
- package：设置 mock 文件的包名，若不设置则为 mock_ 前缀加上文件名（如本文的包名会为 mock_person）

## 生成多个 mock 文件

```go
//go:generate mockgen -self_package github.com/lovoo/goka -package goka -destination mockstorage.go github.com/lovoo/goka/storage Storage
//go:generate mockgen -self_package github.com/lovoo/goka -package goka -destination mocks.go github.com/lovoo/goka TopicManager,Producer,Broker
//go:generate mockgen -self_package github.com/lovoo/goka -package goka -destination mockssarama.go github.com/IBM/sarama Client,ClusterAdmin
//go:generate mockgen -destination=mocks/logger.go -package=gocronmocks . Logger
//go:generate mockgen -destination=../mocks/app/mock_app.go -package=mock_app github.com/rudderlabs/rudder-server/app App
//go:generate mockgen -source=./internal/abs/abstraction.go -destination=./internal/mock/mock.go -package mock -mock_names IQueue=MockQueue,IDelegator=MockDelegator,IWorkerPool=MockWorkerPool
```

**修改 interface 方法**
```go
//go:generate mockgen -source=my_interface.go -destination=mocks/mock_my_interface.go -package=mocks

package mypackage

type MyInterface interface {
    DoSomething() error
}
```

- source 参数指定了包含接口定义的源文件
- destination 参数指定了生成的 mock 文件的位置和名称
- package 参数指定了生成的 mock 文件所在的包名。

**重新生成 mock 文件**
```bash
go generate [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
go generate ./...
```

# gomock

## 打桩(stubs)

### 输入
* Eq(value) 表示与 value 等价的值。
* Any() 可以用来表示任意的入参。
* Not(value) 用来表示非 value 以外的值。
* Nil() 表示 None 值

```go
m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
m.EXPECT().Get(gomock.Any()).Return(630, nil)
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil) 
m.EXPECT().Get(gomock.Nil()).Return(0, errors.New("nil")) 
```

### 返回值

* Return 返回确定的值
* Do Mock 方法被调用时，要执行的操作，忽略返回值。
* DoAndReturn 可以动态地控制返回值。

```go
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
m.EXPECT().Get(gomock.Any()).Do(func(key string) {
    t.Log(key)
})
m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
    if key == "Sam" {
        return 630, nil
    }
    return 0, errors.New("not exist")
}
```

### 调用次数(Times)

* Times() 断言 Mock 方法被调用的次数。
* MaxTimes() 最大次数。
* MinTimes() 最小次数。
* AnyTimes() 任意次数（包括 0 次）。
  
```go
mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)
mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
```

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)
	GetFromDB(m, "ABC")
	GetFromDB(m, "DEF")
}
```

### 调用顺序(InOrder)

```go
firstCall := mockSpider.EXPECT().Foo1()
mockRepo.EXPECT().Foo2().After(firstCall)
```

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
	o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)
	gomock.InOrder(o1, o2)
	GetFromDB(m, "Tom")
	GetFromDB(m, "Sam")
}

```

## 常用方法

|方法|作用|
|:----|:----|
|func (c Call) After(preReq Call)|*Call 指定执行顺序|
|func (c Call) AnyTimes() Call|允许调用次数为 0 次或更多次|
|func (c Call) Do(f interface{}) Call|指定匹配时执行的操作|
|func (c Call) DoAndReturn(f interface{}) Call|指定匹配是执行的操作，并且模拟返回值|
|func (c Call) MaxTimes(n int) Call|指定最大的调用次数为 n 次|
|func (c Call) MinTimes(n int) Call|指定最小的调用次数为 n 次|
|func (c Call) Return(rets …interface{}) Call|模拟返回值|
|func (c Call) SetArg(n int, value interface{}) Call|设置第n个参数的值为value|
|func (c *Call) String() string|返回其字符串形式|
|func (c Call) Times(n int) Call|指定调用次数为 n 次|
