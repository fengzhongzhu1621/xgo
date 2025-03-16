package viper_parser

import (
	"bytes"

	log "github.com/sirupsen/logrus"
)

func SetExtraFromByte(data []byte) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// if it is not nil, do not create a new parser, but add the new configuration information to viper
	if extraParser != nil {
		err = extraParser.parser.ReadConfig(bytes.NewBuffer(data))
		if err != nil {
			log.Errorf("fail to read configure from extra")
			return err
		}
		return nil
	}
	extraParser, err = newViperParser(data)
	if err != nil {
		log.Errorf("fail to read configure from extra")
		return err
	}
	return nil
}
