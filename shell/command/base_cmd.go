package command

import (
	"context"
	"time"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/proto"
	"github.com/fengzhongzhu1621/xgo/str/bytesutils"
	"github.com/fengzhongzhu1621/xgo/str/stringutils"
)

type Cmder interface {
	Name() string // 获得命令的名称
	FullName() string
	Args() []interface{}
	String() string
	stringArg(int) string // 获得指定位置的参数
	firstKeyPos() int8
	setFirstKeyPos(int8)

	readTimeout() *time.Duration
	readReply(rd *proto.Reader) error

	SetErr(error)
	Err() error
}

// 将命令串转换为字符串表示，仅用于显示.
func cmdString(cmd Cmder, val interface{}) string {
	b := make([]byte, 0, 64)

	// 变量命令的所有参数
	for i, arg := range cmd.Args() {
		if i > 0 {
			b = append(b, ' ')
		}
		b = bytesutils.AppendArg(b, arg)
	}

	if err := cmd.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	} else if val != nil {
		b = append(b, ": "...)
		b = bytesutils.AppendArg(b, val)
	}

	return cast.String(b)
}

// 私有接口，定义基类.
type baseCmd struct {
	ctx    context.Context
	args   []interface{} // 命令参数列表
	err    error
	keyPos int8

	_readTimeout *time.Duration
}

// 获得命令的名称，必须是小写，为第一个参数的值.
func (cmd *baseCmd) Name() string {
	if len(cmd.args) == 0 {
		return ""
	}
	// 获得指定位置的参数，Cmd name must be lower cased.
	return stringutils.ToLower(cmd.stringArg(0))
}

// 获得命令的全名.
func (cmd *baseCmd) FullName() string {
	switch name := cmd.Name(); name {
	case "cluster", "command":
		if len(cmd.args) == 1 {
			return name
		}
		if s2, ok := cmd.args[1].(string); ok {
			return name + " " + s2
		}
		return name
	default:
		return name
	}
}

// 获得参数列表.
func (cmd *baseCmd) Args() []interface{} {
	return cmd.args
}

// 获得指定位置的参数值.
func (cmd *baseCmd) stringArg(pos int) string {
	if pos < 0 || pos >= len(cmd.args) {
		return ""
	}
	// 参数值必须是字符串
	s, _ := cmd.args[pos].(string)
	return s
}

func (cmd *baseCmd) firstKeyPos() int8 {
	return cmd.keyPos
}

func (cmd *baseCmd) setFirstKeyPos(keyPos int8) {
	cmd.keyPos = keyPos
}

func (cmd *baseCmd) SetErr(e error) {
	cmd.err = e
}

func (cmd *baseCmd) Err() error {
	return cmd.err
}

func (cmd *baseCmd) readTimeout() *time.Duration {
	return cmd._readTimeout
}

func (cmd *baseCmd) setReadTimeout(d time.Duration) {
	cmd._readTimeout = &d
}
