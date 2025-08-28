package celgo

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestNull(t *testing.T) {
	env, _ := cel.NewEnv()

	ast, iss := env.Compile("null")
	if iss.Err() != nil {
		t.Errorf("compile failed: %v", iss.Err())
		return
	}
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})

	fmt.Println("评估结果:", out) // 0
}
