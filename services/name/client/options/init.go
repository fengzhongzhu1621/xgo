package name

import (
	"time"
)

type InitOptions struct {
	Address string	// 配置中心地址
	AppId string
	Groups []string
	userId string
	UserKey string
	HmacAlgorithm string					// 签名算法
	Heartbeat time.Duration					// 心跳间隔
	RequestTimeout time.Duration         	// 请求超时时间
	RequestPollingTimeout time.Duration  	// 长轮询等待时间
	ProtoType string                        // 协议类型
}

func (options *InitOptions) Init() {
	if options.RequestTimeout == time.Duration(0) {
		// 默认0.5s
		options.RequestTimeout = 500 * time.Millisecond
	}
	if options.RequestPollingTimeout == time.Duration(0) {
		// 默认1分钟
		options.RequestPollingTimeout = 60 * time.Second
	}
	if options.ProtoType == "" {
		options.ProtoType = "http"
	}
	if options.Heartbeat == time.Duration(0) {
		// 默认1分钟
		options.Heartbeat = 60 * time.Second
	}
	if options.HmacAlgorithm == "" {
		options.HmacAlgorithm = "sha256"
	}
}

func (options *InitOptions) SetAddress(address string) {
	options.Address = address
}

func (options *InitOptions) AddGroup(group ...string) {
	options.Groups = append(options.Groups, group...)
}
