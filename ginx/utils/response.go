package utils

import (
	"fmt"
	"net/http"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/network/constant"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/gin-gonic/gin"
)

var (
	BadRequestErrorJSONResponse   = NewErrorJSONResponseFunc(xgo.BadRequestError, "bad request")
	ParamErrorJSONResponse        = NewErrorJSONResponseFunc(xgo.ParamError, "param error")
	ForbiddenJSONResponse         = NewErrorJSONResponseFunc(xgo.ForbiddenError, "no permission")
	UnauthorizedJSONResponse      = NewErrorJSONResponseFunc(xgo.UnauthorizedError, "unauthorized")
	NotFoundJSONResponse          = NewErrorJSONResponseFunc(xgo.NotFoundError, "not found")
	ConflictJSONResponse          = NewErrorJSONResponseFunc(xgo.ConflictError, "conflict")
	TooManyRequestsJSONResponse   = NewErrorJSONResponseFunc(xgo.TooManyRequests, "too many requests")
	SamForbiddenJSONResponse      = NewErrorJSONResponseFunc(xgo.IAMNoPermission, "no permission")
	StaffUnauthorizedJSONResponse = NewErrorJSONResponseFunc(xgo.StaffUnauthorizedError, "unauthorized")
)

// SetError 在上下文中保存错误信息
func SetError(c *gin.Context, err error) {
	// c.Set方法用于在上下文（context）中设置一个键值对。这个方法允许你在处理请求的过程中存储和检索自定义数据。这些数据可以在同一个请求的不同处理函数之间共享。
	// c.Set方法设置的键值对仅在当前请求的上下文中有效。
	c.Set(constant.ErrorIDKey, err)
}

// JSONResponse 转换为标准的响应格式
func JSONResponse(c *gin.Context, status int, code int, message string, data interface{}) {
	if status == http.StatusNoContent {
		c.Status(status)
		return
	}

	body := nethttp.Response{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: GetRequestID(c),
	}
	if code == xgo.NoError {
		body.Result = true
	} else {
		body.Result = false
	}
	c.JSON(status, body)
}

// ErrorJSONResponse 返回错误响应（data 为空，http 状态码为 200）
func ErrorJSONResponse(c *gin.Context, code int, message string) {
	JSONResponse(c, http.StatusOK, code, message, nil)
}

// StatusForbiddenJSONResponse 返回无权限响应
func StatusForbiddenJSONResponse(c *gin.Context, message string, data interface{}) {
	JSONResponse(c, http.StatusForbidden, xgo.IAMNoPermission, message, data)
}

// SuccessRawJsonResponse 返回成功响应（不需要 message）
func SuccessRawJsonResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// SuccessJSONResponse 返回成功响应（code 为 0）
func SuccessJSONResponse(c *gin.Context, data interface{}) {
	message := ""
	JSONResponse(c, http.StatusOK, xgo.NoError, message, data)
}

// SystemErrorJSONResponse 返回系统报错
func SystemErrorJSONResponse(c *gin.Context, err error) {
	var (
		message string
	)
	SetError(c, err)

	// 构造错误信息
	cfg := config.GetGlobalConfig()
	if cfg.Debug {
		message = fmt.Sprintf("system error[request_id=%s]: %s", GetRequestID(c), err.Error())
	} else {
		message = err.Error()
	}

	ErrorJSONResponse(c, xgo.SystemError, message)
}

// NewErrorJSONResponseFunc 追加自定义错误信息
func NewErrorJSONResponseFunc(errorCode int, defaultMessage string) func(c *gin.Context, message string) {
	return func(c *gin.Context, message string) {
		msg := defaultMessage
		if message != "" {
			msg = fmt.Sprintf("%s:%s", msg, message)
		}
		ErrorJSONResponse(c, errorCode, msg)
	}
}

func NewPaginatedRespData(count int64, results any) nethttp.PaginatedResp {
	return nethttp.PaginatedResp{Count: count, Results: results}
}
