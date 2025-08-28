package dict

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestDictGetAttr(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.StringType),
	)

	ast, _ := env.Compile(`
		{'hello': 'world'}.hello == x


    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": "world",
	})
	fmt.Println("评估结果:", out) // true

	out2, _, _ := pro.Eval(map[string]interface{}{
		"x": "unknown",
	})
	fmt.Println("评估结果:", out2) // false
}

func TestStringType(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.MapType(cel.StringType, cel.StringType)),
	)

	ast, _ := env.Compile(`
		x.a == "1"
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": map[string]string{
			"a": "1",
			"b": "2",
		},
	})
	fmt.Println("评估结果:", out) // true
}
