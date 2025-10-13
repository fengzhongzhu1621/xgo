package config

import (
	"github.com/fengzhongzhu1621/xgo/config/hooks"
	"github.com/fengzhongzhu1621/xgo/config/provider"
	"github.com/fengzhongzhu1621/xgo/crypto/unmarshaler"
)

// WithCodec returns an option which sets the codec's name.
func WithCodec(name string) LoadOption {
	return func(c *XGoConfig) {
		c.decoder = unmarshaler.GetCodec(name)
	}
}

// WithProvider returns an option which sets the provider's name.
func WithProvider(name string) LoadOption {
	return func(c *XGoConfig) {
		c.p = provider.GetProvider(name)
	}
}

// WithExpandEnv replaces ${var} in raw bytes with environment value of var.
// Note, method TrpcConfig.Bytes will return the replaced bytes.
func WithExpandEnv() LoadOption {
	return func(c *XGoConfig) {
		c.expandEnv = true
	}
}

// WithWatch returns an option to start watch model
func WithWatch() LoadOption {
	return func(c *XGoConfig) {
		c.watch = true
	}
}

// WithWatchHook returns an option to set log func for config change logger
func WithWatchHook(f func(msg hooks.WatchMessage)) LoadOption {
	return func(c *XGoConfig) {
		c.watchHook = f
	}
}

// options is config option.
type options struct{}

// Option is the option for config provider sdk.
type Option func(*options)
