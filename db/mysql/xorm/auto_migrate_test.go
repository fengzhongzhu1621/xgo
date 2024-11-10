package xorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSyncTable(t *testing.T) {
	dbClient := GetDefaultXormDBClient()
	defer dbClient.Close()

	dbClient.SyncTable(new(XormStudent))
	dbClient.SyncTable(new(XormUser))
	dbClient.SyncTable(new(XormPost))
	dbClient.SyncTable(new(XormUser2))
	dbClient.SyncTable(new(XormUser3))
	dbClient.SyncTable(new(XormCardS))
	dbClient.SyncTable(new(XormCardM))
}
