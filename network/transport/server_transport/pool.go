package server_transport

import (
	"math"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/panjf2000/ants/v2"
)

const defaultBufferSize = 128 * 1024

type handleParam struct {
	req   []byte
	c     *tcpconn
	start time.Time
}

func (p *handleParam) reset() {
	p.req = nil
	p.c = nil
	p.start = time.Time{}
}

var handleParamPool = &sync.Pool{
	New: func() interface{} { return new(handleParam) },
}

// createRoutinePool 创建协程池
func createRoutinePool(size int) *ants.PoolWithFunc {
	if size <= 0 {
		size = math.MaxInt32
	}
	pool, err := ants.NewPoolWithFunc(size, func(args interface{}) {
		param, ok := args.(*handleParam)
		if !ok {
			logging.Tracef("routine pool args type error, shouldn't happen!")
			return
		}
		if param.c == nil {
			logging.Tracef("routine pool tcpconn is nil, shouldn't happen!")
			return
		}
		param.c.handleSync(param.req)
		param.reset()
		handleParamPool.Put(param)
	})
	if err != nil {
		logging.Tracef("routine pool create error:%v", err)
		return nil
	}
	return pool
}
