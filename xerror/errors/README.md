# github.com/pkg/errors

在Go 1.13版本之前，用于扩展Go错误处理能力的一个第三方库。它的功能比标准库的pkg/errors更加丰富。可以保存和打印错误发生时的堆栈信息。
Go 1.13 之后，标准库errors引入了Wrap和Is等方法来处理错误的包装和检查，部分功能和pkg/errors类似。
github.com/pkg/errors 的Wrap函数和标准库的 fmt.Errorf 的 %w 占位符有一些不同。如果用github.com/pkg/errors的代码中混用了 fmt.Errorf 的 % w，需要使用 errors.Unwrap 或 errors.Is、errors.As 来获取原始错误，而不能使用 github.com/pkg/errors 的Cause函数。

## New()

```
var ErrNameEmpty = errors.New("Name can't be empty!")
```

## errors.Wrap()

在已有错误基础上同时附加堆栈信息和新提示信息

* %+v 显示错误信息和堆栈信息，例如函数运行时的文件名、行数等信息
* 占位符 %s 则只显示错误信息

```go
err := errors.Wrap(originalError, "an error occurred")
fmt.Printf("%+v", err) // 打印错误和堆栈信息
```

## errors.WithMessage()

在已有错误基础上附加新提示信息

## errors.WithStack()

在已有错误基础上附加堆栈信息。

# pkg/errors
使用的Go版本在1.13及以上，那么你可以优先选择标准库的 pkg/errors


# 错误码的设计
* 每个错误码是由第一段模块号+错误分类+具体错号组成
* 错误码能够做到全局唯一
* 错误码不可随意追加

## 错误分类
### A 类（1）
用户行为错误，默认错误号是 10000

- 10500 参数非法(为空，或者基础校验不通过)
    - 10501 参数校验失败（page传了负数，手机号位数不对）
    - 10502 参数缺失（没传必传参数）
- 10800 没有业务权限
- 10900 业务操作限频
- 11000 获取用户登录态失败
- 11100 数据不存在异常

### B 类（2）
系统错误，默认错误号是 20000

- 20100 资源耗尽
- 20200 容灾功能被触发（系统限流、有损降级、过载保护）
- 20300 应用程序错误
  - 20301 业务序列化失败
  - 20302 业务未知异常，panic 等
  - 20303 业务异常场景，进入异常分支无法处理等
- 20500 分区表错误
- 20600 数据一致性异常

### C 类（3）
下游服务错误，默认错误号是 30000
- 30100 调用内部下游服务超时
- 30200 调用内部下游服务失败
- 30300 调用外部下游服务超时
- 30400 调用外部下游服务失败

### D 类（4）
依赖组件错误（如Redis、MySQL、ES等），默认错误号是 40000

- 40100 调用MySQL错误
- 40200 调用Kafka错误
- 40300 调用ES错误
- 40301 调用mongodb错误
