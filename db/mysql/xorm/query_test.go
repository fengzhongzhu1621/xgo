package xorm

import (
	"fmt"
	"testing"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestFind(t *testing.T) {
	var (
		err        error
		findResult = make([]XormCardM, 0)
	)
	dbClient := GetDefaultXormDBClient()
	xormEngine := dbClient.DB

	cardM := &XormCardM{}
	err = xormEngine.Sync2(cardM)
	if err != nil {
		panic(err)
	}

	cardM.Addr = map[string]interface{}{"street": "1778 Main st", "city": "Los Angeles"}
	cardM.Id = fmt.Sprintf("something%d", rand.Intn(100000000))
	cardM.Nickname = "test"
	cardM.NumberLast4 = "1111"

	// 获得多条数据
	if err = xormEngine.Cols("number_last_4").Where("nickname=?", "test").Find(&findResult); err != nil {
		panic(err)
	}
	expect := `[
	{
		"addr": null,
		"id": "",
		"nickname": "",
		"numberLast4": "1111"
	}
]`
	assert.Equal(t, expect, tests.ToString(findResult))
}

func TestGet(t *testing.T) {
	var (
		err       error
		has       bool
		getResult = new(XormCardM)
	)
	dbClient := GetDefaultXormDBClient()
	xormEngine := dbClient.DB

	cardM := &XormCardM{}
	err = xormEngine.Sync2(cardM)
	if err != nil {
		panic(err)
	}

	cardM.Addr = map[string]interface{}{"street": "1778 Main st", "city": "Los Angeles"}
	cardM.Id = fmt.Sprintf("something%d", rand.Intn(100000000))
	cardM.Nickname = "test"
	cardM.NumberLast4 = "1111"

	if has, err = xormEngine.Where("nickname = ? AND number_last_4 = ?", cardM.Nickname, cardM.NumberLast4).Get(getResult); err != nil {
		panic(err)
	}
	expect := `{
	"addr": {
		"city": "Los Angeles",
		"street": "1778 Main st"
	},
	"id": "something11989351",
	"nickname": "test",
	"numberLast4": "1111"
}`
	assert.Equal(t, true, has)
	assert.Equal(t, expect, tests.ToString(getResult))
}
