package command

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fengzhongzhu1621/xgo/utils/proto"
)

type Cmd struct {
	baseCmd

	val interface{} // 命令的返回结果
}

var _ Cmder = (*Cmd)(nil)

// NewCmd 创建一个命令.
func NewCmd(ctx context.Context, args ...interface{}) *Cmd {
	return &Cmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *Cmd) String() string {
	return cmdString(cmd, cmd.val)
}

// 获得命令的返回结果.
func (cmd *Cmd) Val() interface{} {
	return cmd.val
}

func (cmd *Cmd) Result() (interface{}, error) {
	return cmd.val, cmd.err
}

func (cmd *Cmd) Text() (string, error) {
	if cmd.err != nil {
		return "", cmd.err
	}
	switch val := cmd.val.(type) {
	case string:
		return val, nil
	default:
		err := fmt.Errorf("redis: unexpected type=%T for String", val)
		return "", err
	}
}

func (cmd *Cmd) Int() (int, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	switch val := cmd.val.(type) {
	case int64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Int", val)
		return 0, err
	}
}

func (cmd *Cmd) Int64() (int64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	switch val := cmd.val.(type) {
	case int64:
		return val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Int64", val)
		return 0, err
	}
}

func (cmd *Cmd) Uint64() (uint64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	switch val := cmd.val.(type) {
	case int64:
		return uint64(val), nil
	case string:
		return strconv.ParseUint(val, 10, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Uint64", val)
		return 0, err
	}
}

func (cmd *Cmd) Float32() (float32, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	switch val := cmd.val.(type) {
	case int64:
		return float32(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, err
		}
		return float32(f), nil
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Float32", val)
		return 0, err
	}
}

func (cmd *Cmd) Float64() (float64, error) {
	if cmd.err != nil {
		return 0, cmd.err
	}
	switch val := cmd.val.(type) {
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Float64", val)
		return 0, err
	}
}

func (cmd *Cmd) Bool() (bool, error) {
	if cmd.err != nil {
		return false, cmd.err
	}
	switch val := cmd.val.(type) {
	case int64:
		return val != 0, nil
	case string:
		return strconv.ParseBool(val)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Bool", val)
		return false, err
	}
}

func (cmd *Cmd) Slice() ([]interface{}, error) {
	if cmd.err != nil {
		return nil, cmd.err
	}
	switch val := cmd.val.(type) {
	case []interface{}:
		return val, nil
	default:
		return nil, fmt.Errorf("redis: unexpected type=%T for Slice", val)
	}
}

// 读取并解析指令的返回结果.
func (cmd *Cmd) readReply(rd *proto.Reader) (err error) {
	cmd.val, err = rd.ReadReply(sliceParser)
	return err
}

// sliceParser implements proto.MultiBulkParse.
func sliceParser(rd *proto.Reader, n int64) (interface{}, error) {
	vals := make([]interface{}, n)
	for i := 0; i < len(vals); i++ {
		v, err := rd.ReadReply(sliceParser)
		if err != nil {
			if err == proto.Nil {
				vals[i] = nil
				continue
			}
			if err, ok := err.(proto.RedisError); ok {
				vals[i] = err
				continue
			}
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}
