package types

import (
	"github.com/fengzhongzhu1621/xgo/version"
)

type ListVersionsResponseData struct {
	// 需要显示的版本
	ShowVersion bool `json:"show_version"`
	// 版本日志
	VersionLogs []version.VersionLog `json:"version_logs"`
	// 最新的版本
	LastVersion string `json:"last_version"`
}

type GetVersionContentRequestData struct {
	// 语言
	Language string `json:"language" form:"language"`
	// 版本号
	Version string `json:"version"  form:"version"  binding:"required"`
	// 操作人
	Operator string `json:"operator" form:"operator" binding:"required"`
}

type IVersionsConfigService interface {
	// 获得版本日志内容
	GetVersionContent(getVersionContentParam *GetVersionContentRequestData) (string, error)
	// 获得指定语言的所有的版本记录（不包含版本记录的内容，只包含版本的名称或最新的版本号）
	ListVersions(language string) (*ListVersionsResponseData, error)
}
