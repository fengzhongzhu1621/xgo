package xerror

import (
	"testing"

	"github.com/duke-git/lancet/v2/xerror"
	"github.com/stretchr/testify/assert"
)

// TestInAndIs godoc
// func (e *XError) Id(id string) *XError 设置XError对象id以在XError.Is中检查是否相等。
// func (e *XError) Is(target error) bool 检查目标错误是否为XError，并且这两个错误的Error.id是否匹配。
func TestInAndIs(t *testing.T) {
	err1 := xerror.New("error").Id("e001")
	err2 := xerror.New("error").Id("e001")
	err3 := xerror.New("error").Id("e003")

	equal := err1.Is(err2)
	notEqual := err1.Is(err3)

	assert.Equal(t, true, equal)
	assert.Equal(t, false, notEqual)
}
