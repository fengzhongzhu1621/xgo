package xerror

import (
	"fmt"
	"strings"
	"testing"

	"github.com/duke-git/lancet/v2/xerror"
	"github.com/stretchr/testify/assert"
)

// TestStackTrace 返回与pkg/errors兼容的堆栈跟踪。
// func (e *XError) StackTrace() StackTrace
func TestStackTrace(t *testing.T) {
	err := xerror.New("error")
	stacks := err.Stacks()

	fmt.Println(stacks[0].Func)
	fmt.Println(stacks[0].Line)

	containFile := strings.Contains(stacks[0].File, "xxx.go")
	assert.Equal(t, false, containFile)
}
