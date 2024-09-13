package kafkax

type Kafka struct {
	Id       string `yaml:"id"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Topic    string `yaml:"topic"`
	Version  string `yaml:"version"`
	LogPath  string `yaml:"log_path"`
}
