package multiplexed

// GetOptions 获取连接的配置参数
type GetOptions struct {
	FP  IFrameParser // 帧解析器，用于解析网络帧
	VID uint32       // 虚拟连接ID，用于多路复用

	CACertFile    string // CA证书文件路径
	TLSCertFile   string // 客户端证书文件路径
	TLSKeyFile    string // 客户端私钥文件路径
	TLSServerName string // 客户端验证服务器的服务名称，如果未填写，默认为http主机名

	LocalAddr string // 建立连接时的本地地址

	network  string // 网络协议类型
	address  string // 目标地址
	isStream bool   // 是否为流式连接
	nodeKey  string // 节点键，用于标识连接
}

// NewGetOptions 创建并初始化GetOptions
// 返回:
//
//	GetOptions: 初始化后的配置选项
func NewGetOptions() GetOptions {
	return GetOptions{}
}

// WithFrameParser 设置单个Get操作的帧解析器
// 参数:
//
//	fp: 帧解析器接口
func (o *GetOptions) WithFrameParser(fp IFrameParser) {
	o.FP = fp
}

// WithDialTLS 设置客户端支持TLS的选项
// 参数:
//
//	certFile: 客户端证书文件路径
//	keyFile: 客户端私钥文件路径
//	caFile: CA证书文件路径
//	serverName: 服务器名称
func (o *GetOptions) WithDialTLS(certFile, keyFile, caFile, serverName string) {
	o.TLSCertFile = certFile
	o.TLSKeyFile = keyFile
	o.CACertFile = caFile
	o.TLSServerName = serverName
}

// WithVID 设置虚拟连接ID的选项
// 参数:
//
//	vid: 虚拟连接ID
func (o *GetOptions) WithVID(vid uint32) {
	o.VID = vid
}

// WithLocalAddr 设置建立连接时本地地址的选项
// 当有多个网卡时，默认随机选择
// 参数:
//
//	addr: 本地地址
func (o *GetOptions) WithLocalAddr(addr string) {
	o.LocalAddr = addr
}

// update 更新配置选项的网络和地址信息
// 参数:
//
//	network: 网络协议类型
//	address: 目标地址
//
// 返回:
//
//	error: 更新过程中发生的错误
func (o *GetOptions) update(network, address string) error {
	// 检查帧解析器是否设置
	if o.FP == nil {
		return ErrFrameParserNil
	}

	// 判断是否为流式连接
	isStream, err := isStream(network)
	if err != nil {
		return err
	}

	// 更新配置信息
	o.isStream = isStream
	o.address = address
	o.network = network
	o.nodeKey = makeNodeKey(o.network, o.address) // 生成节点键
	return nil
}
