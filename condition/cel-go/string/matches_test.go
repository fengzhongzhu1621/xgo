package string

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

func TestMatches(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Variable("x", types.StringType),
		cel.ASTValidators(cel.ValidateRegexLiterals()),
	)

	ast, _ := env.Compile(`
		'hello'.matches('el*')
    `)

	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})
	fmt.Println("评估结果:", out) // true
}
