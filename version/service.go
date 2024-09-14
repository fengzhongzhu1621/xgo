package version

import (
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"go.uber.org/zap"
)

// 版本配置服务
type versionsConfigService struct {
	rootDir string
	logger  *zap.Logger
}

var _ IVersionsConfigService = (*versionsConfigService)(nil)

// getVersionLogs 格式化语言并获得指定目录下指定版本和语言的版本日志
func getVersionLogs(rootDir string, language, versionStr string) (latestVer string, versionLogs []*VersionLog, err error) {
	if language == nethttp.WEB_LANGUAGE_CN {
		language = "zh"
	} else {
		language = "en"
	}
	// 获得指定目录下指定版本和语言的版本日志
	return ListChangelogs(rootDir, language, versionStr)
}

// GetVersionContent 获得版本日志内容
func (s *versionsConfigService) GetVersionContent(param *GetVersionContentRequestParam) (string, error) {
	var (
		err         error
		versionlogs []*VersionLog
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
func (s *versionsConfigService) ListVersions(language string) (*ListVersionsResponse, error) {
	var (
		err         error
		versionlogs []*VersionLog
		respData    ListVersionsResponse
		latestVer   string
	)

	if latestVer, versionlogs, err = getVersionLogs(s.rootDir, language, ""); err != nil {
		s.logger.Error("ListVersions ListBkSamChangelogs error", zap.Any("errmsg", err))
		return nil, err
	}

	respData.LastVersion = latestVer
	for _, v := range versionlogs {
		respData.VersionLogs = append(respData.VersionLogs, VersionLog{
			Version:   v.Version,
			ReleaseAt: v.ReleaseAt,
			Content:   "",
		})
	}

	return &respData, nil
}

// func NewVersionsConfigService(cfg *config.Config) *VersionsConfigService {
// 	obj := versionsConfigService{
// 		rootDir: cfg.RootDir,
// 		logger:  logging.GetWebLogger(),
// 	}

// 	return &obj
// }
