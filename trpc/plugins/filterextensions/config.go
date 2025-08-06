package filterextensions

type cfg struct {
	Client []cfgService `yaml:"client"`
	Server []cfgService `yaml:"server"`
}

type cfgService struct {
	Name    string      `yaml:"name"`
	Methods []cfgMethod `yaml:"methods"`
}

type cfgMethod struct {
	Name    string   `yaml:"name"`
	Filters []string `yaml:"filters"`
}
