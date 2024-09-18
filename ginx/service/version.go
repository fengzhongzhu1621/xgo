package service

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/service/types"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/fengzhongzhu1621/xgo/version"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"go.uber.org/zap"
)

// 版本配置服务
type versionsConfigService struct {
	rootDir string
	logger  *zap.Logger
}

var _ types.IVersionsConfigService = (*versionsConfigService)(nil)

// getVersionLogs 格式化语言并获得指定目录下指定版本和语言的版本日志
func getVersionLogs(rootDir string, language, versionStr string) (latestVer string, versionLogs []*version.VersionLog, err error) {
	if language == nethttp.WEB_LANGUAGE_CN {
		language = "zh"
	} else {
		language = "en"
	}
	// 获得指定目录下指定版本和语言的版本日志
	return version.ListChangelogs(rootDir, language, versionStr)
}

// GetVersionContent 获得版本日志内容
func (s *versionsConfigService) GetVersionContent(param *types.GetVersionContentRequestData) (string, error) {
	var (
		err         error
		versionlogs []*version.VersionLog
		outputHtml  []byte
	)

	// 获得指定目录下指定版本和语言的版本日志
	if _, versionlogs, err = getVersionLogs(s.rootDir, param.Language, param.Version); err != nil {
		s.logger.Error("GetVersionContentConfig ListBkSamChangelogs error", zap.Any("errmsg", err))
		return "", err
	}

	for _, v := range versionlogs {
		if v.Version == param.Version {
			// 转换为Markdown格式
			unsafe := blackfriday.MarkdownCommon([]byte(v.Content))
			// 对其进行安全处理
			outputHtml = bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			break
		}
	}

	// 返回版本日志的内容
	return string(outputHtml), nil
}

// ListVersions 获得指定语言的所有的版本记录（不包含版本记录的内容，只包含版本的名称或最新的版本号）
func (s *versionsConfigService) ListVersions(language string) (*types.ListVersionsResponseData, error) {
	var (
		err         error
		versionlogs []*version.VersionLog
		respData    types.ListVersionsResponseData
		latestVer   string
	)

	// 格式化语言并获得指定目录下指定版本和语言的版本日志
	if latestVer, versionlogs, err = getVersionLogs(s.rootDir, language, ""); err != nil {
		s.logger.Error("ListVersions ListBkSamChangelogs error", zap.Any("errmsg", err))
		return nil, err
	}

	// 构造响应输出
	respData.LastVersion = latestVer
	for _, v := range versionlogs {
		respData.VersionLogs = append(respData.VersionLogs, version.VersionLog{
			Version:   v.Version,
			ReleaseAt: v.ReleaseAt,
			Content:   "",
		})
	}

	return &respData, nil
}

// NewVersionsConfigService 版本配置服务
func NewVersionsConfigService(cfg *config.Config) *versionsConfigService {
	obj := versionsConfigService{
		rootDir: cfg.RootDir,
		logger:  logging.GetWebLogger(),
	}

	return &obj
}
