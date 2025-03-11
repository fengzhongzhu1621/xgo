package xerror

import (
	"context"
	"testing"

	"github.com/duke-git/lancet/v2/xerror"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// TestTryCatch 简单模拟Java风格的try-catch。它与Go的错误处理哲学不符。建议谨慎使用。
// func NewTryCatch(ctx context.Context) *TryCatch
// func (tc *TryCatch) Try(tryFunc func(ctx context.Context) error) *TryCatch
// func (tc *TryCatch) Catch(catchFunc func(ctx context.Context, err error)) *TryCatch
// func (tc *TryCatch) Finally(finallyFunc func(ctx context.Context)) *TryCatch
// func (tc *TryCatch) Do()
func TestTryCatch(t *testing.T) {
	calledFinally := false
	calledCatch := false

	tc := xerror.NewTryCatch(context.Background())

	tc.Try(func(ctx context.Context) error {
		// 抛出异常
		return errors.New("error message    ")
	}).Catch(func(ctx context.Context, err error) {
		// 捕获异常
		calledCatch = true
		// Error in try block at /path/xxx.go:{line_number} - Cause: error message
		// fmt.Println(err.Error())
	}).Finally(func(ctx context.Context) {
		calledFinally = true
	}).Do()

	assert.Equal(t, true, calledCatch)
	assert.Equal(t, true, calledFinally)
}
