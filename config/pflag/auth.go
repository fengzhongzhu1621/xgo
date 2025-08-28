package pflag

import (
	"strconv"
	"sync"

	"github.com/spf13/pflag"
)

var once = sync.Once{}

var (
	EnableAuthFlag *authValue
	enableAuth     = true // 默认值
)

type authValue struct{}

var _ pflag.Value = (*authValue)(nil)

func (a *authValue) String() string {
	return strconv.FormatBool(enableAuth)
}

func (a *authValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	setEnableAuth(v)
	return nil
}

// Type TODO
func (a *authValue) Type() string {
	return "bool"
}

func setEnableAuth(enable bool) {
	once.Do(func() {
		enableAuth = enable
	})
}
