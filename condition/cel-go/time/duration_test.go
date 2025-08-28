package time

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestDuration(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.StringType),
	)

	ast, _ := env.Compile(`
		duration(string(7 * 24) + 'h')
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": "world",
	})
	fmt.Println("评估结果:", out) // 168h0m0s
}

func TestTimestamp(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", cel.StringType),
	)

	ast, _ := env.Compile(`
		timestamp("1970-01-01T00:00:00Z")
    `)
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{
		"x": "world",
	})
	fmt.Println("评估结果:", out) // 1970-01-01 00:00:00 +0000 UTC
}
