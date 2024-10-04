package utils

import (
	"fmt"
	"net/http"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/buildin/reflectutils"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/gin-gonic/gin"
)

func SetError(c *gin.Context, err error) {
	// c.Set方法用于在上下文（context）中设置一个键值对。这个方法允许你在处理请求的过程中存储和检索自定义数据。这些数据可以在同一个请求的不同处理函数之间共享。
	// c.Set方法设置的键值对仅在当前请求的上下文中有效。
	c.Set(nethttp.ErrorIDKey, err)
}

func SystemErrorJSONResponse(c *gin.Context, err error) {
	SetError(c, err)
	nethttp.ErrorJSONResponse(c, xgo.SystemError, err.Error())
}

func SystemErrorJSONResponseWithDebug(c *gin.Context, err error, debug interface{}) {
	if reflectutils.IsNil(debug) {
		// 非 debug 模式
		SystemErrorJSONResponse(c, err)
		return
	}

	// debug 模式
	message := fmt.Sprintf("system error[request_id=%s]: %s", GetRequestID(c), err.Error())
	SetError(c, err)

	body := nethttp.DebugResponse{
		Response: nethttp.Response{
			Result:  false,
			Code:    xgo.SystemError,
			Message: message,
			Data:    gin.H{},
		},
		Debug: debug,
	}

	c.JSON(http.StatusOK, body)
}
