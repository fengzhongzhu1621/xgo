package xerror

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/xerror"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// TestInfo 返回可以打印的xerror信息。
// func (e *XError) Info() *errInfo
func TestInfo(t *testing.T) {
	cause := errors.New("error")
	err := xerror.Wrap(cause, "invalid username").Id("e001").With("level", "high")

	errInfo := err.Info()

	assert.Equal(t, "e001", errInfo.Id)
	fmt.Println(errInfo.Cause) // error
	assert.Equal(t, "high", errInfo.Values["level"])
	assert.Equal(t, "invalid username", errInfo.Message)
}
