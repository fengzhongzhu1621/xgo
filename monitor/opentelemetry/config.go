package opentelemetry

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/config/viper/viper_parser"
	"github.com/fengzhongzhu1621/xgo/network/ssl"
	log "github.com/sirupsen/logrus"
)

var openTelemetryCfg = new(OpenTelemetryConfig)

// OpenTelemetryConfig TODO
type OpenTelemetryConfig struct {
	// 表示是否开启openTelemetry跟踪链接入相关功能，布尔值, 默认值为false不开启
	enable bool
	// openTelemetry跟踪链功能的自定义上报服务地址
	endpoint string

	tlsConf *tls.Config
}

// InitOpenTelemetryConfig init openTelemetry config
func InitOpenTelemetryConfig() error {
	var err error
	maxCnt := 100
	cnt := 0
	for !viper_parser.IsExist("openTelemetry") && cnt < maxCnt {
		log.Infof("waiting openTelemetry config to be init")
		cnt++
		time.Sleep(time.Millisecond * 300)
	}

	if cnt == maxCnt {
		return errors.New("no openTelemetry config is found, the config 'openTelemetry' must exist")
	}

	openTelemetryCfg.enable, err = viper_parser.Bool("openTelemetry.enable")
	if err != nil {
		return fmt.Errorf("config openTelemetry.enable err: %v", err)
	}

	// 如果不需要开启OpenTelemetry，那么后续没有必要再检查配置
	if !openTelemetryCfg.enable {
		return nil
	}

	openTelemetryCfg.endpoint, err = viper_parser.String("openTelemetry.endpoint")
	if err != nil {
		return fmt.Errorf("config openTelemetry.endpoint err: %v", err)
	}

	if !viper_parser.IsExist("openTelemetry.tls.caFile") ||
		!viper_parser.IsExist("openTelemetry.tls.certFile") ||
		!viper_parser.IsExist("openTelemetry.tls.keyFile") {

		return nil
	}

	caFile, err := viper_parser.String("openTelemetry.tls.caFile")
	if err != nil {
		return fmt.Errorf("get openTelemetry.tls.caFile error: %v", err)
	}

	certFile, err := viper_parser.String("openTelemetry.tls.certFile")
	if err != nil {
		return fmt.Errorf("get openTelemetry.tls.certFile error: %v", err)
	}

	keyFile, err := viper_parser.String("openTelemetry.tls.keyFile")
	if err != nil {
		return fmt.Errorf("get openTelemetry.tls.keyFile error: %v", err)
	}

	insecureSkipVerify, err := viper_parser.Bool("openTelemetry.tls.insecureSkipVerify")
	if err != nil {
		return fmt.Errorf("get openTelemetry.tls.insecureSkipVerify error: %v", err)
	}

	var password string
	if viper_parser.IsExist("openTelemetry.tls.password") {
		password, err = viper_parser.String("openTelemetry.tls.password")
		if err != nil {
			return fmt.Errorf("get openTelemetry.tls.password error: %v", err)
		}
	}

	tls, err := ssl.ClientTslConfVerity(caFile, certFile, keyFile, password)
	if err != nil {
		return err
	}
	tls.InsecureSkipVerify = insecureSkipVerify

	openTelemetryCfg.tlsConf = tls

	return nil
}
