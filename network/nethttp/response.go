package nethttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/buildin/reflectutils"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseList struct {
	Total    int64       `json:"total"`
	Page     int64       `json:"page"`
	NumPages int64       `json:"numPages"`
	Results  interface{} `json:"results"`
}

type DebugResponse struct {
	Response
	Debug interface{} `json:"debug"`
}

// IsSuccess 判断响应是否成功返回结果
func (r Response) IsSuccess() bool {
	return r.Code == xgo.NoError
}

func JSONResponse(c *gin.Context, status int, code int, message string, data interface{}) {
	body := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	if code == xgo.NoError {
		body.Result = true
	} else {
		body.Result = false
	}
	c.JSON(status, body)
}

func ErrorJSONResponse(c *gin.Context, code int, message string) {
	JSONResponse(c, http.StatusOK, code, message, gin.H{})
}

func StatusForbiddenJSONResponse(c *gin.Context, message string, data interface{}) {
	JSONResponse(c, http.StatusForbidden, xgo.IAMNoPermission, message, data)
}

func SuccessRawJsonResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func SuccessJSONResponse(c *gin.Context, message string, data interface{}) {
	JSONResponse(c, http.StatusOK, xgo.NoError, message, data)
}

func SuccessJSONResponseWithDebug(c *gin.Context, message string, data interface{}, debug interface{}) {
	if reflectutils.IsNil(debug) {
		// 非 debug 模式
		SuccessJSONResponse(c, message, data)
		return
	}

	// debug 模式
	body := DebugResponse{
		Response: Response{
			Code:    xgo.NoError,
			Message: message,
			Data:    data,
			Result:  true,
		},
		Debug: debug,
	}

	c.JSON(http.StatusOK, body)
}

// NewResponse 创建响应对象，并转换为 []byte
func NewResponse(code int, message string, data interface{}) ([]byte, error) {
	result := false
	if 0 == code {
		result = true
	} else {
		appName := os.Args[0]
		szArr := strings.Split(appName, "/")
		if len(szArr) >= 2 {
			appName = szArr[1]
		}
		message = "(" + appName + "):" + message
	}

	resp := Response{result, code, message, data}
	b, err := json.Marshal(resp)
	if err != nil {
		return []byte(""), err
	}

	return b, nil
}

func NewErrorJSONResponseFunc(errorCode int, defaultMessage string) func(c *gin.Context, message string) {
	return func(c *gin.Context, message string) {
		msg := defaultMessage
		if message != "" {
			msg = fmt.Sprintf("%s:%s", msg, message)
		}
		ErrorJSONResponse(c, errorCode, msg)
	}
}

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

func SetError(c *gin.Context, err error) {
	c.Set(ErrorIDKey, err)
}

func SystemErrorJSONResponse(c *gin.Context, err error) {
	SetError(c, err)
	ErrorJSONResponse(c, xgo.SystemError, err.Error())
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

	body := DebugResponse{
		Response: Response{
			Result:  false,
			Code:    xgo.SystemError,
			Message: message,
			Data:    gin.H{},
		},
		Debug: debug,
	}

	c.JSON(http.StatusOK, body)
}
