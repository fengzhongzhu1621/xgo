package router

import "time"

var _ IRouterConfig = (*RouterConfig)(nil)

type IRouterConfig interface {
	SetDefaults()
	Validate() error
}

// RouterConfig holds the Router's configuration options.
type RouterConfig struct {
	// CloseTimeout determines how long router should work for handlers when closing.
	CloseTimeout time.Duration
}

// SetDefaults 设置路由配置的默认值
func (c *RouterConfig) SetDefaults() {
	if c.CloseTimeout == 0 {
		// 默认30秒超时
		c.CloseTimeout = time.Second * 30
	}
}

// Validate returns Router configuration error, if any.
func (c RouterConfig) Validate() error {
	return nil
}
