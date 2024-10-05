package gomonkey

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestApplyFuncSeq(t *testing.T) {
	// 第一次调用 add 函数时返回 3，第二次调用时返回 5
	outputs := []gomonkey.OutputCell{
		{Values: gomonkey.Params{3}, Times: 1},
		{Values: gomonkey.Params{5}, Times: 1},
	}

	patches := gomonkey.ApplyFuncSeq(add, outputs)
	defer patches.Reset()

	assert.Equal(t, 3, add(1, 2))
	assert.Equal(t, 5, add(1, 2))
}
