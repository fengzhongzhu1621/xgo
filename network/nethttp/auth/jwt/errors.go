package jwt

import "errors"

var (
	ErrUnauthorized = errors.New("jwt auth: token is unauthorized")

	ErrExpired    = errors.New("jwt auth: token is expired")
	ErrNBFInvalid = errors.New("jwt auth: token nbf validation failed")
	ErrIATInvalid = errors.New("jwt auth: token iat validation failed")

	ErrAPIGatewayJWTMissingApp             = errors.New("app not in jwt claims")
	ErrAPIGatewayJWTAppInfoParseFail       = errors.New("app info parse fail")
	ErrAPIGatewayJWTAppInfoNoAppCode       = errors.New("app_code not in app info")
	ErrAPIGatewayJWTAppCodeNotString       = errors.New("app_code not string")
	ErrAPIGatewayJWTAppInfoNoVerified      = errors.New("verified not in app info")
	ErrAPIGatewayJWTAppInfoVerifiedNotBool = errors.New("verified not bool")
	ErrAPIGatewayJWTAppNotVerified         = errors.New("app not verified")

	ErrAPIGatewayJWTAppInfoNoUsername = errors.New("username not in jwt claims")

	ErrJwtRootTokenNone = errors.New("jwt root token not found")
)
