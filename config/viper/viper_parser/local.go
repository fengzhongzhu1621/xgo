package viper_parser

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// SetLocalFile set localParser from file
func SetLocalFile(target string) error {
	var err error
	confLock.Lock()
	defer confLock.Unlock()
	// /data/migrate.yaml -> /data/migrate
	split := strings.Split(target, ".")
	filePath := split[0]
	localParser, err = newViperParserFromFile(filePath)
	if err != nil {
		log.Errorf("set local file failed, target: %s, err: %v", target, err)
		return err
	}
	return nil
}
