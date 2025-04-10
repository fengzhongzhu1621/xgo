package uuid

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/duke-git/lancet/v2/random"
	"github.com/stretchr/testify/assert"
)

func testUniqness(t *testing.T, genFunc func() string) {
	producers := 100
	uuidsPerProducer := 10000

	if testing.Short() {
		producers = 10
		uuidsPerProducer = 1000
	}

	uuidsCount := producers * uuidsPerProducer

	// 并发产生uuid到管道中
	uuids := make(chan string, uuidsCount)
	allGenerated := sync.WaitGroup{}
	allGenerated.Add(producers)
	for i := 0; i < producers; i++ {
		go func() {
			for j := 0; j < uuidsPerProducer; j++ {
				uuids <- genFunc()
			}
			allGenerated.Done()
		}()
	}
	allGenerated.Wait()
	close(uuids)

	// 判断是否产生重复的数据
	uniqueUUIDs := make(map[string]struct{}, uuidsCount)
	for uuid := range uuids {
		if _, ok := uniqueUUIDs[uuid]; ok {
			t.Error(uuid, " has duplicate")
		}
		uniqueUUIDs[uuid] = struct{}{}
	}
}

func TestUUID(t *testing.T) {
	testUniqness(t, NewUUID4)
}

func TestShortUUID(t *testing.T) {
	testUniqness(t, NewShortUUID)
}

func TestULID(t *testing.T) {
	testUniqness(t, NewULID)
}

func TestMD5Hash(t *testing.T) {
	assert.Equal(t, "098f6bcd4621d373cade4e832627b4f6", MD5Hash("test"))
}

func TestMd5(t *testing.T) {
	src := "123456789"
	actual := Md5(src)
	expect := "25f9e794323b453885f5181f1b624d0b"
	assert.Equal(t, expect, actual)
}

func TestGenerateId(t *testing.T) {
	actual := GenerateID()
	s, err := strconv.ParseUint(actual, 10, 64)
	assert.Equal(t, err, nil)
	assert.Equal(t, s > 0, true)
}

// TestUUIdV4 Generate a random UUID of version 4 according to RFC 4122.
// func UUIdV4() (string, error)
func TestUUIdV4(t *testing.T) {
	{
		uuid, err := random.UUIdV4()
		if err != nil {
			return
		}
		fmt.Println(uuid) // c746705a-860f-46cf-a117-ef996fc4defe}
	}

}
