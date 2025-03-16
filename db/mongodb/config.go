package mongodb

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/db"
	"github.com/fengzhongzhu1621/xgo/db/mongodb/local"
)

const (
	// DefaultMaxOpenConns TODO
	// if maxOpenConns isn't configured, use default value
	DefaultMaxOpenConns = 1000
	// MaximumMaxOpenConns TODO
	// if maxOpenConns exceeds maximum value, use maximum value
	MaximumMaxOpenConns = 3000
	// MinimumMaxIdleOpenConns TODO
	// if maxIDleConns is less than minimum value, use minimum value
	MinimumMaxIdleOpenConns = 50
	// DefaultSocketTimeout TODO
	// if timeout isn't configured, use default value
	DefaultSocketTimeout = 10
	// MaximumSocketTimeout TODO
	// if timeout exceeds maximum value, use maximum value
	MaximumSocketTimeout = 30
	// MinimumSocketTimeout TODO
	// if timeout less than the minimum value, use minimum value
	MinimumSocketTimeout = 5
)

// Config config
type Config struct {
	Connect       string
	Address       string
	User          string
	Password      string
	Port          string
	Database      string
	Mechanism     string
	MaxOpenConns  uint64
	MaxIdleConns  uint64
	RsName        string
	SocketTimeout int
	DisableInsert bool
}

// BuildURI return mongo uri according to  https://docs.mongodb.com/manual/reference/connection-string/
// format example: mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
func (c Config) BuildURI() string {
	if c.Connect != "" {
		return c.Connect
	}

	if !strings.Contains(c.Address, ":") && len(c.Port) > 0 {
		c.Address = c.Address + ":" + c.Port
	}

	c.User = url.QueryEscape(c.User)
	c.Password = url.QueryEscape(c.Password)
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authMechanism=%s", c.User, c.Password, c.Address, c.Database, c.Mechanism)
	return uri
}

// GetMongoConf TODO
func (c Config) GetMongoConf() local.MongoConf {
	return local.MongoConf{
		MaxOpenConns:  c.MaxOpenConns,
		MaxIdleConns:  c.MaxIdleConns,
		URI:           c.BuildURI(),
		RsName:        c.RsName,
		SocketTimeout: c.SocketTimeout,
		DisableInsert: c.DisableInsert,
	}
}

// GetMongoClient TODO
func (c Config) GetMongoClient() (db db.DB, err error) {
	mongoConf := local.MongoConf{
		MaxOpenConns:  c.MaxOpenConns,
		MaxIdleConns:  c.MaxIdleConns,
		URI:           c.BuildURI(),
		RsName:        c.RsName,
		SocketTimeout: c.SocketTimeout,
	}
	db, err = local.NewMgo(mongoConf, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	return
}
