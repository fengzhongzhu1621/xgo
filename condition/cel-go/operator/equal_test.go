package operator

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestEqual(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)

	ast, _ := env.Compile(`
		1 + x + 2 == 2 + x + 1
    `)
	pro, _ := env.Program(ast)

	out, _, _ := pro.Eval(map[string]interface{}{
		"x": 3,
	})

	fmt.Println("评估结果:", out) // true
}
