package nethttp

import (
	"fmt"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/iam"
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

type BaseResp struct {
	Result      bool               `json:"result" bson:"result" mapstructure:"result"`
	Code        int                `json:"bk_error_code" bson:"bk_error_code" mapstructure:"bk_error_code"`
	ErrMsg      string             `json:"bk_error_msg" bson:"bk_error_msg" mapstructure:"bk_error_msg"`
	Permissions *iam.IamPermission `json:"permission" bson:"permission" mapstructure:"permission"`
}

func (br *BaseResp) ToString() string {
	return fmt.Sprintf("code:%d, message:%s", br.Code, br.ErrMsg)
}

type JsonCntInfoResp struct {
	BaseResp
	Data CntInfoString `json:"data"`
}

// CntInfoString TODO
type CntInfoString struct {
	Count int64 `json:"count"`
	// info is a json array string field.
	Info string `json:"info"`
}

type JsonStringResp struct {
	BaseResp
	Data string
}

func NewNoPermissionResp(permission *iam.IamPermission) BaseResp {
	return BaseResp{
		Result:      false,
		Code:        NoPermission,
		ErrMsg:      "no permissions",
		Permissions: permission,
	}
}

type SuccessResponse struct {
	BaseResp `json:",inline"`
	Data     interface{} `json:"data"`
}

var SuccessBaseResp = BaseResp{Result: true, Code: Success, ErrMsg: SuccessStr}

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		BaseResp: SuccessBaseResp,
		Data:     data,
	}
}

// AbortTransactionResult abort transaction result
type AbortTransactionResult struct {
	// Retry defines if the transaction needs to retry, the following are the scenario that needs to retry:
	// 1. the write operation in the transaction conflicts with another transaction,
	// then do retry in the scene layer with server times depends on conditions.
	Retry bool `json:"retry"`
}

// AbortTransactionResponse abort transaction response
type AbortTransactionResponse struct {
	BaseResp               `json:",inline"`
	AbortTransactionResult `json:"data"`
}
