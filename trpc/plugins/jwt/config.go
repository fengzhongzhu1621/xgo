package jwt

// Config 插件配置
type Config struct {
	Secret       string   `yaml:"secret"`        // 签名使用的私钥
	Expired      int      `yaml:"expired"`       // 过期时间 seconds
	Issuer       string   `yaml:"issuer"`        // 发行人
	ExcludePaths []string `yaml:"exclude_paths"` // 跳过 jwt 鉴权的 paths, 如登陆接口
}
