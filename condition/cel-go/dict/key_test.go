package dict

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestDictIndex(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.StringType),
	)

	ast, _ := env.Compile(`
		{'hello': 'world'}['hello'] == x
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": "world",
	})
	fmt.Println("评估结果:", out) // true
}

func TestDictOrValue(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.StringType),
	)

	_, iss := env.Compile(`
		{'hello': 'world'}.?hello.orValue('default') == x
    `)
	if iss.Err() != nil {
		fmt.Printf("compile failed: %v", iss.Err())
	}
	// compile failed: ERROR: <input>:2:21: unsupported syntax '.?'
	// |   {'hello': 'world'}.?hello.orValue('default') == x
}

func TestDictOptional(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.StringType),
	)

	_, iss := env.Compile(`
		{'hello': optional.of('world')}['hello'] == x
    `)
	if iss.Err() != nil {
		fmt.Printf("compile failed: %v", iss.Err())
	}
	// ERROR: <input>:2:13: undeclared reference to 'optional' (in container '')
	//	 |   {'hello': optional.of('world')}['hello'] == x
	//	 | ............^
	//
	// ERROR: <input>:2:24: undeclared reference to 'of' (in container ”)
	//
	//	|   {'hello': optional.of('world')}['hello'] == x
}
