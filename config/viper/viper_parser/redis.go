package viper_parser

import (
	"bytes"
	"errors"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/redis/client"
	redis "github.com/fengzhongzhu1621/xgo/db/redis/client"
	log "github.com/sirupsen/logrus"
)

var redisParser *viperParser

// Redis return redis configuration information according to the prefix.
func Redis(prefix string, waitTime int) (redis.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	var parser *viperParser
	for sleepCnt := 0; sleepCnt < waitTime; sleepCnt++ {
		parser = getRedisParser()
		if parser != nil {
			break
		}
		log.Warn("the configuration of redis is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		log.Errorf("can't find redis configuration")
		return client.Config{}, errors.New("can't find redis configuration")
	}

	return client.Config{
		Address:          parser.getString(prefix + ".host"),
		Password:         parser.getString(prefix + ".pwd"),
		Database:         parser.getString(prefix + ".database"),
		MasterName:       parser.getString(prefix + ".masterName"),
		SentinelPassword: parser.getString(prefix + ".sentinelPwd"),
		Enable:           parser.getString(prefix + ".enable"),
		MaxOpenConns:     parser.getInt(prefix + ".maxOpenConns"),
	}, nil
}

func getRedisParser() *viperParser {
	if redisParser != nil {
		return redisParser
	}

	return localParser
}

func SetRedisFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	if redisParser != nil {
		err := redisParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			log.Errorf("fail to read configure from redis")
			return err
		}
		return nil
	}
	redisParser, err = newViperParser(data)
	if err != nil {
		log.Errorf("fail to read configure from redis")
		return err
	}
	return nil
}

func SetRedisFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	redisParser, err = newViperParserFromFile(target)
	if err != nil {
		log.Errorf("fail to read configure from redis")
		return err
	}
	return nil
}
