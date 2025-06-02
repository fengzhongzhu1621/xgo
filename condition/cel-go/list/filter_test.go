package list

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestListFilter(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)

	ast, _ := env.Compile(`
		[1, 2, 3].map(i, [1, 2, 3].map(j, i * j).filter(k, k % 2 == 0))
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // [[2], [2, 4, 6], [6]]

	ast, _ = env.Compile(`
		[1, 2, 3].map(i, [1, 2, 3].map(j, i * j).filter(k, k % 2 == x))
    `)
	pro, _ = env.Program(ast)
	out, _, _ = pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // no such attribute(s): x
	out, _, _ = pro.Eval(map[string]interface{}{
		"x": 1,
	})
	fmt.Println("评估结果:", out) // [[1, 3], [], [3, 9]]
}

func TestListFilterHas(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		[{}, {"a": 1}, {"b": 2}].filter(m, has(m.a))
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // [{a: 1}]

	ast, _ = env.Compile(`
		[{}, {"a": 1}, {"b": 2}].filter(m, has({'a': true}.a))
    `)
	pro, _ = env.Program(ast)
	out, _, _ = pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // [{}, {a: 1}, {b: 2}]
}
