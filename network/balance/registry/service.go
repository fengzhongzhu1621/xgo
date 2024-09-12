package registry

// ServiceNode 服务器节点
type ServiceNode struct {
	ServerIp   string `json:"server_ip"  validate:"required"`        // 服务 IP
	ServerPort int    `json:"server_port" validate:"required,gt=0" ` // 服务 端口
	Weight     int    `json:"weight" validate:"required,gt=0"`       // 权重 最大值为 100
	Proto      string `json:"proto"`                                 // http/https/grpc
}

// Validate 验证服务器节点字段
func (n ServiceNode) Validate() error {
	return g_validator.Struct(&n)
}
