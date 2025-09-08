package plugin

import (
	"errors"
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// pluginInfo is the information of a plugin.
type pluginInfo struct {
	factory IFactory  // 插件工厂，用于注册
	typ     string    // 插件类型
	name    string    // 插件名称
	cfg     yaml.Node // 插件配置
}

// hasDependence decides if any other plugins that this plugin depends on haven't been initialized.
// The input param is the initial status of all plugins.
// The output bool param being true means there are plugins that this plugin depends on haven't been initialized,
// while being false means this plugin doesn't depend on any other plugin or all the plugins that his plugin depends
// on have already been initialized.
func (p *pluginInfo) hasDependence(status map[string]bool) (bool, error) {
	deps, ok := p.factory.(IDepender)
	if ok {
		hasDeps, err := p.checkDependence(status, deps.DependsOn(), false)
		if err != nil {
			return false, err
		}
		if hasDeps { // 个别插件会同时强依赖和弱依赖多个不同插件，当所有强依赖满足后需要再判断弱依赖关系
			return true, nil
		}
	}
	fd, ok := p.factory.(IFlexDepender)
	if ok {
		return p.checkDependence(status, fd.FlexDependsOn(), true)
	}
	// This plugin doesn't depend on any other plugin.
	return false, nil
}

// checkDependence 检查插件的依赖关系
// 参数:
//   - status: 所有插件的初始化状态映射
//   - dependences: 当前插件依赖的其他插件名称列表
//   - flexible: 是否为弱依赖关系，true表示弱依赖，false表示强依赖
// 返回值:
//   - bool: true表示存在未初始化的依赖插件，false表示所有依赖插件都已初始化或不存在依赖
//   - error: 依赖检查过程中的错误，如自我依赖或强依赖的插件不存在
func (p *pluginInfo) checkDependence(status map[string]bool, dependences []string, flexible bool) (bool, error) {
	for _, name := range dependences {
		if name == p.key() {
			return false, errors.New("plugin not allowed to depend on itself")
		}
		setup, ok := status[name]
		if !ok {
			if flexible {
				continue
			}
			return false, fmt.Errorf("depends plugin %s not exists", name)
		}
		if !setup {
			return true, nil
		}
	}
	return false, nil
}

// setup initializes a single plugin.
func (p *pluginInfo) setup() error {
	var (
		ch  = make(chan struct{})
		err error
	)
	go func() {
		// 执行插件初始化，传入节点配置(yaml.Node)
		err = p.factory.Setup(p.name, &YamlNodeDecoder{Node: &p.cfg})
		close(ch)
	}()

	select {
	case <-ch:
	case <-time.After(SetupTimeout):
		return fmt.Errorf("setup plugin %s timeout", p.key())
	}
	if err != nil {
		return fmt.Errorf("setup plugin %s error: %v", p.key(), err)
	}
	return nil
}

// key returns the unique index of plugin in the format of 'type-name'.
func (p *pluginInfo) key() string {
	return p.typ + "-" + p.name
}

// onFinish notifies the plugin that all plugins' loading has been done by tRPC-Go.
func (p *pluginInfo) onFinish() error {
	f, ok := p.factory.(IFinishNotifier)
	if !ok {
		// IFinishNotifier not being implemented means notification of
		// completion of all plugins' loading is not needed.
		return nil
	}
	return f.OnFinish(p.name)
}

func (p *pluginInfo) asCloser() (ICloser, bool) {
	closer, ok := p.factory.(ICloser)
	return closer, ok
}
