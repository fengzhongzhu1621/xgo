package uuid

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/duke-git/lancet/v2/random"
	gofrsUUID "github.com/gofrs/uuid"
	"github.com/google/uuid"
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
	for key := range uuids {
		if _, ok := uniqueUUIDs[key]; ok {
			t.Error(key, " has duplicate")
		}
		uniqueUUIDs[key] = struct{}{}
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

func TestUUIdV1(t *testing.T) {
	idV1, _ := uuid.NewUUID()
	fmt.Println("版本1 UUID:", idV1) // ffe0e9aa-15ba-11f0-a84f-c6eda4bffd12
}

func TestUUIdV3(t *testing.T) {
	// 预定义命名空间(DNS)
	nsDNS := uuid.NameSpaceDNS

	// 生成版本3 UUID(基于MD5)
	u3 := uuid.NewMD5(nsDNS, []byte("example.com"))
	fmt.Println("版本3 UUID:", u3)
}

// TestUUIdV4 Generate a random UUID of version 4 according to RFC 4122.
// func UUIdV4() (string, error)
func TestUUIdV4(t *testing.T) {
	{
		uuidValue, err := random.UUIdV4()
		if err != nil {
			return
		}
		fmt.Println(uuidValue) // c746705a-860f-46cf-a117-ef996fc4defe}
	}

	{
		id := uuid.New()
		fmt.Println("生成的UUID:", id) // 77dd4d61-d08c-4eee-8cd6-9e15a7d3d19d
	}

	{
		// 创建版本4 UUID
		u, _ := gofrsUUID.NewV4()
		fmt.Println("生成的UUID:", u)
	}
}

func TestUUIdV5(t *testing.T) {
	// 预定义命名空间(DNS)
	nsDNS := uuid.NameSpaceDNS

	// 生成版本5 UUID(基于SHA-1)
	u5 := uuid.NewSHA1(nsDNS, []byte("example.com"))
	fmt.Println("版本5 UUID:", u5)
}

func TestGoogleParseUUID(t *testing.T) {
	// 从字符串解析UUID
	parsedUUID, err := uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}
	fmt.Println("解析后的UUID:", parsedUUID) // 6ba7b810-9dad-11d1-80b4-00c04fd430c8
}
