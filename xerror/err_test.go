package xerror

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:noinline
func parent() error {
	if err := child(); err != nil {
		return err
	}
	return nil
}

//go:noinline
func child() error {
	if err := grandson(); err != nil {
		return err
	}
	return nil
}

//go:noinline
func grandson() error {
	return Newf(111, "%s", "inner fail")
}

func TestErrs(t *testing.T) {
	var err *Error
	str := err.Error()
	assert.Contains(t, str, "success")

	e := New(111, "inner fail")
	assert.NotNil(t, e)

	assert.EqualValues(t, 111, Code(e))
	assert.Equal(t, "inner fail", Msg(e))

	err, ok := e.(*Error)
	assert.Equal(t, true, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeBusiness, err.Type)

	str = err.Error()
	assert.Contains(t, str, "business")

	e = NewFrameError(111, "inner fail")
	assert.NotNil(t, e)

	assert.EqualValues(t, 111, Code(e))
	assert.Equal(t, "inner fail", Msg(e))

	err, ok = e.(*Error)
	assert.Equal(t, true, ok)
	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeFramework, err.Type)

	str = err.Error()
	assert.Contains(t, str, "framework")

	assert.EqualValues(t, 0, Code(nil))
	assert.Equal(t, "success", Msg(nil))

	assert.EqualValues(t, 0, Code((*Error)(nil)))
	assert.Equal(t, "success", Msg((*Error)(nil)))

	e = errors.New("unknown error")
	assert.Equal(t, RetUnknown, Code(e))
	assert.Equal(t, "unknown error", Msg(e))

	err.Type = ErrorTypeCalleeFramework
	assert.Contains(t, err.Error(), "type:callee framework")
}

func TestErrsFormat(t *testing.T) {
	err := New(10000, "test error")

	s := fmt.Sprintf("%s", err)
	assert.Equal(t, "type:business, code:10000, msg:test error", s)

	s = fmt.Sprintf("%q", err)
	assert.Equal(t, `"type:business, code:10000, msg:test error"`, s)

	s = fmt.Sprintf("%v", err)
	assert.Equal(t, "type:business, code:10000, msg:test error", s)

	s = fmt.Sprintf("%d", err)
	assert.Equal(t, "%!d(errs.Error=type:business, code:10000, msg:test error)", s)
}

func TestNewFrameError(t *testing.T) {
	ok := true
	SetTraceable(ok)
	e := NewFrameError(111, "inner fail")
	assert.NotNil(t, e)
}

func TestWrapFrameError(t *testing.T) {
	ok := true
	SetTraceable(ok)
	e := WrapFrameError(New(123, "inner fail"), 456, "wrap frame error")
	assert.NotNil(t, e)
	e = WrapFrameError(nil, 456, "wrap frame error")
	assert.Nil(t, e)
}

func TestTraceError(t *testing.T) {

	SetTraceable(true)

	err := parent()
	assert.NotNil(t, err)

	s := fmt.Sprintf("%+v", err)
	br := bufio.NewReader(strings.NewReader(s))

	line, isPrefix, err := br.ReadLine()
	assert.Equal(t, "type:business, code:111, msg:inner fail", string(line))
	assert.Equal(t, isPrefix, false)
	assert.Nil(t, err)

	line, isPrefix, err = br.ReadLine()
	assert.Equal(t, "github.com/fengzhongzhu1621/xgo/xerror.grandson", string(line))
	assert.Equal(t, isPrefix, false)
	assert.Nil(t, err)

	_, _, _ = br.ReadLine()
	line, isPrefix, err = br.ReadLine()
	assert.Equal(t, "github.com/fengzhongzhu1621/xgo/xerror.child", string(line))
	assert.Equal(t, isPrefix, false)
	assert.Nil(t, err)

	_, _, _ = br.ReadLine()
	line, isPrefix, err = br.ReadLine()
	assert.Equal(t, "github.com/fengzhongzhu1621/xgo/xerror.parent", string(line))
	assert.Equal(t, isPrefix, false)
	assert.Nil(t, err)

	// print stack detail info
	err = parent()
	fmt.Printf("%+v", err)
}
