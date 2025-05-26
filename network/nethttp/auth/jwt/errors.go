package jwt

import "errors"

var (
	ErrUnauthorized = errors.New("jwt auth: token is unauthorized")

	ErrExpired    = errors.New("jwt auth: token is expired")
	ErrNBFInvalid = errors.New("jwt auth: token nbf validation failed")
	ErrIATInvalid = errors.New("jwt auth: token iat validation failed")

	ErrAPIGatewayJWTAppNotFound       = errors.New("app not in jwt claims")
	ErrAPIGatewayJWTAppInfoNoAppCode  = errors.New("app_code not in app info")
	ErrAPIGatewayJWTAppCodeNotString  = errors.New("app_code validation failed")
	ErrAPIGatewayJWTAppInfoNoVerified = errors.New("verified not in app info")

	ErrAPIGatewayJWTAppInfoUserNotFound = errors.New("user not in jwt claims")

	ErrJwtRootTokenNone = errors.New("jwt root token not found")
)
