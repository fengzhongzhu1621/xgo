package plugin

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v3"
)

var (
	// SetupTimeout is the timeout for initialization of each plugin.
	// Modify it if some plugins' initialization does take a long time.
	SetupTimeout = 3 * time.Second

	// MaxPluginSize is the max number of plugins.
	MaxPluginSize = 1000
)

// Config is the configuration of all plugins. plugin type => { plugin name => plugin config }
type Config map[string]map[string]yaml.Node

// SetupClosables loads plugins and returns a function to close them in reverse order.
func (c Config) SetupClosables() (close func() error, err error) {
	// load plugins one by one through the config file and put them into an ordered plugin queue.
	// 将配置转换为 pluginInfo 结构体，并放到队列中延迟处理
	plugins, status, err := c.loadPlugins()
	if err != nil {
		return nil, err
	}

	// remove and setup plugins one by one from the front of the ordered plugin queue.
	// 处理插件之间的依赖关系，执行插件Setup()，标记插件的加载状态
	pluginInfos, closes, err := c.setupPlugins(plugins, status)
	if err != nil {
		return nil, err
	}

	// notifies all plugins that plugin initialization is done.
	// 插件全部加载完成后的 OnFinish hook
	if err := c.onFinish(pluginInfos); err != nil {
		return nil, err
	}

	// 插件全部加载完成后的 Close hook
	return func() error {
		for i := len(closes) - 1; i >= 0; i-- {
			if err := closes[i](); err != nil {
				return err
			}
		}
		return nil
	}, nil
}

// loadPlugins 将配置转换为 pluginInfo 结构体，并放到队列中延迟处理
func (c Config) loadPlugins() (chan pluginInfo, map[string]bool, error) {
	var (
		plugins = make(chan pluginInfo, MaxPluginSize) // use channel as plugin queue
		// plugins' status. plugin key => {true: init done, false: init not done}.
		status = make(map[string]bool)
	)
	for typ, factories := range c {
		for name, cfg := range factories {
			factory := Get(typ, name)
			if factory == nil {
				return nil, nil, fmt.Errorf("plugin %s:%s no registered or imported, do not configure", typ, name)
			}
			p := pluginInfo{
				factory: factory,
				typ:     typ,
				name:    name,
				cfg:     cfg,
			}
			select {
			case plugins <- p:
			default:
				return nil, nil, fmt.Errorf("plugin number exceed max limit:%d", len(plugins))
			}
			status[p.key()] = false
		}
	}
	return plugins, status, nil
}

func (c Config) setupPlugins(plugins chan pluginInfo, status map[string]bool) ([]pluginInfo, []func() error, error) {
	var (
		result []pluginInfo
		closes []func() error
		num    = len(plugins)
	)
	for num > 0 {
		for i := 0; i < num; i++ {
			// 从队列中获取待处理的插件
			p := <-plugins

			// 处理插件之间的依赖关系
			// check if plugins that current plugin depends on have been initialized
			if deps, err := p.hasDependence(status); err != nil {
				return nil, nil, err
			} else if deps {
				// There are plugins that current plugin depends on haven't been initialized,
				// move current plugin to tail of the channel.
				plugins <- p
				continue
			}

			// 初始化插件
			if err := p.setup(); err != nil {
				return nil, nil, err
			}

			// 记录关闭函数
			if closer, ok := p.asCloser(); ok {
				closes = append(closes, closer.Close)
			}

			// 标记插件加载成功
			status[p.key()] = true
			// 记录返回成功的插件
			result = append(result, p)
		}

		if len(plugins) == num { // none plugin is setup, circular dependency exists.
			return nil, nil, fmt.Errorf("cycle depends, not plugin is setup")
		}
		num = len(plugins) // continue to process plugins that were moved to tail of the channel.
	}

	return result, closes, nil
}

func (c Config) onFinish(plugins []pluginInfo) error {
	for _, p := range plugins {
		if err := p.onFinish(); err != nil {
			return err
		}
	}
	return nil
}
