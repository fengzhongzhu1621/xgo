package xgo

import "fmt"

var (
	ErrUnauthorized = fmt.Errorf("%s", "unauthorized")
	ErrInvalidArg = fmt.Errorf("%s", "invalid argument")
	ErrInvalidAddress  = fmt.Errorf("%s", "invalid address")
	ErrUnknown = fmt.Errorf("%s", "unknown")
)


var errTable = map[int32]error {

}

/**
 * 将错误码转换为error
 */
func Code2Error(code int32) error {
	if err, ok := errTable[code]; ok {
		return err
	}
	return ErrUnknown
}
