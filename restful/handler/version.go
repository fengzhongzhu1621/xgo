package handler

import (
	"os"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/version"
)

func Version(req *restful.Request, resp *restful.Response) {
	runEnv := os.Getenv("RUN_ENV")
	now := time.Now()
	answer := map[string]any{
		"module":    server_option.GetIdentification(),
		"version":   version.AppVersion,
		"commit":    version.GitCommit,
		"buildTime": version.BuildTime,
		"goVersion": version.GoVersion,
		"env":       runEnv,
		"timestamp": now.Unix(),
		"date":      now,
	}
	resp.WriteJson(answer, "application/json")
}
