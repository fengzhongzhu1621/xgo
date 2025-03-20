package configcenter

// ProcessConfig TODO
type ProcessConfig struct {
	ConfigData []byte
}

// ParseConfigWithData TODO
func ParseConfigWithData(data []byte) *ProcessConfig {

	return &ProcessConfig{ConfigData: data}
}
