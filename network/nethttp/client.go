package nethttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/nethttp/ssl"
)

type HttpClient struct {
	caFile   string
	certFile string
	keyFile  string
	header   map[string]string
	httpCli  *http.Client
}

// NewHttpClient 构造函数
func NewHttpClient() *HttpClient {
	return &HttpClient{
		httpCli: &http.Client{},
		header:  make(map[string]string),
	}
}

func (client *HttpClient) GetClient() *http.Client {
	return client.httpCli
}

// NewTransPort 创建并返回一个 http.Transport 对象
func (client *HttpClient) NewTransPort() *http.Transport {
	return &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second, // TLS 握手超时时间设置为 5 秒
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second, // 连接超时时间为 5 秒
			KeepAlive: 30 * time.Second, // 保持连接活跃时间为 30 秒
		}).Dial,
		ResponseHeaderTimeout: 30 * time.Second, // 响应头超时时间设置为 30 秒
	}
}

// SetTlsNoVerity 设置客户端不验证服务器的 TLS 证书
func (client *HttpClient) SetTlsNoVerity() error {
	// 跳过TLS服务器证书验证
	tlsConf := ssl.ClientTslConfNoVerity()
	// 创建并设置一个 http.Transport 对象
	client.SetTlsVerityConfig(tlsConf)

	return nil
}

// SetTlsVerityServer 设置客户端需要验证服务器的 TLS 证书
func (client *HttpClient) SetTlsVerityServer(caFile string) error {
	client.caFile = caFile

	// 用于配置 TLS 客户端的 SSL/TLS 设置，以便验证服务器的证书
	tlsConf, err := ssl.ClientTslConfVerityServer(caFile)
	if err != nil {
		return err
	}
	// 创建并设置一个 http.Transport 对象
	client.SetTlsVerityConfig(tlsConf)

	return nil
}

// SetTlsVerity 需要验证客户端证书
func (client *HttpClient) SetTlsVerity(caFile, certFile, keyFile, passwd string) error {
	client.caFile = caFile
	client.certFile = certFile
	client.keyFile = keyFile

	// 用于创建一个 TLS 客户端的 SSL/TLS 配置，包括加载 CA 证书、客户端证书和私钥
	tlsConf, err := ssl.ClientTslConfVerity(caFile, certFile, keyFile, passwd)
	if err != nil {
		return err
	}
	// 创建并设置一个 http.Transport 对象
	client.SetTlsVerityConfig(tlsConf)

	return nil
}


// SetTlsVerityConfig 创建并设置一个 http.Transport 对象
func (client *HttpClient) SetTlsVerityConfig(tlsConf *tls.Config) {
	trans := client.NewTransPort()
	trans.TLSClientConfig = tlsConf
	client.httpCli.Transport = trans
}

func (client *HttpClient) SetTransPort(transport http.RoundTripper) {
	client.httpCli.Transport = transport
}

func (client *HttpClient) SetTimeOut(timeOut time.Duration) {
	client.httpCli.Timeout = timeOut
}

func (client *HttpClient) SetHeader(key, value string) {
	client.header[key] = value
}

func (client *HttpClient) SetBatchHeader(headerSet []*HeaderSet) {
	if headerSet == nil {
		return
	}
	for _, header := range headerSet {
		client.header[header.Key] = header.Value
	}
}
