package nethttp

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/fengzhongzhu1621/xgo"
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

// NewResponse 创建响应对象，并转换为 []byte
func NewResponse(code int, message string, data interface{}) ([]byte, error) {
	result := false
	if code == 0 {
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
