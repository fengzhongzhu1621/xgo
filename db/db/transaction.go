package db

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/network/nethttp"
)

// TxnOption TODO
type TxnOption struct {
	// transaction timeout time
	// min value: 5 * time.Second
	// default: 5min
	Timeout time.Duration
}

// TxnCapable TODO
type TxnCapable struct {
	Timeout   time.Duration `json:"timeout"`
	SessionID string        `json:"session_id"`
}

// AbortTransactionResult abort transaction result
type AbortTransactionResult struct {
	// Retry defines if the transaction needs to retry, the following are the scenario that needs to retry:
	// 1. the write operation in the transaction conflicts with another transaction,
	// then do retry in the scene layer with server times depends on conditions.
	Retry bool `json:"retry"`
}

// AbortTransactionResponse abort transaction response
type AbortTransactionResponse struct {
	nethttp.BaseResp       `json:",inline"`
	AbortTransactionResult `json:"data"`
}
