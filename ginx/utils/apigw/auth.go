package apigw

import (
	"encoding/json"
	"net/http"

	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	log "github.com/sirupsen/logrus"
)

type AppAuthConfig struct {
	AppCode   string `json:"bk_app_code,omitempty"`
	AppSecret string `json:"bk_app_secret,omitempty"`
}

// AuthConfig defines the api gateway authorization config
type AuthConfig struct {
	AppAuthConfig `json:",inline"`
	BkToken       string `json:"bk_token,omitempty"`
	BkTicket      string `json:"bk_ticket,omitempty"`
	UserName      string `json:"bk_username,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
}

// SetAuthHeader set api gateway authorization header
func SetAuthHeader(appConf AppAuthConfig, header http.Header) http.Header {
	conf := AuthConfig{
		AppAuthConfig: appConf,
		BkToken:       nethttp.GetUserToken(header),
		BkTicket:      nethttp.GetUserTicket(header),
		UserName:      nethttp.GetUser(header),
	}

	authInfo, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("marshal api auth config %+v failed, err: %v, rid: %s", conf, err, nethttp.GetRid(header))
		return header
	}

	return nethttp.SetBkAuth(header, string(authInfo))
}

// GenDefaultAuthHeader generate api gateway default authorization header
func GenDefaultAuthHeader(conf *ApiGWConfig) (string, error) {
	authConf := AuthConfig{
		AppAuthConfig: AppAuthConfig{
			AppCode:   conf.AppCode,
			AppSecret: conf.AppSecret,
		},
		UserName: conf.Username,
	}

	authInfo, err := json.Marshal(authConf)
	if err != nil {
		log.Errorf("marshal default api auth config %+v failed, err: %v", conf, err)
		return "", err
	}

	return string(authInfo), nil
}

func SetApiGWAuthHeader(conf *ApiGWConfig, header http.Header) http.Header {
	appConf := AppAuthConfig{
		AppCode:   conf.AppCode,
		AppSecret: conf.AppSecret,
	}
	return SetAuthHeader(appConf, header)
}
