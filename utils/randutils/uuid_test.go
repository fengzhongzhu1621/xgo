package randutils

import (
	"sync"
	"testing"
)

func testuUniqness(t *testing.T, genFunc func() string) {
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
	testuUniqness(t, NewUUID)
}

func TestShortUUID(t *testing.T) {
	testuUniqness(t, NewShortUUID)
}

func TestULID(t *testing.T) {
	testuUniqness(t, NewULID)
}
