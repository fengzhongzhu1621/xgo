package operator

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestAndOperator(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
	)

	ast, _ := env.Compile(`
        name == 'John' && age > 30
    `)
	pro, _ := env.Program(ast)

	out, _, err := pro.Eval(map[string]interface{}{
		"name": "John",
		"age":  31,
	})

	if err != nil {
		fmt.Printf("evaluation failed: %v", err)
	} else {
		fmt.Println("评估结果:", out) // [1, 3, 6]
	}
}
