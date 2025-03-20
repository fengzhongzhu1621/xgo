package configcenter

import (
	"github.com/fengzhongzhu1621/xgo/db/zookeeper/confregdiscover"
)

// ConfigCenter TODO
type ConfigCenter struct {
	Type               string // type of configuration center
	ConfigCenterDetail confregdiscover.ConfRegDiscvIf
}

var (
	configCenterGroup []*ConfigCenter
	configCenterType  = DefaultConfigCenter // the default configuration center is zookeeper.
)

// SetConfigCenterType use this function to change the type of configuration center.
func SetConfigCenterType(serverType string) {
	configCenterType = serverType
}

// AddConfigCenter add the configuration center you want to replace.
func AddConfigCenter(configCenter *ConfigCenter) {
	configCenterGroup = append(configCenterGroup, configCenter)
}

// CurrentConfigCenter use this method to return to the configuration center you want to use.
func CurrentConfigCenter() confregdiscover.ConfRegDiscvIf {
	var defaultConfigCenter *ConfigCenter
	for _, center := range configCenterGroup {
		if center.Type == configCenterType {
			return center.ConfigCenterDetail
		}
		if DefaultConfigCenter == center.Type {
			defaultConfigCenter = center
		}
	}
	if nil != defaultConfigCenter {
		return defaultConfigCenter.ConfigCenterDetail
	}
	return nil
}
