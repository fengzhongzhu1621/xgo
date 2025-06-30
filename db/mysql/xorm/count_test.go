package xorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	var (
		err error
	)
	dbClient := GetDefaultXormDBClient()
	xormEngine := dbClient.DB

	cardM := &XormCardM{}
	err = xormEngine.Sync2(cardM)
	if err != nil {
		panic(err)
	}

	count, err := xormEngine.Where("name = ?", "bob").Count(cardM)
	assert.Equal(t, true, count >= 0)
}
