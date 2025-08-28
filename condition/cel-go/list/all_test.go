package list

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestAll(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		["ab", "abc", "abcd", "abcde"].all(x, x.size() > 2)
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // false

	ast2, _ := env.Compile(`
		["a", "b", "ab", "abc"].all(x, x.size() >=1 )
    `)
	pro2, _ := env.Program(ast2)
	out2, _, _ := pro2.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out2) // true
}
