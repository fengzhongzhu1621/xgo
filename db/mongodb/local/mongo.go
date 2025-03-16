package local

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/config/server_option"
	"github.com/fengzhongzhu1621/xgo/db/db"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	// reference doc:
	// https://docs.mongodb.com/manual/core/read-preference-staleness/#replica-set-read-preference-max-staleness
	// this is the minimum value of maxStalenessSeconds allowed.
	// specifying a smaller maxStalenessSeconds value will raise an error. Clients estimate secondaries’ staleness
	// by periodically checking the latest write date of each replica set member. Since these checks are infrequent,
	// the staleness estimate is coarse. Thus, clients cannot enforce a maxStalenessSeconds value of less than
	// 90 seconds.
	maxStalenessSeconds = 90 * time.Second
)

// NewMgo returns new RDB
func NewMgo(config MongoConf, timeout time.Duration) (*Mongo, error) {
	connStr, err := connstring.Parse(config.URI)
	if nil != err {
		return nil, err
	}
	if config.RsName == "" {
		return nil, fmt.Errorf("mongodb rsName not set")
	}
	socketTimeout := time.Second * time.Duration(config.SocketTimeout)
	maxConnIdleTime := 25 * time.Minute
	appName := server_option.GetIdentification()
	// do not change this, our transaction plan need it to false.
	// it's related with the transaction number(eg txnNumber) in a transaction session.
	disableWriteRetry := false
	conOpt := options.ClientOptions{
		MaxPoolSize:     &config.MaxOpenConns,
		MinPoolSize:     &config.MaxIdleConns,
		ConnectTimeout:  &timeout,
		SocketTimeout:   &socketTimeout,
		ReplicaSet:      &config.RsName,
		RetryWrites:     &disableWriteRetry,
		MaxConnIdleTime: &maxConnIdleTime,
		AppName:         &appName,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI), &conOpt)
	if nil != err {
		return nil, err
	}

	if err := client.Connect(context.TODO()); nil != err {
		return nil, err
	}

	// initialize mongodb related metrics
	initMongoMetric()

	mgo := &Mongo{
		dbc:    client,
		dbname: connStr.Database,
		tm:     &TxnManager{},
		conf: &mongoCliConf{
			disableInsert: config.DisableInsert,
		},
	}

	return mgo, nil
}

func checkMongodbVersion(db string, client *mongo.Client) error {
	// 创建 serverStatus 命令
	command := bson.D{{Key: "serverStatus", Value: 1}}

	// 执行命令
	result := client.Database(db).RunCommand(context.Background(), command)

	// 解码结果到 bson.M
	var serverStatus bson.M
	if err := result.Decode(&serverStatus); err != nil {
		return err
	}

	// 提取版本信息
	versionStr, ok := serverStatus["version"]
	if !ok {
		return errors.New("version field not found in serverStatus")
	}

	// 分割版本字符串
	fields := strings.Split(versionStr.(string), ".")
	if len(fields) < 2 {
		return fmt.Errorf("got invalid mongodb version: %v", versionStr)
	}

	// 解析主版本号
	major, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("parse mongodb version %s major failed, err: %v", versionStr, err)
	}
	if major < 4 {
		return errors.New("mongodb version must be >= v4.0.0")
	}

	// 解析次版本号，如果存在
	var minor int
	if len(fields) > 1 {
		minor, err = strconv.Atoi(fields[1])
		if err != nil {
			return fmt.Errorf("parse mongodb version %s minor failed, err: %v", versionStr, err)
		}
	} else {
		minor = 0
	}

	if major == 4 && minor < 2 {
		return errors.New("mongodb version must be >= v4.2.0")
	}

	return nil
}

// InitTxnManager TxnID management of initial transaction
func (c *Mongo) InitTxnManager(r *redis.Client) error {
	return c.tm.InitTxnManager(r)
}

// Close replica client
func (c *Mongo) Close() error {
	c.dbc.Disconnect(context.TODO())
	return nil
}

// Ping replica client
func (c *Mongo) Ping() error {
	return c.dbc.Ping(context.TODO(), nil)
}

// IsDuplicatedError check duplicated error
func (c *Mongo) IsDuplicatedError(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "The existing index") {
			return true
		}
		if strings.Contains(err.Error(), "There's already an index with name") {
			return true
		}
		if strings.Contains(err.Error(), "E11000 duplicate") {
			return true
		}
		if strings.Contains(err.Error(), "IndexOptionsConflict") {
			return true
		}
		if strings.Contains(err.Error(), "all indexes already exist") {
			return true
		}
		if strings.Contains(err.Error(), "already exists with a different name") {
			return true
		}
	}
	return err == db.ErrDuplicated
}

// HasTable 判断是否存在集合
func (c *Mongo) HasTable(ctx context.Context, collName string) (bool, error) {
	cursor, err := c.dbc.Database(c.dbname).ListCollections(ctx, bson.M{"name": collName, "type": "collection"})
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		return true, nil
	}

	return false, nil
}

// IsNotFoundError check the not found error
func (c *Mongo) IsNotFoundError(err error) bool {
	return err == db.ErrDocumentNotFound
}

// GetDBClient TODO
// get db client
func (c *Mongo) GetDBClient() *mongo.Client {
	return c.dbc
}

// GetDBName TODO
// get db name
func (c *Mongo) GetDBName() string {
	return c.dbname
}

// ListTables 获取所有的表名
func (c *Mongo) ListTables(ctx context.Context) ([]string, error) {
	return c.dbc.Database(c.dbname).ListCollectionNames(ctx, bson.M{"type": "collection"})

}

// DropTable 移除集合
func (c *Mongo) DropTable(ctx context.Context, collName string) error {
	return c.dbc.Database(c.dbname).Collection(collName).Drop(ctx)
}

// CreateTable 创建集合 TODO test
func (c *Mongo) CreateTable(ctx context.Context, collName string) error {
	return c.dbc.Database(c.dbname).RunCommand(ctx, map[string]interface{}{"create": collName}).Err()
}

// RenameTable 更新集合名称
func (c *Mongo) RenameTable(ctx context.Context, prevName, currName string) error {
	cmd := bson.D{
		{"renameCollection", c.dbname + "." + prevName},
		{"to", c.dbname + "." + currName},
	}
	return c.dbc.Database("admin").RunCommand(ctx, cmd).Err()
}

// Table collection operation
func (c *Mongo) Table(collName string) db.Table {
	return nil
}

// NextSequences 批量获取新序列号(非事务)
func (c *Mongo) NextSequences(ctx context.Context, sequenceName string, num int) ([]uint64, error) {
	return make([]uint64, 0), nil
}

// NextSequence 获取新序列号(非事务)
func (c *Mongo) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	doc := Idgen{}
	return doc.SequenceID, nil
}

// Idgen TODO
type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}
