package nethttp

import (
	"github.com/fengzhongzhu1621/xgo"
)

type Response struct {
	Result    bool        `json:"result"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"requestID"`
}

type ResponseList struct {
	Total    int64       `json:"total"`
	Page     int64       `json:"page"`
	NumPages int64       `json:"num_pages"`
	Results  interface{} `json:"results"`
}

type PaginatedResp struct {
	Count   int64 `json:"count"`
	Results any   `json:"results"`
}

// IsSuccess 判断响应是否成功返回结果
func (r Response) IsSuccess() bool {
	return r.Code == xgo.NoError
}
