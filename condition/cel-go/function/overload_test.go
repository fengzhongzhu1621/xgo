package function

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func TestDoubleFunc(t *testing.T) {
	// 创建CEL环境并定义自定义函数
	env, err := cel.NewEnv(
		cel.Function("my_double", // 定义函数名
			cel.Overload("my_double_int", // 定义函数的一个重载版本
				[]*cel.Type{cel.IntType}, // 函数参数类型（整数）
				cel.IntType,              // 函数返回值类型（整数）
				cel.UnaryBinding( // 函数体实现
					func(value ref.Val) ref.Val {
						// 将CEL值转换为int类型
						if intVal, ok := value.(types.Int); ok {
							// 返回输入值的两倍
							return types.Int(intVal * 2)
						}
						// 如果参数类型不正确，返回错误
						return types.NewErr("invalid argument type")
					},
				),
			),
		),
	)
	if err != nil {
		log.Fatalf("created environment failed: %v", err) // 如果环境创建失败，记录错误并退出
	}

	// 定义CEL表达式：调用my_double(5)并与10比较
	expr := `my_double(5) == 10` // 调用my_double函数并将其结果和10进行相等比较

	// 编译CEL表达式
	ast, iss := env.Compile(expr) // 解析并检查表达式是否正确
	if iss.Err() != nil {
		log.Fatalf("compilation failed: %v", iss.Err()) // 如果编译失败，记录错误并退出
	}

	// 将AST编译为可执行程序
	program, err := env.Program(ast) // 执行表达式
	if err != nil {
		log.Fatalf("program creation failed: %v", err) // 如果程序创建失败，记录错误并退出
	}

	// 执行表达式求值
	out, _, err := program.Eval(map[string]any{}) // my_double函数我们直接传入5，所以这里传递参数时为空即可
	if err != nil {
		log.Fatalf("evaluation failed: %v", err) // 如果求值失败，记录错误并退出
	}

	// 打印表达式求值结果
	fmt.Println(out) // 输出结果为true
}

func TestExampleCustomGlobalFunction(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("i", cel.StringType),
		cel.Variable("you", cel.StringType),
		cel.Function("shake_hands",
			cel.Overload("shake_hands_string_string",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.StringType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					return types.String(
						fmt.Sprintf("%s and %s are shaking hands.\n", lhs, rhs))
				},
				),
			),
		),
	)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}
	// Check iss for error in both Parse and Check.
	ast, iss := env.Compile(`shake_hands(i,you)`)
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("Program creation error: %v\n", err)
	}

	out, _, err := prg.Eval(map[string]any{
		"i":   "CEL",
		"you": "world",
	})
	if err != nil {
		log.Fatalf("Evaluation error: %v\n", err)
	}

	fmt.Println(out)
	// Output:CEL and world are shaking hands.
}
