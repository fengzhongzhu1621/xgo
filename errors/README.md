# github.com/pkg/errors

在Go 1.13版本之前，用于扩展Go错误处理能力的一个第三方库。它的功能比标准库的pkg/errors更加丰富。可以保存和打印错误发生时的堆栈信息。

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


