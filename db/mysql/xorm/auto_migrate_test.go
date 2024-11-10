package xorm

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSyncTable(t *testing.T) {
	dbClient := GetDefaultXormDBClient()
	dbClient.SyncTable(new(XormStudent))
}
