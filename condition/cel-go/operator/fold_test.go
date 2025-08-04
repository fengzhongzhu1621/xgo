package operator

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestFolding(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("name", cel.StringType),
	)

	ast, _ := env.Compile(`
		[1, 1 + 2, 1 + (2 + 3)]
    `)
	pro, _ := env.Program(ast)

	out, _, err := pro.Eval(map[string]interface{}{})

	if err != nil {
		fmt.Printf("evaluation failed: %v", err)
	} else {
		fmt.Println("评估结果:", out) // [1, 3, 6]
	}
}
