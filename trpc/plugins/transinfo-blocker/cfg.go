package transinfoblocker

import "fmt"

const (
	modeWhitelist = "whitelist"
	modeBlacklist = "blacklist"
	modeNone      = "none"
)

var cfg = &Config{}

// Config 插件配置
type Config struct {
	Default    *ListCfg            `yaml:"default"`      // Default 对所有client生效
	RPCNameCfg map[string]*ListCfg `yaml:"rpc_name_cfg"` // RPCNameCfg 对于特定的rpcname生效
}

// ListCfg 具体内容配置
type ListCfg struct {
	Mode string              `yaml:"mode"` // Mode "blocker, blacklist, none"
	Keys []string            `yaml:"keys"` // Keys blocker or blacklist key
	Set  map[string]struct{} `yaml:"-"`    // Set keys list to set
}

// parseBlockSet 格式化配置
func (l *ListCfg) parseBlockSet() error {
	if l == nil {
		return nil
	}
	if l.Mode != modeNone &&
		l.Mode != modeWhitelist &&
		l.Mode != modeBlacklist {
		return fmt.Errorf("unknown mod: %s", l.Mode)
	}

	// 将配置中的Keys转换为Set
	if l != nil {
		l.Set = make(map[string]struct{})
		for _, k := range l.Keys {
			l.Set[k] = struct{}{}
		}
	}

	return nil
}
