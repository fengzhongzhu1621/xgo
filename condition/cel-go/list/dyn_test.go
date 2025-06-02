package list

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestDyn(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		dyn([1, 2]) + [3.0, 4.0]
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // [1, 2, 3, 4]

	ast, _ = env.Compile(`
		{'a': dyn([1, 2]), 'b': x}
    `)
	pro, _ = env.Program(ast)
	out, _, _ = pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // [1, 2, 3, 4]

	ast, iss := env.Compile(`
		[1, 2] + [3.0, 4.0]
    `)
	if iss.Err() != nil {
		log.Fatalf("compilation failed: %v", iss)
	}
	//  ERROR: <input>:2:10: found no matching overload for '_+_' applied to '(list(int), list(double))'
	// |   [1, 2] + [3.0, 4.0]
	// | .........^
}

func TestDynUseInValue(t *testing.T) {
	env, _ := cel.NewEnv(cel.Variable("x", cel.IntType))

	ast, _ := env.Compile(`
		{'a': dyn([1, 2]), 'b': x}
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": 3,
	})
	fmt.Println("评估结果:", out) // {a: [1, 2], b: 3}
}
