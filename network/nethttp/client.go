package nethttp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/nethttp/ssl"
)
type HttpRespone struct {
	Reply      []byte // 响应内容
	StatusCode int
	Status     string
	Header     http.Header
}

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

// RequestEx do http request, old version
func (client *HttpClient) RequestEx(url, method string, header http.Header, data []byte) (*HttpRespone, error) {
	var req *http.Request
	var errReq error
	httpRsp := &HttpRespone{
		Reply:      nil,
		StatusCode: http.StatusInternalServerError,
		Status:     "Internal Server Error",
	}

	// 用于创建一个新的 HTTP 请求。该函数接受三个参数：请求方法（如 "GET"、"POST" 等），请求的 URL 和可选的请求体
	// （对于某些请求方法，如 "GET"，请求体通常为空）。
	if data != nil {
		req, errReq = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, errReq = http.NewRequest(method, url, nil)
	}

	if errReq != nil {
		return httpRsp, errReq
	}

	req.Close = true

	// 设置 header
	if header != nil {
		req.Header = header
	}
	for key, value := range client.header {
		if req.Header.Get(key) != "" {
			continue
		}
		req.Header.Set(key, value)
	}

	// 发送请求
	rsp, err := client.httpCli.Do(req)
	if err != nil {
		return httpRsp, err
	}

	defer rsp.Body.Close()

	httpRsp.Status = rsp.Status
	httpRsp.StatusCode = rsp.StatusCode
	httpRsp.Header = rsp.Header

	// 读取响应内容
	rpy, err := io.ReadAll(rsp.Body)
	if err != nil {
		return httpRsp, err
	}
	httpRsp.Reply = rpy

	return httpRsp, nil
}

func (client *HttpClient) RequestStream(url, method string, header http.Header, data []byte) (io.ReadCloser, error) {
	var req *http.Request
	var errReq error

	// 用于创建一个新的 HTTP 请求。该函数接受三个参数：请求方法（如 "GET"、"POST" 等），请求的 URL 和可选的请求体
	// （对于某些请求方法，如 "GET"，请求体通常为空）。
	if data != nil {
		req, errReq = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, errReq = http.NewRequest(method, url, nil)
	}

	if errReq != nil {
		return nil, errReq
	}

	req.Close = true

	// 设置 header
	if header != nil {
		req.Header = header
	}
	for key, value := range client.header {
		if req.Header.Get(key) != "" {
			continue
		}
		req.Header.Set(key, value)
	}

	// 发送请求
	rsp, err := client.httpCli.Do(req)
	if err != nil {
		return nil, err
	}

	switch {
	case (rsp.StatusCode >= 200) && (rsp.StatusCode < 300):
		return rsp.Body, nil
	default:
		defer rsp.Body.Close()
		return nil, fmt.Errorf("get stream failed, resp code %d", rsp.StatusCode)
	}
}

func (client *HttpClient) GET(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "GET", header, data)
}

func (client *HttpClient) POST(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "POST", header, data)
}

func (client *HttpClient) DELETE(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "DELETE", header, data)
}

func (client *HttpClient) PUT(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "PUT", header, data)
}

func (client *HttpClient) PATCH(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "PATCH", header, data)
}

func (client *HttpClient) Get(url string, header http.Header, data []byte) (*HttpRespone, error) {
	return client.RequestEx(url, "GET", header, data)
}

func (client *HttpClient) Post(url string, header http.Header, data []byte) (*HttpRespone, error) {
	return client.RequestEx(url, "POST", header, data)
}

func (client *HttpClient) Delete(url string, header http.Header, data []byte) (*HttpRespone, error) {
	return client.RequestEx(url, "DELETE", header, data)
}

func (client *HttpClient) Put(url string, header http.Header, data []byte) (*HttpRespone, error) {
	return client.RequestEx(url, "PUT", header, data)
}

func (client *HttpClient) Patch(url string, header http.Header, data []byte) (*HttpRespone, error) {
	return client.RequestEx(url, "PATCH", header, data)
}

func (client *HttpClient) Request(url, method string, header http.Header, data []byte) ([]byte, error) {
	rsp, err := client.RequestEx(url, method, header, data)
	return rsp.Reply, err
}