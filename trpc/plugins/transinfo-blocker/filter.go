package transinfoblocker

import (
	"context"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
)

// ClientFilter 客户端调用屏蔽metadata内容
func ClientFilter(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) error {
	ParseClientMetadata(ctx)
	return handler(ctx, req, rsp)
}

// ParseClientMetadata 修改metadata
func ParseClientMetadata(ctx context.Context) {
	msg := trpc.Message(ctx)
	metaData := msg.ClientMetaData()

	if len(metaData) == 0 { // 只有当前有传递给下游metaData的时候，才考虑这个逻辑
		return
	}

	transMetaData := make(map[string][]byte)
	callCfg := cfg.Default
	if rpcCfg, ok := cfg.RPCNameCfg[msg.ClientRPCName()]; ok {
		callCfg = rpcCfg
	}

	if callCfg == nil || callCfg.Mode == modeNone {
		return
	}

	for k, v := range metaData {
		if callCfg.Mode == modeWhitelist {
			if _, ok := callCfg.Set[k]; ok {
				transMetaData[k] = v
			}
		}

		if callCfg.Mode == modeBlacklist {
			if _, ok := callCfg.Set[k]; !ok {
				transMetaData[k] = v
			}
		}
	}

	// 重新设置过滤后的元数据
	msg.WithClientMetaData(transMetaData)
}
