package viper_parser

import (
	"errors"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/kafka"
	log "github.com/sirupsen/logrus"
)

// Kafka return kafka configuration information according to the prefix.
func Kafka(prefix string, waitTime int) (kafka.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	var parser *viperParser
	for sleepCnt := 0; sleepCnt < waitTime; sleepCnt++ {
		parser = getCommonParser()
		if parser != nil {
			break
		}
		log.Warn("the configuration of common is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		log.Errorf("can't find kafka configuration")
		return kafka.Config{}, errors.New("can't find kafka configuration")
	}

	return kafka.Config{
		Brokers:   parser.getStringSlice(prefix + ".brokers"),
		GroupID:   parser.getString(prefix + ".groupID"),
		Topic:     parser.getString(prefix + ".topic"),
		Partition: parser.getInt64(prefix + ".partition"),
		User:      parser.getString(prefix + ".user"),
		Password:  parser.getString(prefix + ".password"),
	}, nil
}
