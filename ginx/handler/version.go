package handler

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/ginx/request/header"
	"github.com/fengzhongzhu1621/xgo/ginx/serializer"
	"github.com/fengzhongzhu1621/xgo/ginx/service"
	"github.com/fengzhongzhu1621/xgo/ginx/service/types"
	"github.com/fengzhongzhu1621/xgo/ginx/utils"
	"github.com/fengzhongzhu1621/xgo/network/nethttp"
	"github.com/fengzhongzhu1621/xgo/validator"
	"github.com/fengzhongzhu1621/xgo/version"
	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	runEnv := os.Getenv("RUN_ENV")
	now := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"module":    server_option.GetIdentification(),
		"version":   version.AppVersion,
		"commit":    version.GitCommit,
		"buildTime": version.BuildTime,
		"goVersion": version.GoVersion,
		"env":       runEnv,
		"timestamp": now.Unix(),
		"date":      now,
	})
}

// @Summary 获得版本日志内容
// @Description 获得版本日志内容
// @Router /api/v1/versions/ [get]
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
			utils.BadRequestErrorJSONResponse(c, validator.ValidationErrorMessage(err))
			return
		}

		// 从请求头获取语言
		language, err = header.GetLanguageFromHeader(c)
		if err != nil || language == "" {
			language = nethttp.WEB_LANGUAGE_CN
		}

		// 创建版本服务实例
		svc := service.NewVersionsConfigService(cfg)

		// 构造服务请求参数
		getVersionContentRequestData = &types.GetVersionContentRequestData{
			Version:  listObj.Version,
			Operator: header.GetUsernameFromHeader(c),
			Language: strings.ToLower(language),
		}

		// 调用服务获得版本日志内容
		versionContent, err := svc.GetVersionContent(getVersionContentRequestData)
		if err != nil {
			utils.SystemErrorJSONResponse(c, err)
			return
		}

		utils.SuccessJSONResponse(c, versionContent)
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
		language, err = header.GetLanguageFromHeader(c)
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

		utils.SuccessJSONResponse(c, *listVersionsResp)
	}
}
