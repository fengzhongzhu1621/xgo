package function

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	is := assert.New(t)

	// no error
	{
		transaction := lo.NewTransaction[int]().
			Then(
				// 第一个Then块中，状态增加100，没有错误，所以状态变为 109
				func(state int) (int, error) {
					return state + 100, nil
				},
				// 状态减少100（这是为了回滚准备，但因为没有错误，所以实际上状态仍然是 109）。
				func(state int) int {
					return state - 100
				},
			).
			Then(
				// 状态再增加21，没有错误，所以最终状态为 130。
				func(state int) (int, error) {
					return state + 21, nil
				},
				func(state int) int {
					return state - 21
				},
			)

		state, err := transaction.Process(9)
		is.Equal(130, state)
		is.Equal(nil, err)
	}

	// with error
	{
		transaction := lo.NewTransaction[int]().
			Then(
				func(state int) (int, error) {
					return state + 100, nil
				},
				func(state int) int {
					return state - 100
				},
			).
			Then(
				// 尝试增加21，但返回了一个错误。因此，事务回滚到初始状态，需要 -100，并且状态保持为21。
				// (9 + 100) - 100
				func(state int) (int, error) {
					return state, assert.AnError
				},
				// 未执行
				func(state int) int {
					return state - 21
				},
			).
			Then(
				// 未执行
				func(state int) (int, error) {
					return state + 42, nil
				},
				// 未执行
				func(state int) int {
					return state - 42
				},
			)

		state, err := transaction.Process(9)
		is.Equal(9, state)
		is.Equal(assert.AnError, err)
	}

	// with error + update value
	{
		transaction := lo.NewTransaction[int]().
			Then(
				// 第一个Then块中，状态增加100，没有错误，所以状态变为 109。
				func(state int) (int, error) {
					return state + 100, nil
				},
				// 状态减少100（这是为了回滚准备，但因为没有错误，所以实际上状态仍然是 109）。
				func(state int) int {
					return state - 100
				},
			).
			Then(
				// 尝试增加21，但返回了一个错误。此时，事务尝试回滚，但因为下一个Then块中还有一个更新操作（增加42），所以状态被更新为42
				//（注意：这取决于lo.NewTransaction的具体实现，是否允许在错误发生后仍然执行后续的更新操作）。
				// (9 + 100) + 21 - (100)
				func(state int) (int, error) {
					return state + 21, assert.AnError
				},
				// 未执行
				func(state int) int {
					return state - 21
				},
			).
			Then(
				// 未执行
				func(state int) (int, error) {
					return state + 42, nil
				},
				// 未执行
				func(state int) int {
					return state - 42
				},
			)

		state, err := transaction.Process(9)
		is.Equal(30, state)
		is.Equal(assert.AnError, err)
	}
}
