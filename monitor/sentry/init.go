package sentry

import (
	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

func InitSentry() {
	globalConfig := config.GetGlobalConfig()
	enable := globalConfig.Sentry.Enable
	dsn := globalConfig.Sentry.DSN
	if enable {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		})
		if err != nil {
			log.Errorf("Init Sentry fail: %v", err)
			return
		}
		log.Info("Init Sentry success")
	} else {
		log.Warn("Sentry is not enabled, will not init it")
	}

	SetSentryCapture(enable)
}
