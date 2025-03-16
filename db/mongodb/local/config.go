package local

import (
	"github.com/fengzhongzhu1621/xgo/db/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	dbc    *mongo.Client
	dbname string
	sess   mongo.Session
	tm     *TxnManager
	conf   *mongoCliConf
}

// mongoCliConf is cmdb mongo client config
type mongoCliConf struct {
	// idGenStep is the step of id generator
	idGenStep int
	// disableInsert defines if insert operation for specific tables are disabled
	disableInsert bool
}

var _ db.DB = new(Mongo)

type MongoConf struct {
	TimeoutSeconds int
	MaxOpenConns   uint64
	MaxIdleConns   uint64
	URI            string
	RsName         string
	SocketTimeout  int
	DisableInsert  bool
}
