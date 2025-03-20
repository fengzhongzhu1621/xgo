package coreservice

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo/ginx/rest"
)

type CoreServiceClientInterface interface {
}

type coreService struct {
	restCli rest.ClientInterface
}

// NewCoreServiceClient TODO
func NewCoreServiceClient(c *rest.Capability, version string) CoreServiceClientInterface {
	base := fmt.Sprintf("/api/%s", version)
	return &coreService{
		restCli: rest.NewRESTClient(c, base),
	}
}
