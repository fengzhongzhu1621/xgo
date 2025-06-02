package function

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func TestExampleCustomInstanceFunction(t *testing.T) {
	env, err := cel.NewEnv(cel.Lib(customLib{}))
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}
	// Check iss for error in both Parse and Check.
	ast, iss := env.Compile(`i.greet(you)`)
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("Program creation error: %v\n", err)
	}

	out, _, err := prg.Eval(map[string]any{
		// Native values are converted to CEL values under the covers.
		"i": "CEL",
		// Values may also be lazily supplied.
		"you": func() ref.Val { return types.String("world") },
	})
	if err != nil {
		log.Fatalf("Evaluation error: %v\n", err)
	}

	fmt.Println(out)
	// Output:Hello world! Nice to meet you, I'm CEL.
}

// customLib 是一个自定义CEL库
type customLib struct{}

// CompileOptions 返回编译选项
func (customLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("i", cel.StringType),
		cel.Variable("you", cel.StringType),
		cel.Function("greet",
			cel.MemberOverload("string_greet_string",
				[]*cel.Type{cel.StringType, cel.StringType}, // 函数参数类型：两个字符串
				cel.StringType, // 函数返回值类型：字符串
				// 函数体实现
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					// 将CEL值转换为字符串
					return types.String(
						fmt.Sprintf("Hello %s! Nice to meet you, I'm %s.\n", rhs, lhs))
				}),
			),
		),
	}
}

// ProgramOptions 返回程序选项
func (customLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
