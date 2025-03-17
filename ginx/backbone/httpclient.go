package backbone

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/fengzhongzhu1621/xgo/monitor/opentelemetry"
	"github.com/fengzhongzhu1621/xgo/network/ssl"
)

// NewClient create a new http client
func NewClient(c *TLSClientConfig, conf ...ExtraClientConfig) (*http.Client, error) {
	tlsConf := new(tls.Config)
	if c != nil && len(c.CAFile) != 0 && len(c.CertFile) != 0 && len(c.KeyFile) != 0 {
		var err error
		tlsConf, err = ssl.ClientTslConfVerity(c.CAFile, c.CertFile, c.KeyFile, c.Password)
		if err != nil {
			return nil, err
		}
	}

	if c != nil {
		tlsConf.InsecureSkipVerify = c.InsecureSkipVerify
	}

	// set api request timeout to 25s, so that we can stop the long request like searching all hosts
	responseHeaderTimeout := 25 * time.Second
	if len(conf) > 0 {
		if timeout := conf[0].ResponseHeaderTimeout; timeout != 0 {
			responseHeaderTimeout = timeout
		}
	}
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     tlsConf,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost:   100,
		ResponseHeaderTimeout: responseHeaderTimeout,
	}

	client := &http.Client{
		Transport: transport,
	}

	opentelemetry.WrapperTraceClient(client)

	return client, nil
}

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
