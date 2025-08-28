package operator

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestInTrue(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, _ := env.Compile(`
		6 in [1, 1 + 2, 1 + (2 + 3)]
    `)

	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})

	fmt.Println("评估结果:", out) // true
}

func TestKeyIsVar(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, iss := env.Compile(`
		x in [1, 1 + 2, 1 + (2 + 3)]
    `)
	if iss.Err() != nil {
		// 环境没有声明 x 变量
		fmt.Printf("compile failed: %v", iss.Err())
	}
	// compile failed: ERROR: <input>:2:3: undeclared reference to 'x' (in container '')
	// |   x in [1, 1 + 2, 1 + (2 + 3)]
	// | ..^PASS
	//

	env, _ = cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)
	ast, iss = env.Compile(`
		x in [1, 1 + 2, 1 + (2 + 3)]
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": 1,
	})
	fmt.Println("评估结果:", out) // true

	out2, _, _ := pro.Eval(map[string]interface{}{
		"x": "1",
	})
	fmt.Println("评估结果:", out2) // false
}

func TestValueContainVar(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)

	ast, _ := env.Compile(`
		1 in [1, x + 2, 1 + (2 + 3)]
    `)

	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": 1,
	})

	fmt.Println("评估结果:", out) // true
}

func TestInMap(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("name", cel.StringType),
	)

	ast, _ := env.Compile(`
		name in {'hello': false, 'world': true}
    `)

	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"name": "hello",
	})

	fmt.Println("评估结果:", out) // true
}
