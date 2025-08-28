package operator

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestQuestion(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)

	ast, _ := env.Compile(`
		true ? x + 1 : x + 2
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": 3,
	})
	fmt.Println("评估结果:", out) // 4
}
