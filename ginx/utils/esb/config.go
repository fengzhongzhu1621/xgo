package esb

import (
	"net/http"

	"github.com/fengzhongzhu1621/xgo/ginx/utils/apigw"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
)

// EsbConfig TODO
type EsbConfig struct {
	Addrs     string
	AppCode   string
	AppSecret string
}

// EsbCommParams is esb common parameters
type EsbCommParams struct {
	SupplierID string `json:"bk_supplier_id"`
}

// SetEsbAuthHeader set esb authorization header
func SetEsbAuthHeader(esbConfig EsbConfig, header http.Header) http.Header {
	appConf := apigw.AppAuthConfig{
		AppCode:   esbConfig.AppCode,
		AppSecret: esbConfig.AppSecret,
	}
	return apigw.SetAuthHeader(appConf, header)
}

// GetEsbRequestParams get esb request parameters
func GetEsbRequestParams(esbConfig EsbConfig, header http.Header) *EsbCommParams {
	return &EsbCommParams{
		SupplierID: nethttp.GetSupplierAccount(header),
	}
}
