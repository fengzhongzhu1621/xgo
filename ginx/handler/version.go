package handler

import (
	"strings"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/ginx/serializer"
	"github.com/fengzhongzhu1621/xgo/ginx/service"
	"github.com/fengzhongzhu1621/xgo/ginx/service/types"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/gin-gonic/gin"
)

// GetVersionContent 获得版本日志内容
func GetVersionContent(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err                          error
			listObj                      serializer.GetVersionContentSerializer
			getVersionContentRequestData *types.GetVersionContentRequestData
			language                     string
		)

		// 从请求参数中获得版本号
		// 将请求的查询参数绑定到一个结构体实例上。这个方法会自动处理查询参数的类型转换，并将结果存储在结构体的相应字段中。
		if err = c.ShouldBindQuery(&listObj); err != nil {
			nethttp.BadRequestErrorJSONResponse(c, validator.ValidationErrorMessage(err))
			return
		}

		// 从请求头获取语言
		language, err = nethttp.GetBluekingLanguageFromHeader(c)
		if err != nil || language == "" {
			language = nethttp.WEB_LANGUAGE_CN
		}

		// 创建版本服务实例
		svc := service.NewVersionsConfigService(cfg)

		// 构造服务请求参数
		getVersionContentRequestData = &types.GetVersionContentRequestData{
			Version:  listObj.Version,
			Operator: nethttp.GetUsernameFromHeader(c),
			Language: strings.ToLower(language),
		}

		// 调用服务获得版本日志内容
		versionContent, err := svc.GetVersionContent(getVersionContentRequestData)
		if err != nil {
			utils.SystemErrorJSONResponse(c, err)
			return
		}

		nethttp.SuccessJSONResponse(c, versionContent)
	}
}

// ListVersions 获得指定语言的所有的版本记录（不包含版本记录的内容，只包含版本的名称或最新的版本号）
func ListVersions(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err              error
			listVersionsResp *types.ListVersionsResponseData
			language         string
		)

		// 从请求头获取语言
		language, err = nethttp.GetBluekingLanguageFromHeader(c)
		if err != nil || language == "" {
			language = nethttp.WEB_LANGUAGE_CN
		}

		// 创建版本服务实例
		svc := service.NewVersionsConfigService(cfg)

		listVersionsResp, err = svc.ListVersions(strings.ToLower(language))
		if err != nil {
			utils.SystemErrorJSONResponse(c, err)
			return
		}

		nethttp.SuccessJSONResponse(c, *listVersionsResp)
	}
}
