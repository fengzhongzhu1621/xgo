package cacheimpls

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cache "github.com/fengzhongzhu1621/xgo/cache/common"
	"github.com/fengzhongzhu1621/xgo/crypto/randutils"
)

var (
	ErrAPIGatewayJWTCacheNotFound     = errors.New("not found")
	ErrAPIGatewayJWTClientIDNotString = errors.New("clientID not string")
	ErrAPIGatewayJWTUsernameNotString = errors.New("userName not string")
)

// APIGatewayJWTClientIDCacheKey cache key for JWTToken
type APIGatewayJWTClientIDCacheKey struct {
	JWTToken string
}

// Key ...
func (k APIGatewayJWTClientIDCacheKey) Key() string {
	return randutils.MD5Hash(k.JWTToken)
}

var _ cache.Key = (*APIGatewayJWTClientIDCacheKey)(nil)

func retrieveAPIGatewayJWTClientID(ctx context.Context, key cache.Key) (interface{}, error) {
	// NOTE: this func not work
	return "", nil
}

func getJwtTokenVal(jwtToken string) (value interface{}, ok bool) {
	key := APIGatewayJWTClientIDCacheKey{
		JWTToken: jwtToken,
	}
	ctx := context.TODO()
	// 从缓存中获取用户信息
	return LocalAPIGatewayJWTClientIDCache.DirectGet(ctx, key)
}

// GetJWTTokenClientID will retrieve the clientID of a jwtTOken
func GetJWTTokenClientID(jwtToken string) (clientID string, err error) {
	// 从缓存中获取用户信息
	value, ok := getJwtTokenVal(jwtToken)
	if !ok {
		err = ErrAPIGatewayJWTCacheNotFound
		return
	}

	// 获得用户 id
	clientID, ok = value.(string)
	if !ok {
		err = ErrAPIGatewayJWTClientIDNotString
		return
	}
	return clientID, nil
}

// SetJWTTokenClientID will set the jwtToken-clientID int cache
func SetJWTTokenClientID(jwtToken string, clientID string) {
	key := APIGatewayJWTClientIDCacheKey{
		JWTToken: jwtToken,
	}
	ctx := context.TODO()
	LocalAPIGatewayJWTClientIDCache.Set(ctx, key, clientID)
}

// GetJWTTokenClientIDAndUsername
// GetJWTTokenUsername will retrieve the clientID of a jwtTOken
func GetJWTTokenClientIDAndUsername(jwtToken string) (clientID, userName string, err error) {
	value, ok := getJwtTokenVal(jwtToken)
	if !ok {
		err = ErrAPIGatewayJWTCacheNotFound
		return
	}

	val, ok := value.(string)
	if !ok {
		err = ErrAPIGatewayJWTUsernameNotString
		return
	}

	array := strings.Split(val, "#")
	if len(array) < 2 {
		err = ErrAPIGatewayJWTCacheNotFound
		return
	}
	clientID = array[0]
	userName = array[1]
	return clientID, userName, nil
}

// SetJWTTokenClientIDAndUsername will set the jwtToken-clientID and username int cache
func SetJWTTokenClientIDAndUsername(jwtToken string, clientID, userName string) {
	key := APIGatewayJWTClientIDCacheKey{
		JWTToken: jwtToken,
	}
	val := fmt.Sprintf("%s#%s", clientID, userName)
	ctx := context.TODO()
	// 缓存 token
	LocalAPIGatewayJWTClientIDCache.Set(ctx, key, val)
}