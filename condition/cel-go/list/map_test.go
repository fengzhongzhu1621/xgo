package list

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestListMap(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		[1, 2, 3].map(i, [1, 2, 3].map(j, i * j))
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // [[1, 2, 3], [2, 4, 6], [3, 6, 9]]
}
