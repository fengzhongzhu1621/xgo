package nethttp

import (
	"fmt"
	"time"
)

type UserLocalSessionLock struct {
	Operator string
}

func (t *UserLocalSessionLock) WriteLockKey() string {
	return fmt.Sprintf("%s#%s", LocalSesionLockPrefix, t.Operator)
}

func (t *UserLocalSessionLock) TTL() time.Duration {
	return 120 * time.Second
}
