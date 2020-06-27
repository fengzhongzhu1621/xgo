package request

import (
	"net"
	"net/http"
	"time"
	"github.com/fengzhongzhu1621/xgo/services/name/client/address"
	opts "github.com/fengzhongzhu1621/xgo/services/name/client/options"
)

type HttpRequest struct {
	ServiceAddress    address.NameServiceAddress
	InitOptions       opts.InitOptions
	HttpClient        *http.Client
	PollingHttpClient *http.Client
}

func (request *HttpRequest) Init(initOpts opts.InitOptions) (err error) {
	// 解析配置中心服务器地址
	err = request.ServiceAddress.ParseAddress(initOpts.Address)
	if err != nil {
		return err
	}
	request.InitOptions = initOpts
	// 创建HTTP连接
	request.HttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   initOpts.RequestTimeout,		// 三次握手超时时间
				KeepAlive: initOpts.RequestPollingTimeout + 3*time.Second,	// 连接保持时间
			}).DialContext,
			TLSHandshakeTimeout: initOpts.RequestTimeout,	// 等待TLS握手完成的最长时间
		},
		Timeout: initOpts.RequestTimeout,
	}
	request.PollingHttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   initOpts.RequestPollingTimeout,
				KeepAlive: initOpts.RequestPollingTimeout + 3*time.Second,
			}).DialContext,
			TLSHandshakeTimeout: initOpts.RequestTimeout,
		},
		Timeout: initOpts.RequestPollingTimeout + 1*time.Second,
	}
	return nil
}
