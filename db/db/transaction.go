package db

import (
	"net/http"
	"time"
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

// GetTXId get transaction id from http header
func GetTXId(header http.Header) string {
	return header.Get(TransactionIdHeader)
}

// GetTXTimeout get transaction timeout from http header
func GetTXTimeout(header http.Header) string {
	return header.Get(TransactionTimeoutHeader)
}

// SetTXId set transaction id to http header
func SetTXId(header http.Header, value string) {
	header.Set(TransactionIdHeader, value)
}

// SetTXTimeout set transaction timeout to http header
func SetTXTimeout(header http.Header, value string) {
	header.Set(TransactionTimeoutHeader, value)
}
