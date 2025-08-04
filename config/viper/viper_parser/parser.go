package viper_parser

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	mongodbParser *viperParser
	commonParser  *viperParser
	extraParser   *viperParser
	localParser   *viperParser
)

var confLock sync.RWMutex

type viperParser struct {
	parser *viper.Viper
}

func newViperParser(data []byte) (*viperParser, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return &viperParser{parser: v}, nil
}

func newViperParserFromFile(target string) (*viperParser, error) {
	v := viper.New()
	v.SetConfigName(path.Base(target))
	v.AddConfigPath(path.Dir(target))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	v.WatchConfig()
	return &viperParser{parser: v}, nil
}

func (vp *viperParser) getString(path string) string {
	return vp.parser.GetString(path)
}

func (vp *viperParser) getInt(path string) int {
	return vp.parser.GetInt(path)
}

func (vp *viperParser) getUint64(path string) uint64 {
	return vp.parser.GetUint64(path)
}

func (vp *viperParser) getBool(path string) bool {
	return vp.parser.GetBool(path)
}

func (vp *viperParser) getDuration(path string) time.Duration {
	return vp.parser.GetDuration(path)
}

func (vp *viperParser) isSet(path string) bool {
	return vp.parser.IsSet(path)
}

func (vp *viperParser) getInt64(path string) int64 {
	return vp.parser.GetInt64(path)
}

func (vp *viperParser) getStringSlice(path string) []string {
	return vp.parser.GetStringSlice(path)
}

func (vp *viperParser) isConfigIntType(path string) bool {
	val := vp.parser.GetString(path)
	_, err := strconv.Atoi(val)
	return err == nil
}

func (vp *viperParser) isConfigBoolType(path string) bool {
	val := vp.parser.GetString(path)
	if val != "true" && val != "false" {
		return false
	}
	return true
}

func (vp *viperParser) unmarshalKey(key string, val interface{}) error {
	return vp.parser.UnmarshalKey(key, val)
}

// IsExist checks if key exists in all config files
func IsExist(key string) bool {
	confLock.RLock()
	defer confLock.RUnlock()

	// 在所有的配置文件中判断
	if (localParser == nil || !localParser.isSet(key)) &&
		(commonParser == nil || !commonParser.isSet(key)) &&
		(extraParser == nil || !extraParser.isSet(key)) {
		return false
	}
	return true
}

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func UnmarshalKey(key string, val interface{}) error {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return err
	}

	return parser.unmarshalKey(key, val)
}

// String return the string value of the configuration information according to the key.
func String(key string) (string, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return "", err
	}

	return parser.getString(key), nil
}

// Int return the int value of the configuration information according to the key.
func Int(key string) (int, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return 0, err
	}

	if !parser.isConfigIntType(key) {
		return 0, errors.New("config is not int type")
	}
	return parser.getInt(key), nil
}

// Int64 return the int value of the configuration information according to the key.
func Int64(key string) (int64, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return 0, err
	}

	if !parser.isConfigIntType(key) {
		return 0, errors.New("config is not int type")
	}
	return parser.getInt64(key), nil
}

// Bool return the bool value of the configuration information according to the key.
func Bool(key string) (bool, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return false, err
	}

	if !parser.isConfigBoolType(key) {
		return false, errors.New("config is not bool type")
	}
	return parser.getBool(key), nil
}

// StringSlice return the stringSlice value of the configuration information according to the key.
func StringSlice(key string) ([]string, error) {
	confLock.RLock()
	defer confLock.RUnlock()

	parser, err := getKeyValueParser(key)
	if err != nil {
		return nil, err
	}

	return parser.getStringSlice(key), nil
}

// getKeyValueParser get viper parser for common key value in the order of migrate->common->extra
func getKeyValueParser(key string) (*viperParser, error) {
	if localParser != nil && localParser.isSet(key) {
		return localParser, nil
	}

	if commonParser != nil && commonParser.isSet(key) {
		return commonParser, nil
	}

	if extraParser != nil && extraParser.isSet(key) {
		return extraParser, nil
	}

	return nil, fmt.Errorf("%s key's config not found", key)
}
