package router

import "time"

// RouterConfig holds the Router's configuration options.
type RouterConfig struct {
	// CloseTimeout determines how long router should work for handlers when closing.
	CloseTimeout time.Duration
}

// setDefaults 设置路由配置的默认值
func (c *RouterConfig) setDefaults() {
	if c.CloseTimeout == 0 {
		// 默认30秒超时
		c.CloseTimeout = time.Second * 30
	}
}

// Validate returns Router configuration error, if any.
func (c RouterConfig) Validate() error {
	return nil
}
