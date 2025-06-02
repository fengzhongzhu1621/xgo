package function

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestType(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, iss := env.Compile("type(1)")
	if iss.Err() != nil {
		t.Errorf("compile failed: %v", iss.Err())
		return
	}
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})

	fmt.Println("评估结果:", out) // int
}
