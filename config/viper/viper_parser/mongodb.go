package viper_parser

import (
	"bytes"
	"errors"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/mongodb"
	mongo "github.com/fengzhongzhu1621/xgo/db/mongodb"
	log "github.com/sirupsen/logrus"
)

// Mongo return mongo configuration information according to the prefix.
func Mongo(prefix string, waitTime int) (mongodb.Config, error) {
	confLock.RLock()
	defer confLock.RUnlock()
	var parser *viperParser
	for sleepCnt := 0; sleepCnt < waitTime; sleepCnt++ {
		parser = getMongodbParser()
		if parser != nil {
			break
		}
		log.Warn("the configuration of mongo is not ready yet")
		time.Sleep(time.Duration(1) * time.Second)
	}

	if parser == nil {
		log.Errorf("can't find mongo configuration")
		return mongodb.Config{}, errors.New("can't find mongo configuration")
	}

	c := mongo.Config{
		Address:   parser.getString(prefix + ".host"),
		Port:      parser.getString(prefix + ".port"),
		User:      parser.getString(prefix + ".usr"),
		Password:  parser.getString(prefix + ".pwd"),
		Database:  parser.getString(prefix + ".database"),
		Mechanism: parser.getString(prefix + ".mechanism"),
		RsName:    parser.getString(prefix + ".rsName"),
	}

	if c.RsName == "" {
		log.Errorf("rsName not set")
	}
	if c.Mechanism == "" {
		c.Mechanism = "SCRAM-SHA-1"
	}

	maxOpenConns := prefix + ".maxOpenConns"
	if !parser.isSet(maxOpenConns) {
		log.Errorf("can not find config %s, set default value: %d", maxOpenConns, mongo.DefaultMaxOpenConns)
		c.MaxOpenConns = mongo.DefaultMaxOpenConns
	} else {
		c.MaxOpenConns = parser.getUint64(maxOpenConns)
	}

	if c.MaxIdleConns > mongo.MaximumMaxOpenConns {
		log.Errorf("config %s exceeds maximum value, use maximum value %d", maxOpenConns, mongo.MaximumMaxOpenConns)
		c.MaxIdleConns = mongo.MaximumMaxOpenConns
	}

	maxIdleConns := prefix + ".maxIdleConns"
	if !parser.isSet(maxIdleConns) {
		log.Errorf("can not find config %s, set default value: %d", maxIdleConns, mongo.MinimumMaxIdleOpenConns)
		c.MaxIdleConns = mongo.MinimumMaxIdleOpenConns
	} else {
		c.MaxIdleConns = parser.getUint64(maxIdleConns)
	}

	if c.MaxIdleConns < mongo.MinimumMaxIdleOpenConns {
		log.Errorf("config %s less than minimum value, use minimum value %d",
			maxIdleConns, mongo.MinimumMaxIdleOpenConns)
		c.MaxIdleConns = mongo.MinimumMaxIdleOpenConns
	}

	if !parser.isSet(prefix + ".socketTimeoutSeconds") {
		log.Errorf("can not find mongo.socketTimeoutSeconds config, use default value: %d",
			mongo.DefaultSocketTimeout)
		c.SocketTimeout = mongo.DefaultSocketTimeout
		return c, nil
	}

	c.SocketTimeout = parser.getInt(prefix + ".socketTimeoutSeconds")
	if c.SocketTimeout > mongo.MaximumSocketTimeout {
		log.Errorf("mongo.socketTimeoutSeconds config %d exceeds maximum value, use maximum value %d",
			c.SocketTimeout, mongo.MaximumSocketTimeout)
		c.SocketTimeout = mongo.MaximumSocketTimeout
	}

	if c.SocketTimeout < mongo.MinimumSocketTimeout {
		log.Errorf("mongo.socketTimeoutSeconds config %d less than minimum value, use minimum value %d",
			c.SocketTimeout, mongo.MinimumSocketTimeout)
		c.SocketTimeout = mongo.MinimumSocketTimeout
	}

	return c, nil
}

func getMongodbParser() *viperParser {
	if mongodbParser != nil {
		return mongodbParser
	}

	return localParser
}

// SetMongodbFromByte TODO
func SetMongodbFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()

	if mongodbParser != nil {
		err = mongodbParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			log.Errorf("fail to read configure from mongodb")
			return err
		}
		return nil
	}
	mongodbParser, err = newViperParser(data)
	if err != nil {
		log.Errorf("fail to read configure from mongodb")
		return err
	}

	return nil
}

// SetMongodbFromFile TODO
func SetMongodbFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()

	mongodbParser, err = newViperParserFromFile(target)
	if err != nil {
		log.Errorf("fail to read configure from mongodb")
		return err
	}
	return nil
}
