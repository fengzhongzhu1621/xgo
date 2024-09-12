package auth

import (
	"encoding/base64"

	"github.com/fengzhongzhu1621/xgo/cast"
)

func BasicAuthAuthorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(cast.StringToBytes(base))
}
