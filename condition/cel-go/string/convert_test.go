package string

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestConvertToString(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		string(7 * 24)
    `)

	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})

	fmt.Println("评估结果:", out) // 168
}
