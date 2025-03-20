package apigw

import (
	"github.com/fengzhongzhu1621/xgo/config/viper/viper_parser"
	"github.com/fengzhongzhu1621/xgo/ginx/rest"
	log "github.com/sirupsen/logrus"
)

// ApiGWOptions is the api gateway client options
type ApiGWOptions struct {
	Config     *ApiGWConfig
	Auth       string
	Capability rest.Capability
}

// ApiGWSrv api gateway service
type ApiGWSrv struct {
	Client rest.ClientInterface
	Config *ApiGWConfig
	Auth   string
}

// NewApiGW new api gateway service
func NewApiGW(options *ApiGWOptions, addressPath string) (*ApiGWSrv, error) {
	address, err := viper_parser.StringSlice(addressPath)
	if err != nil {
		log.Errorf("parse %s api gateway address failed, err: %v", addressPath, err)
		return nil, err
	}

	capability := options.Capability
	capability.Discover = &ApiGWDiscovery{
		Servers: address,
	}

	apigw := &ApiGWSrv{
		Client: rest.NewRESTClient(&capability, "/"),
		Config: options.Config,
		Auth:   options.Auth,
	}

	return apigw, nil
}
