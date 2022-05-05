package viper

import "github.com/spf13/pflag"

// FlagValueSet is an interface that users can implement
// to bind a set of flags to viper.
// 使用代理模式访问 pflag.FlagSet，只支持遍历FlagValue .
type FlagValueSet interface {
	// 遍历标签
	VisitAll(fn func(FlagValue))
}

// FlagValue is an interface that users can implement
// to bind different flags to viper.
// 使用代理模式访问 pflag.Flag，只支持部分接口.
type FlagValue interface {
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
}

// pflagValueSet is a wrapper around *pflag.ValueSet
// that implements FlagValueSet.
type PFlagValueSet struct {
	flags *pflag.FlagSet
}

// VisitAll iterates over all *pflag.Flag inside the *pflag.FlagSet.
// 遍历所有的 FlagValue .
func (p PFlagValueSet) VisitAll(fn func(flag FlagValue)) {
	p.flags.VisitAll(func(flag *pflag.Flag) {
		fn(PFlagValue{flag})
	})
}

// pflagValue is a wrapper aroung *pflag.flag
// that implements FlagValue .
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
