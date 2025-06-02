package celgo

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

// 表达式是一个字符串，没有参数
func TestCompileOnlyStr(t *testing.T) {
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}

	// 字符串解析不报错
	_, iss := env.Compile("a")
	if iss.Err() != nil {
		fmt.Printf("compile failed: %v", iss.Err())
		return
	}

	// compile failed: ERROR: <input>:1:1: undeclared reference to 'a' (in container '')
	// | a
}

// 表达式是一个字符串，但是非法字符
func TestCompileInvalidChar(t *testing.T) {
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}

	_, iss := env.Compile("-")
	if iss.Err() != nil {
		fmt.Printf("compile failed: %v", iss.Err())
		return
	}
	// compile failed: ERROR: <input>:1:2: Syntax error: no viable alternative at input '-'
	//  | -
	//  | .^
	// ERROR: <input>:1:2: Syntax error: mismatched input '<EOF>' expecting {'[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}
	//  | -
}
