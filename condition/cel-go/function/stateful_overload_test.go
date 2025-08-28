package function

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type contextString string

func TestExample_statefulOverload(t *testing.T) {
	// makeFetch produces a consistent function signature with a different function
	// implementation depending on the provided context.
	// 根据上下文注册一个函数
	makeFetch := func(ctx any) cel.EnvOption {
		fn := func(arg ref.Val) ref.Val {
			return types.NewErr("stateful context not bound")
		}
		if ctx != nil {
			fn = func(resource ref.Val) ref.Val {
				return types.DefaultTypeAdapter.NativeToValue(
					// 从上下文获取数据
					ctx.(context.Context).Value(contextString(string(resource.(types.String)))),
				)
			}
		}
		return cel.Function("fetch",
			cel.Overload("fetch_string",
				[]*cel.Type{cel.StringType}, cel.StringType,
				cel.UnaryBinding(fn), // 专门用于处理一元操作符（即只需要一个参数的操作符）的绑定方式。
			),
		)
	}

	// The base environment declares the fetch function with a dummy binding that errors
	// if it is invoked without being replaced by a subsequent call to `baseEnv.Extend`
	baseEnv, err := cel.NewEnv(
		// Identifiers used within this expression.
		cel.Variable("resource", cel.StringType),
		// Function to fetch a resource.
		//    fetch(resource)
		makeFetch(nil),
	)
	if err != nil {
		log.Fatalf("environment creation error: %s\n", err)
	}

	// 编译表达式
	ast, iss := baseEnv.Compile("fetch('my-resource') == 'my-value'")
	if iss.Err() != nil {
		log.Fatalf("Compile() failed: %v", iss.Err())
	}

	// 使用 makeFetch(ctx) 扩展基础环境，提供有状态的 fetch 函数实现。
	// The runtime environment extends the base environment with a contextual binding for
	// the 'fetch' function.
	ctx := context.WithValue(context.TODO(), contextString("my-resource"), "my-value")
	runtimeEnv, err := baseEnv.Extend(makeFetch(ctx))
	if err != nil {
		log.Fatalf("baseEnv.Extend() failed with error: %s\n", err)
	}

	// 从扩展后的环境创建一个可执行程序 prg。
	prg, err := runtimeEnv.Program(ast)
	if err != nil {
		log.Fatalf("runtimeEnv.Program() error: %s\n", err)
	}

	// 求值编译后的表达式，没有传入任何变量
	out, _, err := prg.Eval(cel.NoVars())
	if err != nil {
		log.Fatalf("runtime error: %s\n", err)
	}

	fmt.Println(out)
	// Output:true
}
