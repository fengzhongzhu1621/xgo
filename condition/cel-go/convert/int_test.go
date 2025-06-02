package function

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestConvertToInt(t *testing.T) {
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}

	ast, iss := env.Compile("int(1)")
	if iss.Err() != nil {
		t.Errorf("compile failed: %v", iss.Err())
		return
	}
	pro, _ := env.Program(ast)
	out, _, _ := pro.Eval(map[string]interface{}{})

	if err != nil {
		fmt.Printf("evaluation failed: %v", err)
	} else {
		fmt.Println("评估结果:", out) // nil
	}
}
