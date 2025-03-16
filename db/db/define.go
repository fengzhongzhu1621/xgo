package db

import "time"

const (
	TransactionIdHeader      = "transaction_id_string"
	TransactionTimeoutHeader = "transaction_timeout"

	// mongodb default transaction timeout is 1 minute.
	TransactionDefaultTimeout = 2 * time.Minute
)
