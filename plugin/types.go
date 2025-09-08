package plugin

// Decoder is the interface used to decode plugin configuration.
type IDecoder interface {
	Decode(cfg interface{}) error // the input param is the custom configuration of the plugin
}

// IDepender is the interface for "Strong Dependence".
// If plugin a "Strongly" depends on plugin b, b must exist and
// a will be initialized after b's initialization.
// IDepender 是 "强依赖" 的接口。
// 如果插件 a "强烈" 依赖插件 b，那么 b 必须存在，强依赖要求被依赖的插件必须存在，不存在框架会 panic。
// a 将在 b 初始化之后进行初始化。
type IDepender interface {
	// DependsOn returns a list of plugins that are relied upon.
	// The list elements are in the format of "type-name" like [ "selector-polaris" ].
	// DependsOn 返回依赖的插件列表。
	// 列表元素的格式为 "类型-名称"，例如 [ "selector-polaris" ]。
	DependsOn() []string
}

// IFlexDepender is the interface for "Weak Dependence".
// If plugin a "Weakly" depends on plugin b and b does exist,
// a will be initialized after b's initialization.
// IFlexDepender 是 "弱依赖" 的接口。
// 如果插件 a "弱" 依赖插件 b，并且 b 确实存在，弱依赖则不会 panic。
// 那么 a 将在 b 初始化之后进行初始化。
type IFlexDepender interface {
	FlexDependsOn() []string
}

// IFinishNotifier is the interface used to notify that all plugins' loading has been done by tRPC-Go.
// Some plugins need to implement this interface to be notified when all other plugins' loading has been done.
type IFinishNotifier interface {
	// 插件全部加载完成后的 hook
	OnFinish(name string) error
}

// ICloser is the interface used to provide a close callback of a plugin.
type ICloser interface {
	// 插件全部加载完毕后，每个插件需要执行清理操作
	Close() error
}
