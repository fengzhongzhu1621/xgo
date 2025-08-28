package operator

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestAddJoinString(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("name", cel.StringType),
	)

	ast, _ := env.Compile(`
		"Hello world! I'm " + name + "."
    `)
	pro, _ := env.Program(ast)

	out, detail, _ := pro.Eval(map[string]interface{}{
		"name": "John",
	})

	fmt.Println("detail:", detail)
	fmt.Println("评估结果:", out)
}
