package setnx

import (
	"fmt"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/redis/client"
	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {
	redisClient := client.NewTestRedisClient()

	lock := NewLocker(redisClient)

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	locked, err := lock.Lock(prefix+"lock1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"lock2", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"lock3", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"lock1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)
}

func TestUnlock(t *testing.T) {
	redisClient := client.NewTestRedisClient()

	lock := NewLocker(redisClient)

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	locked, err := lock.Lock(prefix+"unlock", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	err = lock.Unlock()
	require.NoError(t, err)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"unlock", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)
}

func TestUnlockLockErr(t *testing.T) {
	redisClient := client.NewTestRedisClient()

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	lockSucc := NewLocker(redisClient)
	locked, err := lockSucc.Lock(prefix+"UnlockLockErr", time.Minute*2)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock := NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"UnlockLockErr", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	// lock failure, unlock no use
	err = lock.Unlock()
	require.NoError(t, err)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"UnlockLockErr", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	err = lockSucc.Unlock()
	require.NoError(t, err)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"UnlockLockErr", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)
}

func TestUnlockLockExpire(t *testing.T) {
	redisClient := client.NewTestRedisClient()

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	lockSucc := NewLocker(redisClient)
	locked, err := lockSucc.Lock(prefix+"UnlockLockExpire", time.Second*2)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock := NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"UnlockLockExpire", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	// lock failure, unlock no use
	err = lock.Unlock()
	require.NoError(t, err)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"UnlockLockExpire", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

	time.Sleep(time.Second * 2)

	lock = NewLocker(redisClient)
	locked, err = lock.Lock(prefix+"test1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)
}
