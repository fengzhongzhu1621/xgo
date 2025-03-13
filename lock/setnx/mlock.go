package setnx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
)

type mlock struct {
	cache *redis.Client
	keys  []string

	needUnlock bool // 是否需要释放key
	isFirst    bool
}

type MLocker interface {
	MLock(rid string, retryTimes int, expire time.Duration, values ...string) (locked bool, err error)
	MUnlock() error
	MPartialUnlock(keys []string) error
}

func NewMLocker(cache *redis.Client) MLocker {
	return &mlock{
		isFirst:    false,
		cache:      cache,
		keys:       []string{},
		needUnlock: false,
	}
}

func (l *mlock) MLock(rid string, retry int, expire time.Duration, keys ...string) (locked bool, err error) {
	if l.isFirst {
		return false, errors.New("repeat lock")
	}

	var (
		bResultFlag bool
		delKeys     []string
	)
	pipeRes := make(map[bool][]string)
	l.isFirst = true

	for i := 0; i < retry; i++ {
		// 加分布式锁
		ctx := context.Background()
		bPipeResultFlag := false
		delKeys = []string{}
		pipe := l.cache.TxPipeline()
		l.keys = []string{}
		for _, key := range keys {
			uuid := xid.New().String()
			l.keys = append(l.keys, key)
			// 加共享锁
			pipe.SetNX(ctx, key, uuid, 0)
			// 设置 key 的过期时间
			pipe.Expire(ctx, key, expire)
		}
		res, err := pipe.Exec(ctx)
		if err != nil {
			// exec error try it again
			continue
		}

		for k, r := range res {
			// mlock contain setnx and expire two commonds, you should only mark setnx if success or not
			if k%2 == 0 {
				key, bResult := getExecSetNxBoolResult(r.String())
				if !bResult {
					// obtain lock fail ones
					pipeRes[false] = append(pipeRes[false], key)
				} else {
					// obtain lock success ones
					bPipeResultFlag = true
					pipeRes[true] = append(pipeRes[true], key)
				}
			}
		}

		// 部分失败，则回滚成功的 key
		if bPipeResultFlag && len(pipeRes[false]) > 0 {
			// if some setnx fail, need release the success ones
			err := l.cache.Del(context.Background(), pipeRes[true]...).Err()
			if err != nil {
				// if del fail, need to del it when unlock
				for _, v := range pipeRes[true] {
					delKeys = append(delKeys, v)
				}
				log.Errorf("delete key fail. the key: %v,rid: %s", pipeRes[true], rid)
			}
		} else {
			// 全部成功，则退出重试
			bResultFlag = true
			break
		}

		// 重置
		for k := range pipeRes {
			// release the key-lockResult  pair
			delete(pipeRes, k)
		}

		time.Sleep(100 * time.Millisecond)
	}

	if bResultFlag {
		// obtain lock success, you should release it
		l.needUnlock = true
		return true, nil
	}

	// delete fail,unlock retry delete it
	if len(delKeys) > 0 {
		l.needUnlock = true
		l.keys = delKeys // 需要调用 MUnlock 重试
	}

	// set map nil for gc
	pipeRes = nil
	return false, errors.New("obtain lock fail")
}

// getExecSetNxBoolResult 解析 Redis 事务执行结果，以确定 SETNX 命令是否成功。
// 然而，该函数存在一些设计和实现上的问题，可能导致其在实际应用中表现不稳定或难以维护。
// 第一个返回值是键名（string）。
// 第二个返回值是布尔值，表示 SETNX 是否成功。
func getExecSetNxBoolResult(result string) (string, bool) {
	// 使用空格分割 result，并检查是否至少有三个部分。
	keySlice := strings.Split(result, " ")
	if len(keySlice) < 3 {
		return "", false
	}

	// 使用冒号分割 result，并检查是否至少有两个部分。
	resultSlice := strings.Split(result, ":")
	if len(resultSlice) < 2 {
		return "", false
	}

	// 获取分割后的最后一个部分，去除空格后判断是否为 "true"。
	ResString := strings.TrimSpace(resultSlice[len(resultSlice)-1])
	if ResString == "true" {
		return keySlice[1], true
	}

	// 根据判断结果返回相应的键名和布尔值。
	return keySlice[1], false
}

func (l *mlock) MUnlock() error {
	if !l.needUnlock {
		return nil
	}
	return l.cache.Del(context.Background(), l.keys...).Err()
}

func (l *mlock) MPartialUnlock(keys []string) error {

	keysMap := make(map[string]struct{}, 0)
	for _, key := range l.keys {
		keysMap[key] = struct{}{}
	}

	for _, key := range keys {
		if _, ok := keysMap[key]; !ok {
			return fmt.Errorf("some key not in keys, key: %v, all keys: %v", key, l.keys)
		}
	}

	if !l.needUnlock {
		return nil
	}

	return l.cache.Del(context.Background(), keys...).Err()
}
