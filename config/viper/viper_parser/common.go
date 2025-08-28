package viper_parser

import (
	"bytes"

	log "github.com/sirupsen/logrus"
)

func getCommonParser() *viperParser {
	if commonParser != nil {
		return commonParser
	}

	return localParser
}

func SetCommonFromFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	commonParser, err = newViperParserFromFile(target)
	if err != nil {
		log.Errorf("fail to read configure from common")
		return err
	}
	return nil
}

func SetCommonFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// if it is not nil, do not create a new parser, but add the new configuration information to viper
	if commonParser != nil {
		err = commonParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			log.Errorf("fail to read configure from common")
			return err
		}
		return nil
	}
	commonParser, err = newViperParser(data)
	if err != nil {
		log.Errorf("fail to read configure from common")
		return err
	}
	return nil
}
