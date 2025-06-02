package operator

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestGteOperator(t *testing.T) {
	// 创建环境，定义变量
	env, err := cel.NewEnv(
		cel.Variable("age", cel.IntType), // 定义变量age，类型为整数
	)
	if err != nil {
		log.Fatalf("created environment failed: %v", err) // 如果创建环境失败，记录错误并退出
	}

	// 定义一个CEL表达式字符串
	expr := "age >= 18" // 定义表达式

	// 将表达式编译为抽象语法树(AST)
	// ast：编译后的抽象语法树
	// iss：包含编译过程中可能出现的问题的IssueSet
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		log.Fatalf("compilation failed: %v", iss)
	}

	// 将AST编译为可执行的程序
	// pro：可执行的CEL程序
	// err：程序创建过程中的错误
	pro, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program creation failed: %v", err)
	}

	// 执行CEL程序并求值，将age设置为18，并将其传入到表达式中
	// out：表达式的求值结果
	// _：忽略第二个返回值（通常是类型信息）
	// err：求值过程中的错误
	out, _, err := pro.Eval(map[string]any{
		"age": 18,
	})
	if err != nil {
		log.Fatalf("evaluation failed: %v", err)
	}

	// 输出结果
	fmt.Println(out) // 输出结果为true
}

func TestGte2(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		100+200 > 300
    `)
	pro, _ := env.Program(ast)

	out, detail, _ := pro.Eval(map[string]interface{}{})

	fmt.Println("detail:", detail)
	fmt.Println("评估结果:", out) // false
}
