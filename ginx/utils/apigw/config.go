package apigw

import (
	"github.com/fengzhongzhu1621/xgo/config/viper/viper_parser"
	"github.com/fengzhongzhu1621/xgo/ginx/utils/tls"
	"github.com/fengzhongzhu1621/xgo/network/ssl"
	log "github.com/sirupsen/logrus"
)

// ApiGWConfig api gateway config
type ApiGWConfig struct {
	AppCode   string
	AppSecret string
	Username  string
	TLSConfig *ssl.TLSClientConfig
}

// ParseApiGWConfig parse api gateway config
func ParseApiGWConfig(path string) (*ApiGWConfig, error) {
	appCode, err := viper_parser.String(path + ".appCode")
	if err != nil {
		log.Errorf("get api gateway appCode config error, err: %v", err)
		return nil, err
	}

	appSecret, err := viper_parser.String(path + ".appSecret")
	if err != nil {
		log.Errorf("get api gateway appSecret config error, err: %v", err)
		return nil, err
	}

	username, err := viper_parser.String(path + ".username")
	if err != nil {
		log.Errorf("get api gateway username config error, err: %v", err)
		return nil, err
	}

	tlsConfig, err := tls.NewTLSClientConfigFromConfig(path + ".tls")
	if err != nil {
		log.Errorf("get api gateway tls config error, err: %v", err)
		return nil, err
	}

	apiGWConfig := &ApiGWConfig{
		AppCode:   appCode,
		AppSecret: appSecret,
		Username:  username,
		TLSConfig: &tlsConfig,
	}
	return apiGWConfig, nil
}
