// Package plugin implements a general plugin factory system which provides plugin registration and loading.
// It is mainly used when certain plugins must be loaded by configuration.
// This system is not supposed to register plugins that do not rely on configuration like codec. Instead, plugins
// that do not rely on configuration should be registered by calling methods in certain packages.
package plugin

var plugins = make(map[string]map[string]IFactory) // plugin type => { plugin name => plugin factory }

// Factory is the interface for plugin factory abstraction.
// Custom Plugins need to implement this interface to be registered as a plugin with certain type.
type IFactory interface {
	// Type returns type of the plugin, i.e. selector, log, config, tracing.
	Type() string
	// Setup loads plugin by configuration.
	// The data structure of the configuration of the plugin needs to be defined in advance。
	Setup(name string, dec IDecoder) error
}

// Register registers a plugin factory.
// Name of the plugin should be specified.
// It is supported to register instances which are the same implementation of plugin Factory
// but use different configuration.
func Register(name string, f IFactory) {
	factories, ok := plugins[f.Type()]
	if !ok {
		factories = make(map[string]IFactory)
		plugins[f.Type()] = factories
	}
	factories[name] = f
}

// Get returns a plugin Factory by its type and name.
func Get(typ string, name string) IFactory {
	return plugins[typ][name]
}
