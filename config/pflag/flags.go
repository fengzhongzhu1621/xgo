package flagutils

import (
	"github.com/spf13/pflag"
)

// IFlagValueSet is an interface that users can implement
// to bind a set of flags to viper.
// 使用代理模式访问 pflag.FlagSet，只支持遍历FlagValue .
type IFlagValueSet interface {
	// VisitAll 遍历标签
	VisitAll(fn func(IFlagValue))
}

// IFlagValue is an interface that users can implement
// to bind different flags to viper.
// 使用代理模式访问 pflag.Flag，只支持部分接口.
type IFlagValue interface {
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
}

var _ IFlagValueSet = (*PFlagValueSet)(nil)
var _ IFlagValue = (*PFlagValue)(nil)

// PFlagValueSet is a wrapper around *pflag.ValueSet
// that implements IFlagValueSet.
type PFlagValueSet struct {
	flags *pflag.FlagSet
}

// VisitAll iterates over all *pflag.Flag inside the *pflag.FlagSet.
// 遍历所有的 IFlagValue .
func (p PFlagValueSet) VisitAll(fn func(flag IFlagValue)) {
	p.flags.VisitAll(func(flag *pflag.Flag) {
		fn(PFlagValue{flag})
	})
}

// PFlagValue is a wrapper aroung *pflag.flag
// that implements IFlagValue .
type PFlagValue struct {
	flag *pflag.Flag
}

// HasChanged returns whether the flag has changes or not.
// If the user set the value (or if left to default) .
func (p PFlagValue) HasChanged() bool {
	return p.flag.Changed
}

// Name returns the name of the flag.
func (p PFlagValue) Name() string {
	return p.flag.Name
}

// ValueString returns the value of the flag as a string.
func (p PFlagValue) ValueString() string {
	return p.flag.Value.String()
}

// ValueType returns the type of the flag as a string.
func (p PFlagValue) ValueType() string {
	return p.flag.Value.Type()
}
