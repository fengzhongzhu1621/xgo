package xorm

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"
)

func TestInsertOne(t *testing.T) {
	var err error
	dbClient := GetDefaultXormDBClient()
	xormEngine := dbClient.DB
	defer xormEngine.Close()

	cardS := &XormCardS{}
	cardM := &XormCardM{}
	err = xormEngine.Sync2(cardS, cardM)
	if err != nil {
		panic(err)
	}

	cardM.Addr = map[string]interface{}{"street": "1778 Main st", "city": "Los Angeles"}
	cardM.Id = fmt.Sprintf("something%d", rand.Intn(100000000))
	cardM.Nickname = "test"
	cardM.NumberLast4 = "1111"

	// Using the jsonb field as struct will succeed...
	// +------------------------------------------------+-------------------+----------+---------------+
	// | addr                                           | id                | nickname | number_last_4 |
	// +------------------------------------------------+-------------------+----------+---------------+
	// | {"city":"Los Angeles","street":"1778 Main st"} | something11989351 | test     | 1111          |
	// +------------------------------------------------+-------------------+----------+---------------+
	affected, err := xormEngine.InsertOne(cardM)
	if affected != 1 || err != nil {
		panic(
			fmt.Sprintf(
				"Failed to insert cardM, got (%d) row affected and/or error (%s)",
				affected,
				err,
			),
		)
	} else {
		fmt.Printf("****** Success, inserted (%d) rows\n", affected)
	}

	// This will fail... "no primary key for col addr"
	cardS.Addr = Address{"1778 Main st", "Los Angeles"}
	cardS.Id = fmt.Sprintf("something%d", rand.Intn(100000000))
	cardS.Nickname = "test"
	cardS.NumberLast4 = "1111"

	// +------------------------------------------------+-------------------+------------+----------+---------------+
	// | addr                                           | id                | is_default | nickname | number_last_4 |
	// +------------------------------------------------+-------------------+------------+----------+---------------+
	// | {"street":"1778 Main st","city":"Los Angeles"} | something33091121 |          0 | test     | 1111          |
	// +------------------------------------------------+-------------------+------------+----------+---------------+
	affected, err = xormEngine.InsertOne(cardS)
	if affected != 1 || err != nil {
		panic(
			fmt.Sprintf(
				"Failed to insert cardS, got (%d) row affected and/or error (%s)",
				affected,
				err,
			),
		)
	}
}

func TestInsert(t *testing.T) {
	dbClient := GetDefaultXormDBClient()
	engine := dbClient.DB
	defer engine.Close()

	_ = engine.Sync2(new(XormUser4))

	// 插入单个用户
	newUser := &XormUser4{
		Name:  "username_a",
		Age:   30,
		Email: "a@example.com",
	}
	affected, _ := engine.Insert(newUser)
	fmt.Printf("插入了 %d 条记录，新用户的 ID 是 %d\n", affected, newUser.Id)

	// 插入多个用户
	users := []XormUser4{
		{Name: "username_b", Age: 25, Email: "a@example.com"},
		{Name: "username_c", Age: 28, Email: "b@example.com"},
	}
	affected, _ = engine.Insert(&users)
	fmt.Printf("批量插入了 %d 条记录\n", affected)
}
