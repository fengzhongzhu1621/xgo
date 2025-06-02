package list

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestSize(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		[1, 1 + 1, 1 + 2, 2 + 3].size()
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // 4
}
