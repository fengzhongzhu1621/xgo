package celgo

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

// 表达式是一个数字，没有参数
func TestCompileOnlyNumber(t *testing.T) {
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}

	// 字符串解析不报错
	ast, iss := env.Compile("123")
	if iss.Err() != nil {
		fmt.Printf("compile failed: %v", iss.Err())
		return
	}

	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})

	fmt.Println("评估结果:", out) // 123
}
